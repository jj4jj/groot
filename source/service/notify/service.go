package notify

import (
	"groot/comm/conf"
	"groot/comm/constk"
	"groot/proto/ssmsg"
	"groot/service/notify/crud"
	"groot/service/notify/model"
	"groot/service/stream"
	"groot/sfw/app"
	"groot/sfw/db"
	sfw_event "groot/sfw/event"
	"groot/sfw/task"
	"groot/sfw/util"
)

type (
	//logic state
	NotifyService struct {
		notifyStream NotifyUserStreamLogic
	}
)



func (logic *NotifyService) SyncDb() error {
	config := conf.GetAppConfig()
	return sfw_db.AutoMigrateSplitDbAndTables(crud.NotifyDbTbNameBase, config.Notify.UseDb.SplitDbNum,
		config.Notify.UseDb.SplitTbNum, &model.DbNotifyUserEvent{})
}


func initCsRpc(svc *app.ServiceCtx){
	//cs rpc
	svc.BindRpc("user_event_get_list", UserEventGetList)
	svc.BindRpc("user_event_remove", UserEventRemove)
	svc.BindRpc("user_event_update", UserEventUpdate)
}

func startNotifyWorker(){
	//1 pull 1 send
	topicEmailSend := util.GetBrokerTopicName(constk.BRK_TOPIC_NOTIFY_SERVICE, NOTIFY_CHANNEL_SMS)
	go task.RunWorker(topicEmailSend, 2, NotifyRealSendMsg)

	topicSmsSend := util.GetBrokerTopicName(constk.BRK_TOPIC_NOTIFY_SERVICE, NOTIFY_CHANNEL_EMAIL)
	go task.RunWorker(topicSmsSend, 2, NotifyRealSendMsg)
}


func (logic *NotifyService) Init(svc *app.ServiceCtx) error {


	initCsRpc(svc)

	startNotifyWorker()

	logic.notifyStream = NotifyUserStreamLogic {
		Base: stream.MakeLogicBase(&logic.notifyStream),
	}

	//listen user event task
	sfw_event.ListenEvent(ssmsg.SvcEventType_SVC_EVT_NOTIFY_BROADCAST, logic.notifyStream.MulticastUserEventMsg)
	sfw_event.ListenEvent(ssmsg.SvcEventType_SVC_EVT_NOTIFY_MULTICAST, logic.notifyStream.MulticastUserEventMsg)
	sfw_event.ListenEvent(ssmsg.SvcEventType_SVC_EVT_NOTIFY_USER_EVENT, logic.notifyStream.OnRecvUserEventEvent)

	//listen user websocket but no request receive (reqInst is nil) //csmsg.CSMsgHttpFrame{}
	svc.BindStream("push", logic.notifyStream.Base, nil, 60)


	return nil
}

func Initialize() {
	//todo use config
	sfw_db.AddDbUseMap(crud.NotifyDbTbNameBase, "prod-misc")
	app.RegisterService("notify", &NotifyService{}, "main", "need_auth")

}
