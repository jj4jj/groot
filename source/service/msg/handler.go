package msg

import (
	"github.com/gin-gonic/gin"
	"groot/proto/comm"
	"groot/proto/cserr"
	"groot/proto/csmsg"
	"groot/service/account/ctx"
	"groot/service/chat"
	"groot/service/friend"
)


func SendMsgToFriendUser(c * gin.Context, req *csmsg.CSMsgMsgSendReq, rsp * csmsg.CSMsgMsgSendRsp) cserr.ICSErrCodeError {
	uinMine := ctx.NeedAuthContextGetUin(c)
	uinFriend := req.Target.TargetId
	sendStatus := friend.DbFriendGetStatus(uinMine, uinFriend)
	if sendStatus != comm.FriendStatus_FRIEND_STATUS_NORMAL {
		return cserr.ErrFriendNotExist
	}
	recvStatus := friend.DbFriendGetStatus(uinFriend, uinMine)
	if recvStatus != comm.FriendStatus_FRIEND_STATUS_NORMAL {
		if sendStatus == comm.FriendStatus_FRIEND_STATUS_BLACK {
			return cserr.ErrFriendSetBeBlack
		}
		if recvStatus == comm.FriendStatus_FRIEND_STATUS_STRANGER {
			return cserr.ErrFriendSetBeRemoved
		}
		return cserr.ErrFriendNotExist
	}
	//
	sessionId := chat.GetFriendChatSessionId(uinMine, uinFriend)
	lastMsgUniqId,cerr := chat.DbChatSessionGetLastMsgUniqIdWithCheckSeq(sessionId, uint64(req.CheckSendSessionSeq))
	if cerr != nil {
		return cerr
	}
	newMsgUniqId,cerr2 := InsertChatMsgByHint(sessionId,uinMine, req.Content, lastMsgUniqId)
	if cerr2 != nil {
		return cerr
	}
	cerr = chat.DbChatSessionUpdateLastMsgUniqId(sessionId, newMsgUniqId)
	if cerr != nil {
		return cerr
	}
	//notify user new msg
	//notify.NotifyUserEvent(uinFriend,)

	return nil
}

func SendMsgToGroup(c * gin.Context, req *csmsg.CSMsgMsgSendReq, rsp * csmsg.CSMsgMsgSendRsp) cserr.ICSErrCodeError {
	return cserr.ErrNoImplement
}


func SendMsgToSystem(c * gin.Context, req *csmsg.CSMsgMsgSendReq, rsp * csmsg.CSMsgMsgSendRsp) cserr.ICSErrCodeError {
	return cserr.ErrNoImplement
}



func MsgSend(c *gin.Context, req *csmsg.CSMsgMsgSendReq, rsp *csmsg.CSMsgMsgSendRsp)(err cserr.ICSErrCodeError) {
	err = cserr.ErrInternal
	switch req.Target.TargetType  {
	case comm.MsgTargetType_TARGET_TYPE_USER:
		err = SendMsgToFriendUser(c, req, rsp)
	case comm.MsgTargetType_TARGET_TYPE_GROUP:
		err = SendMsgToGroup(c, req, rsp)
	case comm.MsgTargetType_TARGET_TYPE_SYSTEM:
		err = SendMsgToSystem(c, req, rsp)
	default:
		err = cserr.ErrBadRequest
	}
	return

}
