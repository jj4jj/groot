package crud

import (
	"github.com/lexkong/log"
	"groot/service/profile/model"
)

func DbUserRegisterInfoCreate(info *model.DbUserRegisterInfo){
	db,tb := GetProfileDbAndTbNameByUin(info.Uin, DbUserRegisterInfoTbNameBase)
	if e := db.Table(tb).Create(info).Error; e != nil {
		log.Errorf(e, "db create user:%s register fail", info.Uin)
		return
	}
}