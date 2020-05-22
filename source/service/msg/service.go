package msg

import (
	"groot/comm/conf"
	"groot/service/msg/model"
	"groot/sfw/app"
	"groot/sfw/db"
)


type MsgService struct {}

//split tb = 1000 , split db num = 10000
func (c *MsgService) SyncDb() error {
	config := conf.GetAppConfig()
	MsgDbNameBase := "db_msg"
	var err error
	if config.RunEnv == "dev" {
		//using simple msg tb 0 , 0
		err = sfw_db.AutoMigrateSplitDbAndTables(MsgDbNameBase,config.Msg.EnableSplitDbBegIdx,
			1, &model.DbChatMsg{})
	} else {
		//end+1=total num
		err = sfw_db.AutoMigrateSplitDbAndTables(MsgDbNameBase,config.Msg.EnableSplitDbEndIdx+1,
			model.DbChatMsgTbMaxNum, &model.DbChatMsg{})
	}
	return err
}


func (c *MsgService) Init(service *app.ServiceCtx) error {

	if e:= initMsgDbCtx(); e != nil {
		return e
	}

	service.BindRpc("send", MsgSend)


	return nil
}

func Initialize() {
	//
	config := conf.GetAppConfig()
	//todo @hex use config
	sfw_db.UpdateDbUseMap(conf.AuxGenDbUseMapRoundRobinDbConnx(MsgDbNameBase, int(config.Msg.EnableSplitDbEndIdx+1),
		"prod-msg", 1))
	app.RegisterService("msg", &MsgService{}, "main", "need_auth")
}
