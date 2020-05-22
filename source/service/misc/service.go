package misc

import (
	"groot/sfw/app"
)

type MiscService struct {}

func (c *MiscService) SyncDb() error {
	//config := conf.GetAppConfig()
	return nil
}


func (c *MiscService) Init(service *app.ServiceCtx) error {
	return nil
}

func Initialize() {
	app.RegisterService("misc", &MiscService{}, "main", "need_auth")
}


