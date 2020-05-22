package crud

import (
	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"groot/service/profile/model"
)

func DbUserBaseInfoGetByUin(uin string) *model.DbUserBaseInfo {
	dbInfo := model.DbUserBaseInfo{}
	db, tb := GetProfileDbAndTbNameByUin(uin, DbUserBaseInfoTbNameBase)
	if e := db.Table(tb).First(dbInfo, "uin = ?", uin).Error; e != nil {
		if e == gorm.ErrRecordNotFound {
			return nil
		}
		log.Errorf(e, "db get base user:%s fail", uin)
		return nil
	}
	return &dbInfo
}

func DbUserRegisterInfoGetByUin(uin string) *model.DbUserRegisterInfo {
	dbInfo := model.DbUserRegisterInfo{}
	db, tb := GetProfileDbAndTbNameByUin(uin, DbUserRegisterInfoTbNameBase)
	if e := db.Table(tb).First(dbInfo, "uin = ?", uin).Error; e != nil {
		if e == gorm.ErrRecordNotFound {
			return nil
		}
		log.Errorf(e, "db get register user:%s fail", uin)
		return nil
	}
	return &dbInfo
}

func DbUserSettingGetByUin(uin string) *model.DbUserSetting {
	dbInfo := model.DbUserSetting{}
	db, tb := GetProfileDbAndTbNameByUin(uin, DbUserSettingTbNameBase)
	if e := db.Table(tb).First(dbInfo, "uin = ?", uin).Error; e != nil {
		if e == gorm.ErrRecordNotFound {
			return nil
		}
		log.Errorf(e, "db get setting user:%s fail", uin)
		return nil
	}
	return &dbInfo
}

func DbUserProfileInfoGetByUin(uin string) *model.DbUserProfileInfo {
	dbInfo := model.DbUserProfileInfo{}
	db, tb := GetProfileDbAndTbNameByUin(uin, DbUserProfileInfoTbNameBase)
	if e := db.Table(tb).First(dbInfo, "uin = ?", uin).Error; e != nil {
		if e == gorm.ErrRecordNotFound {
			return nil
		}
		log.Errorf(e, "db get profile user:%s fail", uin)
		return nil
	}
	return &dbInfo
}

