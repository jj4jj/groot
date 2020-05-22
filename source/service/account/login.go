package account

import (
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"groot/comm/conf"
	"groot/comm/constk"
	"groot/proto/comm"
	"groot/proto/cserr"
	"groot/proto/csmsg"
	"groot/proto/ssmsg"
	"groot/service/account/crud"
	"groot/service/account/model"
	"groot/service/captcha"
	"groot/service/notify"
	"groot/sfw/crypto/token"
	"groot/sfw/db"
	"groot/sfw/event"
	"groot/sfw/task"
	"groot/sfw/util"
	"time"
)

//注册发送验证码.如果账号已经存在则返回已存在错误.否则签发验证码.
func LoginOrRegisterGetSmsCode(c *gin.Context,
	req *csmsg.CSMsgAccountLoginOrRegisterGetSmsCodeReq,rsp *csmsg.CSMsgAccountLoginOrRegisterGetSmsCodeRsp) (err cserr.ICSErrCodeError) {
	//validate request todo (check phone formating)

	//Search the user is exist?
	phone := model.DbUserLoginPhone {
		Phone: req.Tel,
	}

	db, tb := crud.GetAccountDbAndTbNameByTelephone(req.Tel, crud.AccountDbUserLoginPhoneTbNameBase)

	//存在并且绑定了User
	var captchaType constk.CaptchaSceneType
	var notifyScene notify.NotifySceneType
	if db.Table(tb).First(&phone).RecordNotFound() || phone.Uin == "" {
		//user is exist already , login
		captchaType = constk.CAPTCHA_ACCOUNT_REGISTER
		notifyScene = notify.NOTIFY_SCENE_ACCOUNT_REGISTER
		rsp.AsRegister = true
	} else {
		captchaType = constk.CAPTCHA_ACCOUNT_LOGIN
		notifyScene = notify.NOTIFY_SCENE_ACCOUNT_LOGIN
		rsp.AsRegister = false
	}

	//申请验证码.
	code, err := captcha.ApplyCode(req.Tel, captchaType,
		captcha.CaptchaGetSMSCodeDefaultApplyOptions())
	if !util.ErrOK(err) {
		return cserr.ErrCaptcha
	}

	//发送验证码.调用grpc服务.
	if e := notify.Send(req.Tel, notifyScene, notify.NOTIFY_CHANNEL_SMS,
		notify.M{"Code": code}); e != nil {
		log.Errorf(err, "notify code send error !")
		return cserr.ErrCaptcha
	}

	//@hex todo just for debug
	rsp.Ticket = ""
	if conf.GetAppConfig().RunEnv == "dev" {
		rsp.Ticket = code
	}
	return nil
}

//just register
func userRegisterConfirm(c *gin.Context, req *csmsg.CSMsgAccountLoginOrRegisterReq,
	rsp *csmsg.CSMsgAccountLoginOrRegisterRsp) cserr.ICSErrCodeError {
	log.Infof("user register device id is %s tel:%s push token:%s", req.LoginDevice.DeviceId,
		req.AuthId, req.LoginDevice.PushToken)
	//create user by telephone (telephone)
	dbuser, err := crud.DbUserCreateByTelephone(req.AuthId)
	if !util.ErrOK(err) || dbuser == nil {
		log.Errorf(err, "create DB user fail with tel:%s !", req.AuthId)
		return cserr.ErrDb
	}

	//其他子系统数据初始化通知.其他子系统收到通知后初始化
	sfw_event.FireEvent(ssmsg.SvcEventType_SVC_EVT_ACCOUNT_REGISTER, &ssmsg.SvcEvtAccountRegister{
		Uin: dbuser.Uin,
	})

	//login
	return userLoginConfirm(c, dbuser, req, rsp)
}

//just login
func userLoginConfirm(c *gin.Context, dbuser *model.DbUser, req *csmsg.CSMsgAccountLoginOrRegisterReq,
	rsp *csmsg.CSMsgAccountLoginOrRegisterRsp) cserr.ICSErrCodeError {
	if dbuser.RegisterInitTimeStamp == 0 {
		//todo here fill register info
		dbRegisterInfo := &ssmsg.TcsProfileDbUserRegisterInfoSave{}
		dbRegisterInfo.Uin = dbuser.Uin
		dbRegisterInfo.RegisterInfo = &comm.UserRegisterInfo{}
		dbRegister := dbRegisterInfo.RegisterInfo
		dbRegister.TimeStamp = time.Now().UnixNano()
		dbRegister.BundleId = req.LoginDevice.BundleId
		dbRegister.ClientVersion = int64(req.LoginDevice.ClientVersion)
		dbRegister.DeviceInfo = req.LoginDevice.Imei
		dbRegister.Imei = req.LoginDevice.Imei
		dbRegister.Ip = req.LoginDevice.Ip
		dbRegister.Language = req.LoginDevice.Language
		dbRegister.MacAddr = req.LoginDevice.MacAddr
		dbRegister.OsType = req.LoginDevice.OsType
		dbRegister.AuthType = req.LoginDevice.AuthType
		dbRegister.Scene = int32(req.AuthType)
		dbRegister.DeviceInfo = req.LoginDevice.DeviceInfo
		dbRegister.RealCountry = req.LoginDevice.RealCountry
		dbRegister.TimeZone = req.LoginDevice.TimeZone
		//grpc todo instead of
		task.TcsCall("profile", "DbUserRegisterInfoSave", dbRegisterInfo)
		////////////////////////////////////////////////////////////////////////
		if e := crud.DbUserUpdateFieldsByValueMap(dbuser.Uin,  sfw_db.UpdateFieldsMap{
			"RegisterInitTimeStamp" : time.Now().UnixNano(),
		}); !util.ErrOK(e) {
			return cserr.ErrDb
		}
	}

	//生成用户初始device 信息
	var device model.DbUserLoginDevice

	db,tb := crud.GetAccountDbAndTbNameByUin(dbuser.Uin, crud.AccountDbUserLoginDeviceTbNameBase)

	device.Uin = dbuser.Uin
	device.DeviceID = req.LoginDevice.DeviceId
	err := db.Table(tb).Where("uin = ? AND device_id = ?", dbuser.Uin, req.LoginDevice.DeviceId).
		FirstOrCreate(&device).Error
	if err != nil {
		log.Errorf(err, "create or get user tel:%s device:%s error", req.AuthId, req.LoginDevice.DeviceId)
		return cserr.ErrDb
	}

	device.LoginTimestamp = time.Now().UnixNano()
	device.IsOnline = true
	device.LoginClientVersion = int64(req.LoginDevice.ClientVersion)
	device.DeviceName = req.LoginDevice.DeviceName
	device.DeviceType = req.LoginDevice.DeviceType
	device.LoginAuthType = req.LoginDevice.AuthType
	device.LoginDeviceInfo = req.LoginDevice.DeviceInfo
	device.LoginIP = req.LoginDevice.Ip
	device.LoginRealCountry = req.LoginDevice.RealCountry
	device.BundleID = req.LoginDevice.BundleId
	device.Language = req.LoginDevice.Language
	device.TimeZone = req.LoginDevice.TimeZone
	device.DevicePushToken = req.LoginDevice.PushToken
	device.AutoAuthKey = util.RandomString(128)
	device.ClientDBEncryptKey = util.RandomString(36)
	device.ClientDBEncryptInfo = req.LoginDevice.DeviceId + "|" + device.ClientDBEncryptKey

	//设备
	util.CheckError(crud.DbUserLoginDeviceSave(&device), "save login device:%s error",device.DeviceID)

	tokenSecKey := util.RandomString(32)
	Uin := dbuser.Uin
	if e := crud.DbUserTokenUpdate(Uin, req.LoginDevice.DeviceId, tokenSecKey, constk.JWT_TOKEN_TIMEOUT_SECOND); e != nil {
		return cserr.ErrDb
	}
	//签发认证
	//sign the json web token
	t, e := token.Sign(c, token.AuthContext{
		UIN:           Uin,
		DeviceType:    req.LoginDevice.DeviceType,
		DeviceID:      req.LoginDevice.DeviceId,
		ClientVersion: req.LoginDevice.ClientVersion,
		TokenSecKey: tokenSecKey,
	},"")
	if e != nil {
		log.Errorf(e, "error token sign !")
		return cserr.ErrAccountTokenInvalid
	}

	//sign the sessionkey
	sessionKey, e := token.SignLongConnAuthToken(token.LongConnAuthContext{
		UIN:           Uin,
		DeviceType:    req.LoginDevice.DeviceType,
		DeviceID:      req.LoginDevice.DeviceId,
		ClientVersion: req.LoginDevice.ClientVersion,
	},"")
	if e != nil {
		log.Errorf(e, "error token sign !")
		return cserr.ErrAccountTokenInvalid
	}

	//sign the fileKey
	fileKey, e := token.SignFileToken(token.FileContext{
		UIN:           Uin,
		DeviceType:    req.LoginDevice.DeviceType,
		DeviceID:      req.LoginDevice.DeviceId,
		ClientVersion: req.LoginDevice.ClientVersion,
	},"")
	if e != nil {
		log.Errorf(e, "error token sign !")
		return cserr.ErrAccountTokenInvalid
	}

	//rsp.User = DbUserToCsUser(dbuser, &device)
	rsp.Uin = dbuser.Uin
	rsp.PasswordStat = &csmsg.CSUserPasswordStat{
		LastUpdateUsernameTimeStamp: dbuser.ModifyUserNameLastTimeStamp,
		RegisterTimeStamp: dbuser.RegisterInitTimeStamp,
		LastUpdatePasswdTimeStamp: dbuser.ModifyPasswdLastTimeStamp,
	}
	rsp.AuthInfo = &csmsg.CSLoginUserAuthInfo{}
	rsp.AuthInfo.Token = t
	rsp.AuthInfo.SessionKey = sessionKey
	rsp.AuthInfo.FileKey = fileKey
	rsp.AuthInfo.ClientDbEncKey = device.ClientDBEncryptKey
	rsp.AuthInfo.ClientDbEncInfo = device.ClientDBEncryptInfo
	rsp.AuthInfo.AutoAuthKey = device.AutoAuthKey
	rsp.ClientNewVersion = ""

	sfw_event.FireEvent(ssmsg.SvcEventType_SVC_EVT_ACCOUNT_LOGIN, &ssmsg.SvcEvtAccountLogin{
		Uin: Uin,
	})

	return nil
}

//注册发送验证码.如果账号已经存在则返回已存在错误.否则签发验证码.验证码确认后直接登录
func LoginOrRegister(c *gin.Context, req *csmsg.CSMsgAccountLoginOrRegisterReq,
	rsp *csmsg.CSMsgAccountLoginOrRegisterRsp) (serr cserr.ICSErrCodeError) {
	//获取或者创建用户
	if req.AuthType == comm.UserAuthType_USER_AUTH_SMS_CODE {
		code := req.Ticket
		tel := req.AuthId
		//

		phone := model.DbUserLoginPhone{
			Phone: req.AuthId,
		}
		db, tb := crud.GetAccountDbAndTbNameByTelephone(req.AuthId, crud.AccountDbUserLoginPhoneTbNameBase)
		//user is exist
		var userExist = true
		var captchType = constk.CAPTCHA_ACCOUNT_LOGIN
		if db.Table(tb).First(&phone, "phone = ?", tel).RecordNotFound() || phone.Uin == "" {
			captchType = constk.CAPTCHA_ACCOUNT_REGISTER
			userExist = false
		}

		if err := captcha.CheckCodeByTargetScene(tel, code, captchType); !util.ErrOK(err) {
			log.Errorf(err, "captcha check error ")
			serr = cserr.ErrCaptcha
			return
		}

		dbuser, err := crud.DbUserGetByTel(tel)
		if !util.ErrOK(err) && err != cserr.ErrAccountUserNotExist {
			serr = cserr.ErrDb
			return
		}
		if userExist == false || dbuser == nil {
			serr = userRegisterConfirm(c, req, rsp)
			return
		} else {
			serr = userLoginConfirm(c, dbuser, req, rsp)
			return
		}
	} else {
		//just login
		//using username password login (user name validator todo)
		//here, middleware should check session code in case of passwd guessing todo
		//login with username
		db, tb := crud.GetAccountDbAndTbNameByTelephone(req.AuthId, crud.AccountDbUserLoginNameTbNameBase)

		login_uname := model.DbUserLoginName{}
		err := db.Table(tb).First(&login_uname, "user_name = ?", req.AuthId).Error
		if !util.ErrOK(err)  {
			log.Errorf(err, "user name:%s login error", req.AuthId)
			serr = cserr.ErrAccountUserNotExist
			return
		}

		var verifyUserPass = false
		var device *model.DbUserLoginDevice = nil
		//using default pass login
		if len(req.AuthCode) > 0 {
			verifyUserPass = SafeCheckAuthPass(req.AuthKey, login_uname.Password)
		} else {
			//
			device, err = crud.DbUserLoginDeviceGet(login_uname.Uin, req.LoginDevice.DeviceId)
			if !util.ErrOK(err) || device == nil {
				serr = cserr.ErrAccountTokenInvalid
				return
			}
			//using auto pass login
			verifyUserPass = SafeCheckAuthPass(req.AuthKey, device.AutoAuthKey)
		}
		if verifyUserPass == false {
			log.Warnf("user name:%s check pass error !", req.AuthId)
			serr = cserr.ErrAccountTokenInvalid
			return
		}

		dbuser,uerr := crud.DbUserGetByUin(login_uname.Uin)
		if util.CheckError(uerr, "get dbuser by uin:%s fail", login_uname.Uin) {
			serr = cserr.ErrDb
			return
		}

		serr = userLoginConfirm(c, dbuser, req, rsp)
		return
	}
}


//authed success , invlaid user ticket
func Logout(c *gin.Context, req *csmsg.CSMsgAccountLogoutReq, rsp *csmsg.CSMsgAccountLogoutRsp) (err cserr.ICSErrCodeError) {
	//deivce exit ， invalid deice jwt secret
	dbUserTokenState, ok := c.Get("db.UserTokenState")
	if !ok || dbUserTokenState == nil {
		log.Errorf(nil, "db auth device error")
		return cserr.ErrAccountTokenInvalid
	}
	tokenState := dbUserTokenState.(*model.DbUserTokenState)
	//expired now
	tokenState.JwtExpiredTime = time.Now().Unix() - 1

	err = crud.DbUserTokenUpdate(tokenState.Uin, tokenState.DeviceID, tokenState.TokenSecKey, tokenState.JwtExpiredTime)
	if util.CheckError(err, "token state update error") {
		return cserr.ErrDb
	}
	sfw_event.FireEvent(ssmsg.SvcEventType_SVC_EVT_ACCOUNT_LOGOUT, &ssmsg.SvcEvtAccountLogout{
		Uin: tokenState.Uin,
		DeviceId: tokenState.DeviceID,
	})
	err = nil
	return
}
