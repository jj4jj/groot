package account

import "C"
import (
	"github.com/gin-gonic/gin"
	"groot/proto/comm"
	"groot/proto/cserr"
	"groot/proto/csmsg"
	"groot/proto/ssmsg"
	"groot/service/account/ctx"
	"groot/service/account/crud"
	"groot/service/account/model"
	"groot/sfw/crypto"
	"time"
)


func CheckUserExist(c *gin.Context, req *csmsg.CSMsgAccountCheckUserExistReq, rsp *csmsg.CSMsgAccountCheckUserExistRsp) (err cserr.ICSErrCodeError) {
	phone := model.DbUserLoginPhone{
		Phone: req.CheckStr,
	}
	db, tb := crud.GetAccountDbAndTbNameByTelephone(req.CheckStr, crud.AccountDbUserLoginPhoneTbNameBase)
	if db.Table(tb).NewRecord(&phone) {
		rsp.Exist = false
	} else {
		rsp.Exist = true
	}
	return nil
}


func GenerateUserCheckTicketWithMsg(UinMine,UinCheck string, source comm.UserSourceType, timeoutSeconds int, cookie string) string {
	checkRet := ssmsg.UserCheckResultTicket {
		Source:    source,
		UinSigner: UinMine,
		UinCheck:  UinCheck,
		Cookie:    cookie,
	}
	if timeoutSeconds == 0 {
		timeoutSeconds = 86400
	}
	secret := crud.GetAccountTicketSecretByUin(UinMine)
	return crypto.GenerateMsgTicket(&checkRet, timeoutSeconds, []byte(secret))
}

func VerifyUserCheckTicketWithUin(uin,ticket string) *ssmsg.UserCheckResultTicket {
	checkRet := ssmsg.UserCheckResultTicket{}
	secret := crud.GetAccountTicketSecretByUin(uin)
	if false == crypto.VerifyMsgTicket(ticket,[]byte(secret), &checkRet) {
		return nil
	}
	return &checkRet
}

func SearchUser(c *gin.Context, req *csmsg.CSMsgAccountSearchUserReq, rsp *csmsg.CSMsgAccountSearchUserRsp) (err cserr.ICSErrCodeError) {
	err = cserr.ErrInternal
	uin := ctx.NeedAuthContextGetUin(c)
	var dbUser *model.DbUser
	var source comm.UserSourceType
	//
	if len(req.Code) > 0 {
		//qr
		source = comm.UserSourceType_USER_SOURCE_SCAN_QR
		//todo @hex
	} else {
		source = comm.UserSourceType_USER_SOURCE_SEARCH_UID
		dbUser, err = crud.DbUserGetByUserName(req.Search)
		if dbUser == nil {
			dbUser, err = crud.DbUserGetByTel(req.Search)
			if dbUser != nil {
				source = comm.UserSourceType_USER_SOURCE_SEARCH_TEL
			}
		}
	}
	if dbUser != nil {
		rsp.Source = source
		rsp.Uin = uin
		expiredTimeSeconds := 86400*3
		rsp.TicketExpiredTimeStamp = time.Now().Unix() + int64(expiredTimeSeconds)
		rsp.UserTicket =
			GenerateUserCheckTicketWithMsg(uin, dbUser.Uin, source, expiredTimeSeconds, "account.search")
		err = nil
	} else {
		err = cserr.ErrAccountUserNotExist
	}
	return
}






