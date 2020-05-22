package notify

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/lexkong/log"
	"groot/comm/constk"
	"groot/proto/comm"
	"groot/proto/ssmsg"
	"groot/service/notify/crud"
	"groot/sfw/event"
	"groot/sfw/task"
	"groot/sfw/util"
)

type (
	NotifySceneType   int
	NotifyChannelType int
	M                 map[string]interface{}
	UserEventParamType	proto.Message
)

const (
	NOTIFY_SCENE_ACCOUNT_REGISTER NotifySceneType = 1 + iota
	NOTIFY_SCENE_ACCOUNT_LOGIN
	NOTIF_SCENE_ACCOUNT_RESET_PASSWD
)

const (
	NOTIFY_CHANNEL_SMS NotifyChannelType = 1 + iota
	NOTIFY_CHANNEL_EMAIL
	NOTIFY_CHANNEL_ONLINE
	NOTIFY_CHANNEL_PUSH
)

//return task id and error
func Send(target string, scene NotifySceneType, ch NotifyChannelType, ctx M) error {
	//just pub a task
	topic := util.GetBrokerTopicName(constk.BRK_TOPIC_NOTIFY_SERVICE, ch)
	m := & ssmsg.NotifyServiceMsgParam{
		Channel: int32(ch),
		Scene:   int32(scene),
		Target:  target,
		Content: []byte(fmt.Sprint(ctx)),
		From:    "todo",
	}
	log.Debugf("notify send topic:%s msg:%s ...", topic, m.String())
	return task.AddTask(topic, m)
}


func NotifyBroadcastUserEvent(evType comm.UserEventType, intParam int64, strParam string, filter *ssmsg.EvtTargetFilter) {
	broadcastMsg := ssmsg.EvtBroadcastMsg {
		EventTemplate: & comm.CSUserEvent {
			Uin: "",
			EventId: "",
			EvtType: evType,
			IntParam:intParam,
			StrParam:strParam,
		},
		Filter: filter,
	}
	sfw_event.FireEvent(ssmsg.SvcEventType_SVC_EVT_NOTIFY_BROADCAST, &broadcastMsg)
}


func NotifyMultiUserEventNew(saveDb bool,UinList []string, evType comm.UserEventType, intParam int64,
	strParam string, evParam UserEventParamType) {
	multicastMsg := ssmsg.EvtMulticastMsg {
		EventTemplate: & comm.CSUserEvent {
			Uin: "",
			EventId: "",
			EvtType: evType,
			IntParam:intParam,
			StrParam:strParam,
		},
		Ctx: &ssmsg.MultiUserEventCtx {
		},
	}
	for _,Uin := range UinList {
		if saveDb {
			//todo user group add db todo
			dbEvent,e :=  crud.DbUserNotifyEventAdd(Uin, evType, intParam, strParam, evParam)
			if util.ErrOK(e) {
				log.Errorf(e, "user notify event db fail")
			}
			multicastMsg.Ctx.ParamList = append(multicastMsg.Ctx.ParamList,
				&ssmsg.UserEventParam{
					Uin: Uin,
					EventId: dbEvent.UserEventId,
				})
		} else {
			multicastMsg.Ctx.ParamList = append(multicastMsg.Ctx.ParamList,
				&ssmsg.UserEventParam{
					Uin: Uin,
					EventId: "",
				})
		}
	}
	sfw_event.FireEvent(ssmsg.SvcEventType_SVC_EVT_NOTIFY_USER_EVENT, &multicastMsg)
}

func NotifyUserEventNewWithSaveDb(Uin string, evType comm.UserEventType, intParam int64, strParam string, evParam UserEventParamType) {
	dbEvent,e :=  crud.DbUserNotifyEventAdd(Uin, evType, intParam, strParam, evParam)
	if util.ErrOK(e) {
		log.Errorf(e, "user notify event db fail")
	}
	evmsg := comm.CSUserEvent {
		Uin: Uin,
		EventId: dbEvent.UserEventId,
		EvtType: evType,
		IntParam:intParam,
		StrParam:strParam,
	}
	sfw_event.FireEvent(ssmsg.SvcEventType_SVC_EVT_NOTIFY_USER_EVENT, &evmsg)
}

func NotifyUserEventNewWithoutSaveDb(Uin string, evType comm.UserEventType, intParam int64, strParam string, evParam UserEventParamType) {
	dbEvent,e :=  crud.DbUserNotifyEventAdd(Uin, evType, intParam, strParam, evParam)
	if util.ErrOK(e) {
		log.Errorf(e, "user notify event db fail")
	}
	evmsg := comm.CSUserEvent {
		Uin: Uin,
		EventId: dbEvent.UserEventId,
		EvtType: evType,
		IntParam:intParam,
		StrParam:strParam,
	}
	sfw_event.FireEvent(ssmsg.SvcEventType_SVC_EVT_NOTIFY_USER_EVENT, &evmsg)
}


func NotifyUserUpdateEvent(Uin string, updateEvent *comm.CSUserEvent) {
	sfw_event.FireEvent(ssmsg.SvcEventType_SVC_EVT_NOTIFY_USER_EVENT, updateEvent)
}


func NotifyRealSendMsg(msg *ssmsg.NotifyServiceMsgParam) {
	if msg != nil {
		log.Debugf("todo : notify real send msg:%s", msg.String())
	} else {
		log.Errorf(nil, "msg error !")
	}
}




