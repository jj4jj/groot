package notify

import (
	"github.com/gin-gonic/gin"
	"groot/proto/cserr"
	"groot/proto/csmsg"
	"groot/service/account/ctx"
	"groot/service/notify/crud"
	"groot/service/notify/model"
	"groot/sfw/util"
)


func UserEventGetList(c * gin.Context, req *csmsg.CSMsgNotifyUserEventGetListReq,
	rsp *csmsg.CSMsgNotifyUserEventGetListRsp) (err cserr.ICSErrCodeError) {
	uin := ctx.NeedAuthContextGetUin(c)
	err = cserr.ErrInternal
	//limit
	if req.Limit > 100 {
		return cserr.ErrBadRequest
	}

	userEventList,e := crud.DbNotifyUserEventGetList(uin, req.Offset, req.Limit, req.EventIdFrom)
	if !util.ErrOK(e) {
		return e
	}

	for i := range userEventList {
		rsp.EventList = append(rsp.EventList, model.DbNotifyUserEventToCsUserEvent(userEventList[i]))
	}

	err = nil
	return

}

func UserEventRemove(c * gin.Context, req *csmsg.CSMsgNotifyUserEventRemoveReq,
	rsp *csmsg.CSMsgNotifyUserEventRemoveRsp) (err cserr.ICSErrCodeError) {
	err = cserr.ErrNoImplement
	uin := ctx.NeedAuthContextGetUin(c)
	err = cserr.ErrInternal
	if len(req.EventIdList) > 100 {
		return cserr.ErrBadRequest
	}
	err = crud.DbNotifyUserEventRemove(uin, req.EventIdList)
	return
}


func NotifyUpdateUserEvent(uin string,event *model.DbNotifyUserEvent, updateType csmsg.UserEventUpdateType,
	intParam int64, strParam string, updateParam []byte ) cserr.ICSErrCodeError {
	return cserr.ErrTodo
}


func UserEventUpdate(c * gin.Context, req *csmsg.CSMsgNotifyUserEventUpdateReq,
	rsp *csmsg.CSMsgNotifyUserEventUpdateRsp) (err cserr.ICSErrCodeError) {
	err = cserr.ErrNoImplement
	uin := ctx.NeedAuthContextGetUin(c)
	//
	dbEvent,e := crud.DbNotifyUserEventGetByUserEventId(req.EventId)
	if !util.ErrOK(e) {
		return e
	}

	//
	e = NotifyUpdateUserEvent(uin, dbEvent, req.UpdateType, req.IntParam, req.StrParam, req.UpdateParam)
	if !util.ErrOK(e) {
		return e
	}


	return nil
}


