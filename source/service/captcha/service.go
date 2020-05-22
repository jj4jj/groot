package captcha

import (
	"github.com/jinzhu/gorm"
	"groot/comm/conf"
	"groot/service/captcha/model"
	"groot/sfw/app"
	"groot/sfw/db"
)

var (
	CaptchaDb                * gorm.DB
	CaptchaDbNameBase        string
	DbCaptchaStateTbNameBase string
)

func init(){
	CaptchaDbNameBase = "db_captcha"
	DbCaptchaStateTbNameBase = sfw_db.GetModelTableName(&model.DbCaptchaState{})
}

type CaptchaService struct {
}

func (c *CaptchaService) SyncDb() error {
	//migrate
	return sfw_db.AutoMigrateTable(CaptchaDbNameBase, &model.DbCaptchaState{})
}

func (c *CaptchaService) Init(service *app.ServiceCtx) error {
	CaptchaDb = sfw_db.GetDb(CaptchaDbNameBase)
	return nil
}

func Initialize() {
	config := conf.GetAppConfig()
	sfw_db.UpdateDbUseMap(config.Captcha.UseDb.DbUseMap)
	app.RegisterService("captcha", &CaptchaService{}, "main")
}
