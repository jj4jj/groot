package friend

import (
	"github.com/lexkong/log"
	"groot/comm/conf"
	"groot/proto/cserr"
	"groot/proto/ssmsg"
	"groot/service/friend/model"
	"groot/sfw/app"
	sfw_db "groot/sfw/db"
	sfw_event "groot/sfw/event"
	"groot/sfw/util"
)

//logic state of account servica
type FriendService struct {
}


func(f *FriendService) SyncDb() error {
	config := conf.GetAppConfig()
	return sfw_db.AutoMigrateSplitDbAndTables(FriendDbNameBase, config.Friend.UseDb.SplitDbNum, config.Friend.UseDb.SplitTbNum,
		&model.DbFriendship{}, &model.DbFriendSetting{})
}


func (f *FriendService) Init(svc *app.ServiceCtx) error {
	if util.CheckError(initFriendDbCtx(),"init friend db ctx") {
		return cserr.ErrDb
	}

	svc.BindRpc("ask", FriendAsk)
	svc.BindRpc("answer", FriendAnswer)
	svc.BindRpc("set_block", FriendSetBlock)
	svc.BindRpc("remove", FriendRemove)
	svc.BindRpc("list", FriendList)
	svc.BindRpc("update_setting", FriendUpdateSetting)

	//code listen account register and login
	sfw_event.ListenEvent(ssmsg.SvcEventType_SVC_EVT_ACCOUNT_REGISTER, func(evt ssmsg.SvcEventType, param *ssmsg.SvcEvtAccountRegister) {
		log.Debugf("listen event:%d fired param:%v", evt, param)
	})
	sfw_event.ListenEvent(ssmsg.SvcEventType_SVC_EVT_ACCOUNT_LOGIN, func(evt ssmsg.SvcEventType, param *ssmsg.SvcEvtAccountLogin) {
		log.Debugf("listen event:%d fired param:%v", evt, param)
	})
	return nil
}

func Initialize() {
	config := conf.GetAppConfig()
	sfw_db.UpdateDbUseMap(config.Friend.UseDb.DbUseMap)
	app.RegisterService("friend", &FriendService{}, "main", "need_auth")
}
