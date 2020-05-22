package notify

import (
	"github.com/golang/protobuf/proto"
	"groot/proto/comm"
	"groot/proto/ssmsg"
	"groot/service/stream"
)

type (
	NotifyUserStreamLogic struct {
		Base	* stream.LogicBase
	}
)

func (n *NotifyUserStreamLogic)NeedAuth() bool {
	return false
}


func (logic *NotifyUserStreamLogic) OnRecvUserEventEvent (evt ssmsg.SvcEventType, param *comm.CSUserEvent){
	//check online device and notify
	if logic.Base.UserIsOnline(param.Uin) {
		logic.Base.SendUserMsg(param.Uin, param)
	}
}

func MakeMulticastMsg(param *ssmsg.EvtMulticastMsg) []*comm.CSUserEvent {
	var list []*comm.CSUserEvent
	for i := range param.Ctx.ParamList {
		eventBase := comm.CSUserEvent{}
		proto.Merge(&eventBase, param.EventTemplate)
		eventBase.Uin = param.Ctx.ParamList[i].Uin
		eventBase.EventId = param.Ctx.ParamList[i].EventId
		list = append(list, &eventBase)
	}
	return list
}

func(l *NotifyUserStreamLogic) MulticastUserEventMsg(evt ssmsg.SvcEventType, param *ssmsg.EvtMulticastMsg){
	msgList := MakeMulticastMsg(param)
	l.Base.Multicast(msgList)
}

func MakeBroadcastMsg(param *ssmsg.EvtBroadcastMsg) *comm.CSUserEvent {
	var msg comm.CSUserEvent
	proto.Merge(&msg, param.EventTemplate)
	return  &msg
}

func(l *NotifyUserStreamLogic) BroadcastUserEventMsg(evt ssmsg.SvcEventType, param *ssmsg.EvtBroadcastMsg){
	msg := MakeBroadcastMsg(param)
	l.Base.Broadcast(msg)
}


