package csuser

import (
	"groot/proto/comm"
	"groot/proto/csmsg"
	"groot/sfw/util"
)

func (user *CSUser)AddFriend(tel string) {
	ticket := user.SearchTel(tel)
	req := &csmsg.CSMsgFriendAskReq{
		UserTicket: ticket,
	}
	rsp := &csmsg.CSMsgFriendAskRsp{}
	util.CheckError(user.server.Call("friend.ask", req, rsp),"")
}

func (user * CSUser)GetFriendList(){
	req := &csmsg.CSMsgFriendListReq{
		TagFilter: comm.FriendStatus_FRIEND_STATUS_NONE,
		Limit: 100,
	}
	rsp := &csmsg.CSMsgFriendListRsp{}
	util.CheckError(user.server.Call("friend.ask", req, rsp),"")
	user.friend_list = rsp.FriendList
}

func (user * CSUser)SendFriendMsg(friendUin string, msg string){
	req := &csmsg.CSMsgMsgSendReq{
		Target: &comm.MsgTarget{
			TargetType: comm.MsgTargetType_TARGET_TYPE_USER,
			TargetId: friendUin,
		},
		Content: &csmsg.MsgContent{
			MsgType: comm.MsgType_MSG_TYPE_CHAT,
			Content: []byte(msg),
		},
		CheckSendSessionSeq: user.check_session_msg_seq,
	}
	rsp := &csmsg.CSMsgMsgSendRsp{}
	util.CheckError(user.server.Call("msg.send", req, rsp),"")
	user.check_session_msg_seq =  rsp.CheckSendSessionSeq
}


