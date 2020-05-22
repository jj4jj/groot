package model

import (
	"github.com/lexkong/log"
	"groot/proto/comm"
	"groot/proto/csmsg"
	"groot/proto/ssmsg"
)

func TcsUserRegisterToDbUserRegisterInfo(tcsUserRegister * ssmsg.TcsProfileDbUserRegisterInfoSave) *DbUserRegisterInfo {
	if tcsUserRegister == nil {return nil}
	if tcsUserRegister.RegisterInfo == nil {
		log.Errorf(nil, "tcs save uin:%s register info is nil", tcsUserRegister.Uin)
		return nil
	}
	return &DbUserRegisterInfo{
		Uin: 		   tcsUserRegister.Uin,
		Ip:            tcsUserRegister.RegisterInfo.Ip,
		Ticket:        tcsUserRegister.RegisterInfo.Ticket,
		MacAddr:       tcsUserRegister.RegisterInfo.MacAddr,
		Imei:          tcsUserRegister.RegisterInfo.Imei,
		Scene:         tcsUserRegister.RegisterInfo.Scene,
		OsType:        tcsUserRegister.RegisterInfo.OsType,
		DeviceInfo:    tcsUserRegister.RegisterInfo.DeviceInfo,
		AuthType:      tcsUserRegister.RegisterInfo.AuthType,
		TimeZone:      tcsUserRegister.RegisterInfo.TimeZone,
		Language:      tcsUserRegister.RegisterInfo.Language,
		RealCountry:   tcsUserRegister.RegisterInfo.RealCountry,
		ClientVersion: tcsUserRegister.RegisterInfo.ClientVersion,
		BundleId:      tcsUserRegister.RegisterInfo.BundleId,
		TimeStamp:     tcsUserRegister.RegisterInfo.TimeStamp,
	}
}

func DbUserBaseToCsUserBase(dbUserBase *DbUserBaseInfo) *comm.UserBase {
	if dbUserBase == nil {
		return nil
	}
	return &comm.UserBase{
		Uin:      dbUserBase.Uin,
		NickName: dbUserBase.NickName,
		Sex:      dbUserBase.SexType,
		Location: &comm.Location{
			Desc:       dbUserBase.LocationDesc,
			Longtitude: dbUserBase.Longitude,
			Latitude:   dbUserBase.Latitude,
		},
		IconUrl: dbUserBase.HeadIcon,
	}
}

func DbUserRegisterToCsUserRegisterInfo(dbUserRegister *DbUserRegisterInfo) *csmsg.CSUserRegisterInfo {
	if dbUserRegister == nil {
		return nil
	}
	return &csmsg.CSUserRegisterInfo{
		RegisterInfo: &comm.UserRegisterInfo{
			Ip:            dbUserRegister.Ip,
			Ticket:        dbUserRegister.Ticket,
			MacAddr:       dbUserRegister.MacAddr,
			Imei:          dbUserRegister.Imei,
			Scene:         dbUserRegister.Scene,
			OsType:        dbUserRegister.OsType,
			DeviceInfo:    dbUserRegister.DeviceInfo,
			AuthType:      dbUserRegister.AuthType,
			TimeZone:      dbUserRegister.TimeZone,
			Language:      dbUserRegister.Language,
			RealCountry:   dbUserRegister.RealCountry,
			ClientVersion: dbUserRegister.ClientVersion,
			BundleId:      dbUserRegister.BundleId,
			TimeStamp:     dbUserRegister.TimeStamp,
		},
	}
}

func DbUserProfileToCsProfileInfo(dbProfile *DbUserProfileInfo) *csmsg.CSUserProfileInfo {
	if dbProfile == nil {
		return nil
	}
	return &csmsg.CSUserProfileInfo{
		College:   dbProfile.College,
		Job:       dbProfile.Job,
		Education: dbProfile.Education,
	}
}


func DbUserSettingToCsUserSetting(dbSetting *DbUserSetting) * csmsg.CSUserSetting {
	if dbSetting == nil {
		return nil
	}
	return &csmsg.CSUserSetting{
		JsonSetting: dbSetting.GetJsonUserSetting(),
		EnableSearchByTel: dbSetting.EnableSearchByTel == 1,
		CustomSetting: dbSetting.CustomSetting,
	}
}
