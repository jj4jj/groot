package account

import (
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"groot/proto/cserr"
	"groot/proto/csmsg"
	"groot/service/account/ctx"
	"groot/service/account/crud"
	"groot/service/account/model"
	"groot/sfw/crypto/token"
	"groot/sfw/db"
	"groot/sfw/util"
	"time"
)

func TokenRenew(c * gin.Context, req *csmsg.CSMsgAccountTokenRenewReq, rsp *csmsg.CSMsgAccountTokenRenewRsp) (err cserr.ICSErrCodeError) {
	uin := ctx.NeedAuthContextGetUin(c)
	log.Debugf("user:%v token ticket:%s renew", uin, req.Ticket)
	tokenState := ctx.NeedAuthContextGetUserTokenState(c)
	//签发认证


	device, err := crud.DbUserLoginDeviceGet(uin, tokenState.DeviceID)
	if util.CheckError(err, "get device error !") {
		return cserr.ErrDb
	}

	//sign the json web token
	tokenSecKey := util.RandomString(32)
	t, e := token.Sign(c, token.AuthContext{
		UIN:           uin,
		DeviceType:    device.DeviceType,
		DeviceID:      device.DeviceID,
		ClientVersion: uint32(device.LoginClientVersion),
		TokenSecKey:  tokenSecKey,
	},"")
	if e != nil {
		log.Errorf(e, "error token sign !")
		return cserr.ErrInternal
	}
	rsp.Token = t
	var renewTimeSeconds int64 = 3600
	rsp.ExpiredTimeStamp = time.Now().Unix() + renewTimeSeconds
	return crud.DbUserTokenUpdate(tokenState.Uin, tokenState.DeviceID, tokenSecKey, renewTimeSeconds)

}

func ChangePassword(c * gin.Context, req *csmsg.CSMsgAccountChangePasswordReq, rsp *csmsg.CSMsgAccountChangePasswordRsp) (err cserr.ICSErrCodeError) {
	uin := ctx.NeedAuthContextGetUin(c)

	//todo add algorithm
	if req.Passwd != req.PasswdConfirm {
		return cserr.ErrAccountWrongPasswd
	}
	//
	dbUser,e := crud.DbUserGetByUin(uin)
	if util.CheckError(e, "get user by uin:%s", uin) {
		return e
	}
	db := crud.GetAccountDbByKey([]byte(dbUser.LoginUserName))
	tbname := crud.GetAccountTbNameByKey([]byte(dbUser.LoginUserName), crud.AccountDbUserLoginNameTbNameBase)

	dbUserLoginName := model.DbUserLoginName{}
	if err := db.Table(tbname).First(&dbUserLoginName, "user_name = %s AND uin = %s",
		dbUser.LoginUserName, uin).Error; err != nil {
		log.Errorf(err, "db find user:%s uin:%s fail", dbUser.LoginUserName, uin)
		return cserr.ErrDb
	}

	//update passwd
	if dbUser.ModifyPasswdLastTimeStamp != 0 {
		log.Warnf("change password last time is not 0")
		return cserr.ErrBadRequest
	}
	dbUserLoginName.Password = req.PasswdConfirm
	if e:= db.Table(tbname).Where("user_name = %s AND uin = %s", dbUser.LoginUserName,
		uin).Update("Password", req.PasswdConfirm).Error; e!= nil {
		log.Errorf(e, "change password db fail")
		return cserr.ErrDb
	}

	util.CheckError(crud.DbUserUpdateFieldsByValueMap(uin, sfw_db.UpdateFieldsMap{
		"ModifyPasswdLastTimeStamp": time.Now().Unix(),
	}),"update user:%s passwd modify time", uin)
	return nil
}



func ChangeTelephoneOldTelephoneGetSmsCode(c * gin.Context, req *csmsg.CSMsgAccountChangeTelephoneOldTelephoneGetSmsCodeReq,
	rsp *csmsg.CSMsgAccountChangeTelephoneOldTelephoneGetSmsCodeRsp) (err cserr.ICSErrCodeError) {
	return cserr.ErrNoImplement
}


func ChangeTelephoneNewTelephoneGetSmsCode(c * gin.Context, req *csmsg.CSMsgAccountChangeTelephoneNewTelephoneGetSmsCodeReq,
	rsp *csmsg.CSMsgAccountChangeTelephoneNewTelephoneGetSmsCodeRsp) (err cserr.ICSErrCodeError) {
	return cserr.ErrNoImplement
}



func ChangeTelephoneConfirm(c * gin.Context, req *csmsg.CSMsgAccountChangeTelephoneConfirmReq,
	rsp *csmsg.CSMsgAccountChangeTelephoneConfirmRsp) (err cserr.ICSErrCodeError) {
	return cserr.ErrNoImplement
}

