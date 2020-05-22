package crud

import (
	"github.com/jinzhu/gorm"
	"groot/comm/conf"
	"groot/service/profile/model"
	"groot/sfw/db"
)

var (
	ProfileDbNameBase string
	DbUserBaseInfoTbNameBase 	string
	DbUserSettingTbNameBase 	string
	DbUserProfileInfoTbNameBase	 string
	DbUserRegisterInfoTbNameBase	 string
)

func init(){
	ProfileDbNameBase = "db_profile"
	DbUserBaseInfoTbNameBase = sfw_db.GetModelTableNameBase(model.DbUserBaseInfo{})
	DbUserSettingTbNameBase = sfw_db.GetModelTableNameBase(model.DbUserSetting{})
	DbUserProfileInfoTbNameBase = sfw_db.GetModelTableNameBase(model.DbUserProfileInfo{})
	DbUserRegisterInfoTbNameBase = sfw_db.GetModelTableNameBase(model.DbUserRegisterInfo{})
}

func InitProfileDbCtx() error {
	return nil
}

func GetProfileDbAndTbNameByUin(uin,tbNameBase string)(*gorm.DB, string){
	config := conf.GetAppConfig()
	b := []byte(uin)
	dbname := sfw_db.GetDbNameByKey(b, config.Profile.UseDb.SplitDbNum, ProfileDbNameBase)
	tbname := sfw_db.GetTbNameByKey(b, config.Profile.UseDb.SplitDbNum, config.Profile.UseDb.SplitTbNum, tbNameBase)
	return sfw_db.GetDb(dbname), tbname
}



