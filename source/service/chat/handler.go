package chat

import (
	"github.com/gin-gonic/gin"
	"groot/proto/cserr"
	"groot/proto/csmsg"
	"groot/service/account/ctx"
	"groot/service/chat/model"
)

func UpdateSetting(c *gin.Context, req *csmsg.CSMsgChatUpdateSettingReq, rsp *csmsg.CSMsgChatUpdateSettingRsp) (err cserr.ICSErrCodeError) {
	err = cserr.ErrInternal

	uin := ctx.NeedAuthContextGetUin(c)

	dbSetting := model.DbChatSetting{
		IsQuiet: req.IsQuiet,
		ChatBackGround: req.Background,
	}
	err = DbChatSettingUpdate(uin, req.Target.TargetType, req.Target.TargetId, &dbSetting)
	return
}


func Complain(c * gin.Context, req *csmsg.CSMsgChatComplainReq, rsp *csmsg.CSMsgChatComplainRsp)(err cserr.ICSErrCodeError) {
	err = cserr.ErrInternal
	uin := ctx.NeedAuthContextGetUin(c)
	err = DbChatCompainInsert(uin, req.Target.TargetType, req.Target.TargetId, req.ComplainType, req.Content)
	return
}