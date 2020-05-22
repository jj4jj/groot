package profile

import (
	"github.com/lexkong/log"
	"groot/proto/ssmsg"
	"groot/service/profile/crud"
	"groot/service/profile/model"
)

func TcsUserRegisterInfoSave(ssinfo * ssmsg.TcsProfileDbUserRegisterInfoSave) {
	info := model.TcsUserRegisterToDbUserRegisterInfo(ssinfo)
	if info == nil || info.Uin == ""{
		log.Errorf(nil, "tcs user register to db model fail !")
		return
	}
	crud.DbUserRegisterInfoCreate(info)
	return
}