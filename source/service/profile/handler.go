package profile

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"groot/proto/comm"
	"groot/proto/cserr"
	"groot/proto/csmsg"
	"groot/service/account"
	"groot/service/account/crud"
	"groot/service/account/ctx"
	account_model "groot/service/account/model"
	crud2 "groot/service/profile/crud"
	"groot/service/profile/model"
	"groot/sfw/db"
	"groot/sfw/util"
)

func ProfileSetSex(c * gin.Context, req *csmsg.CSMsgProfileSetSexReq, rsp *csmsg.CSMsgProfileSetSexRsp)(err cserr.ICSErrCodeError) {
	err = cserr.ErrInternal
	uin := ctx.NeedAuthContextGetUin(c)
	err = crud2.DbUserBaseInfoUpdateFieldByName(uin, "SexType", req.Sex)
	return
}

func ProfileSetDesc(c * gin.Context, req *csmsg.CSMsgProfileSetDescReq, rsp *csmsg.CSMsgProfileSetDescRsp)(err cserr.ICSErrCodeError) {
	err = cserr.ErrInternal
	uin := ctx.NeedAuthContextGetUin(c)
	err = crud2.DbUserBaseInfoUpdateFieldByName(uin, "Desc", req.Desc)
	return
}

func ProfileSetIcon(c * gin.Context, req *csmsg.CSMsgProfileSetIconReq, rsp *csmsg.CSMsgProfileSetIconRsp)(err cserr.ICSErrCodeError) {
	err = cserr.ErrInternal
	//todo check url head
	uin := ctx.NeedAuthContextGetUin(c)
	err = crud2.DbUserBaseInfoUpdateFieldByName(uin, "HeadIcon", req.IconUrl)
	return

}

func ProfileSetNickName(c * gin.Context, req *csmsg.CSMsgProfileSetNickNameReq, rsp *csmsg.CSMsgProfileSetNickNameRsp)(err cserr.ICSErrCodeError) {
	err = cserr.ErrInternal
	//todo check url head
	uin := ctx.NeedAuthContextGetUin(c)
	err = crud2.DbUserBaseInfoUpdateFieldByName(uin, "NickName", req.NickName)
	return

}


func ProfileSetLocation(c * gin.Context, req *csmsg.CSMsgProfileSetLocationReq, rsp *csmsg.CSMsgProfileSetLocationRsp)(err cserr.ICSErrCodeError) {
	err = cserr.ErrInternal
	uin := ctx.NeedAuthContextGetUin(c)
	err = crud2.DbUserBaseInfoUpdateFieldsByValueMap(uin, sfw_db.UpdateFieldsMap{
		"LocationDesc": req.Location.Desc,
		"Longitude":req.Location.Longtitude,
		"Latitude":req.Location.Latitude,
	})
	return
}

func ProfileUpdateSetting(c * gin.Context, req *csmsg.CSMsgProfileUpdateSettingReq, rsp *csmsg.CSMsgProfileUpdateSettingRsp)(err cserr.ICSErrCodeError) {
	err = cserr.ErrInternal
	uin := ctx.NeedAuthContextGetUin(c)
	kv := sfw_db.UpdateFieldsMap{}
	if len(req.CustomSetting) > 0  {
		kv["CustomSetting"] = req.CustomSetting
	}
	if req.JsonSetting != nil {
		jb,e := json.Marshal(&req.JsonSetting)
		if e != nil {
			kv["JsonSetting"] = jb
		}
	}
	if len(kv) > 0 {
		err = crud2.DbUserSettingUpdateFieldsByValueMap(uin, kv)
	} else {
		err = cserr.ErrBadRequest
	}

	return

}

func ProfileSetEnableSearchByTel(c * gin.Context, req *csmsg.CSMsgProfileSetEnableSearchByTelReq,
	rsp *csmsg.CSMsgProfileSetEnableSearchByTelRsp)(err cserr.ICSErrCodeError) {
	err = cserr.ErrInternal
	uin := ctx.NeedAuthContextGetUin(c)
	enableVal := 0
	if req.Enable {
		enableVal = 1
	}
	err = crud2.DbUserSettingUpdateFieldsByValueMap(uin, sfw_db.UpdateFieldsMap{
		"EnableSearchByTel": enableVal,
	})
	return
}

func GetCsUserInfoByDbUser(user * account_model.DbUser) *csmsg.CSUser {
	if user == nil {
		return nil
	}
	dbUserBase := crud2.DbUserBaseInfoGetByUin(user.Uin)
	if dbUserBase == nil {
		return nil
	}
	csUser := &csmsg.CSUser{}
	csUser.Base = model.DbUserBaseToCsUserBase(dbUserBase)
	//fix user type
	csUser.Base.UserName = user.LoginUserName
	csUser.Base.UserType = user.UserType

	dbUserRegister := crud2.DbUserRegisterInfoGetByUin(user.Uin)
	csUser.RegisterInfo = model.DbUserRegisterToCsUserRegisterInfo(dbUserRegister)

	dbUserProfile := crud2.DbUserProfileInfoGetByUin(user.Uin)
	csUser.Profile = model.DbUserProfileToCsProfileInfo(dbUserProfile)

	return csUser
}

func GetSelfUserInfo(c * gin.Context, req *csmsg.CSMsgProfileGetSelfUserInfoReq, rsp *csmsg.CSMsgProfileGetSelfUserInfoRsp)(err cserr.ICSErrCodeError) {
	err = cserr.ErrInternal
	dbUser := ctx.NeedUserContextGetDbUser(c)
	if dbUser == nil {
		log.Errorf(nil, "need user context not found")
		err = cserr.ErrBadRequest
		return
	}
	rsp.UserInfo = GetCsUserInfoByDbUser(dbUser)
	if rsp.UserInfo == nil {
		log.Errorf(nil, "get cs user info fail by uin:%s", dbUser.Uin)
		err = cserr.ErrInternal
	}
	dbSetting := crud2.DbUserSettingGetByUin(dbUser.Uin)
	if dbSetting != nil {
		rsp.Setting = model.DbUserSettingToCsUserSetting(dbSetting)
	}

	err = nil
	return
}

func GetUserBase(c * gin.Context, req *csmsg.CSMsgProfileGetUserBaseReq, rsp *csmsg.CSMsgProfileGetUserBaseRsp)(err cserr.ICSErrCodeError) {
	err = cserr.ErrInternal
	uin := ctx.NeedAuthContextGetUin(c)
	for i := range req.UinTicketList {
		ret := account.VerifyUserCheckTicketWithUin(uin, req.UinTicketList[i])
		if ret == nil {
			log.Warnf("get user ticket check fail")
			continue
		}
		if ret.UinSigner != uin {
			log.Warnf("user check ticket (%s.%s) fail", ret.UinSigner, uin)
			continue
		}
		if ret.Source == comm.UserSourceType_USER_SOURCE_SEARCH_TEL {
			info := crud2.DbUserSettingGetByUin(ret.UinCheck)
			if info != nil && info.EnableSearchByTel == 0 {
				log.Debugf("user disable search by tel")
				continue
			}
		}
		baseInfo := crud2.DbUserBaseInfoGetByUin(ret.UinCheck)
		if baseInfo != nil {
			userBase := model.DbUserBaseToCsUserBase(baseInfo)
			dbUserCheck, e := crud.DbUserGetByUin(baseInfo.Uin)
			if util.ErrOK(e) && dbUserCheck != nil {
				userBase.UserType = dbUserCheck.UserType
				userBase.UserName = dbUserCheck.LoginUserName
			}
			rsp.UserBaseList = append(rsp.UserBaseList, userBase)
		} else {
			log.Errorf(nil, "get user:%s base fail", ret.UinCheck)
		}
	}
	err = nil
	return
}






