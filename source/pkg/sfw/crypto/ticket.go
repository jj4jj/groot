package crypto

import (
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"github.com/lexkong/log"
	"strings"
	"time"
)

//给任意结构生成票据.使得不需要存储DB,只需要检查状态即可
func GenerateMsgTicket(pbmsg proto.Message, timeout int, secret []byte) string {
	//
	b, e := proto.Marshal(pbmsg)
	if e != nil {
		log.Errorf(e, "gen msg ticket error for msg marshal error !")
		return ""
	}
	timeNow := uint32(time.Now().Unix())
	ctx := MsgTicketBaseCtx{
		IssueTime:   timeNow,
		ExpiredTime: timeNow + uint32(timeout),
		PayloadCtx:  b,
	}

	sign := GenerateSha256WithRsa(b)

	b, _ = proto.Marshal(&ctx)
	//aes
	aes := EncryptWithAES(b, secret)

	//base64
	b64 := base64.StdEncoding.EncodeToString(aes)

	return sign + "." + b64
}

func VerifyMsgTicket(ticket string, secret []byte, appmsg proto.Message) bool {

	ctx := ParseMsgTicket(ticket, secret, appmsg)
	if ctx == nil {
		return false
	}

	//time
	timeNow := time.Now().Unix()
	if uint32(timeNow) > ctx.ExpiredTime {
		log.Warnf("ticket expird ")
		return false
	}

	return true
}

func ParseMsgTicket(ticket string, secret []byte, pbmsg proto.Message) *MsgTicketBaseCtx {

	bs := strings.Split(ticket, ".")
	if len(bs) != 2 {
		return nil
	}
	sign := bs[0]
	tp := bs[1]

	db64, e := base64.StdEncoding.DecodeString(tp)
	if e != nil {
		log.Errorf(e, "base64 decode ticket error")
		return nil
	}

	dx := DecryptWithAES(db64, secret)

	if VerifySha256WithRsa(dx, sign) == false {
		log.Warnf("signature verify fail")
		return nil
	}

	var ctx MsgTicketBaseCtx
	e = proto.Unmarshal(dx, &ctx)
	if e != nil {
		log.Errorf(e, "unmashal error")
		return nil
	}

	if pbmsg != nil {
		e = proto.Unmarshal(ctx.PayloadCtx, pbmsg)
		if e != nil {
			log.Warnf("unmarshal payload ctx error !")
		}
	}

	return &ctx

}
