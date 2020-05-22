package feedback

import (
	"github.com/gin-gonic/gin"
	"groot/proto/cserr"
	"groot/proto/csmsg"
	"groot/service/account/ctx"
)

func FeedbackCreate(c *gin.Context, req *csmsg.CSMsgFeedbackCreateReq, rsp *csmsg.CSMsgFeedbackCreateRsp)(err cserr.ICSErrCodeError) {
	err = cserr.ErrInternal
	uin := ctx.NeedAuthContextGetUin(c)
	err = DbFeedbackCreate(uin, req.Content, req.Qq, req.Ext,req.Tel, req.Device)
	return
}