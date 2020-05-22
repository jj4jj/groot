package account

import (
	"github.com/lexkong/log"
	"groot/comm/conf"
	"groot/service/account/crud"
	"groot/service/account/model"
	"groot/sfw/app"
	sfw_db "groot/sfw/db"
	"groot/sfw/util"
)

//logic state of account service
type AccountService struct {}


func (logic * AccountService) SyncDb() error {
	config := conf.GetAppConfig()
	if e := sfw_db.AutoMigrateSplitDbAndTables(crud.AccountDbNameBase, config.Account.UseDb.SplitDbNum,
		config.Account.UseDb.SplitTbNum, &model.DbUser{}, &model.DbUserLoginEmail{},
		&model.DbUserLoginPhone{}, &model.DbUserLoginName{}, &model.DbUserLoginDevice{},
		&model.DbUserTokenState{}); e != nil {
		return e
	}

	//create system user
	if e := crud.CreateSystemDbUser(); e != nil {
		return e
	}

	return nil
}


func (logic *AccountService) Init(svc *app.ServiceCtx) error {
	//init db
	//using
	//post
	svc.BindRpc("login_or_register_get_sms_code",LoginOrRegisterGetSmsCode)
	svc.BindRpc("login_or_register",  LoginOrRegister)
	svc.BindRpc("logout",Logout, "need_auth","need_user")
	//token 续期逻辑[换取新的token.每10分钟可执行一次]
	svc.BindRpc("token_renew", TokenRenew, "need_auth")
	svc.BindRpc("change_password", ChangePassword, "need_auth")
	svc.BindRpc("check_user_exist", CheckUserExist, "need_auth")
	svc.BindRpc("change_telephone_old_telephone_get_sms_code", ChangeTelephoneOldTelephoneGetSmsCode, "need_auth")
	svc.BindRpc("change_telephone_new_telephone_get_sms_code", ChangeTelephoneNewTelephoneGetSmsCode, "need_auth")
	svc.BindRpc("change_telephone_confirm",ChangeTelephoneConfirm,
	"need_auth")
	svc.BindRpc("search_user", SearchUser, "need_auth")


	return nil
}

func Initialize() {
	config := conf.GetAppConfig()
	util.CheckError(LoadAuthRSAPrivateKey(config.AuthRsaPrivateKey),"load auth key")
	//
	if e := crud.InitAccountDbCtx(); e != nil {
		log.Errorf(e, "init account db fail")
	}

	app.AddMiddleWare("need_auth", AccountNeedAuthMiddleWare())
	app.AddMiddleWare("need_user", AccountNeedUserMiddleWare())
	app.RegisterService("account", &AccountService{}, "main")

}
