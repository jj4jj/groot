package chat

import (
	"groot/comm/conf"
	"groot/service/chat/model"
	"groot/sfw/app"
	sfw_db "groot/sfw/db"
)


type ChatService struct {
}


func (c * ChatService) SyncDb() error {
	config := conf.GetAppConfig()

	if e := sfw_db.AutoMigrateTable(ChatComplainDbName, &model.DbChatComplain{});
		e != nil {
		return e
	}
	if e := sfw_db.AutoMigrateSplitDbAndTables(ChatDbNameBase, config.Chat.UseDb.SplitDbNum,
		config.Chat.UseDb.SplitDbNum, &model.DbChatSetting{}); e != nil {
		return e
	}
		return nil
}

func (c *ChatService) Init(service *app.ServiceCtx) error {

	if e:= initChatDbCtx(); e!= nil {
		return e
	}

	service.BindRpc("update_setting", UpdateSetting)
	service.BindRpc("complain", Complain)

	return nil
}

func Initialize() {
	//
	config := conf.GetAppConfig()
	sfw_db.UpdateDbUseMap(config.Chat.UseDb.DbUseMap)
	app.RegisterService("chat", &ChatService{}, "main", "need_auth")
}
