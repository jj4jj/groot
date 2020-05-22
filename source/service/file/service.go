package file

import (
	"github.com/jinzhu/gorm"
	"groot/comm/conf"
	"groot/sfw/app"
)

var (
	DB *gorm.DB
)

//logic state of account servica
type FileService struct {}


func (logic * FileService) SyncDb () error {return nil}

func (logic *FileService) Init(svc *app.ServiceCtx) error {
	svc.RouteRaw("GET", "get", GetFile)
	svc.RouteRaw("POST", "upload", UploadFile)
	return nil
}


func Initialize() {
	conf := conf.GetAppConfig()
	UploadFileDir = conf.DevEnv.UploadFileDir
	app.RegisterService("file", &FileService{}, "main")
}
