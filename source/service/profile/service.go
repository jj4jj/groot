package profile

import (
	"groot/comm/conf"
	"groot/service/profile/crud"
	"groot/service/profile/model"
	"groot/sfw/app"
	"groot/sfw/db"
)

type ProfileService struct {

}

func (p * ProfileService) SyncDb() error  {
	config := conf.GetAppConfig()
	return sfw_db.AutoMigrateSplitDbAndTables(crud.ProfileDbNameBase, config.Profile.UseDb.SplitDbNum, config.Profile.UseDb.SplitTbNum,
		model.DbUserBaseInfo{}, model.DbUserProfileInfo{}, model.DbUserRegisterInfo{}, model.DbUserRegisterInfo{})
}


func (p * ProfileService) Init (svc * app.ServiceCtx) error {

	if e:= crud.InitProfileDbCtx(); e!=nil {
		return e
	}

	svc.BindRpc("set_sex", ProfileSetSex)
	svc.BindRpc("set_icon", ProfileSetIcon)
	svc.BindRpc("set_location", ProfileSetLocation)
	svc.BindRpc("update_setting", ProfileUpdateSetting)
	svc.BindRpc("set_nick_name", ProfileSetNickName)
	svc.BindRpc("set_desc", ProfileSetDesc)
	svc.BindRpc("set_enable_search_by_tel", ProfileSetEnableSearchByTel)
	svc.BindRpc("get_user_base", GetUserBase)
	svc.BindRpc("get_self_user_info", GetSelfUserInfo)


	svc.TcsListen("TcsUserRegisterInfoSave", TcsUserRegisterInfoSave)


	return nil
}



func Initialize() {
	config := conf.GetAppConfig()

	//todo use config
	sfw_db.UpdateDbUseMap(config.Profile.UseDb.DbUseMap)

	app.RegisterService("profile", &ProfileService{}, "main",
		"need_auth","need_user")

}
