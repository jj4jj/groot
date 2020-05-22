package crud

import (
	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"groot/proto/comm"
	"groot/proto/cserr"
	"groot/service/account/model"
	"groot/sfw/util"

)



func DbUserLoginDeviceGet(uin,deviceId string) (device *model.DbUserLoginDevice, err cserr.ICSErrCodeError) {
	db, tbname := GetAccountDbAndTbNameByUin(uin, AccountDbUserLoginDeviceTbNameBase)
	device = &model.DbUserLoginDevice{}
	if e := db.Table(tbname).First(&device, "uin = ? AND device_id = ?", uin, deviceId).Error; err != nil {
		if e == gorm.ErrRecordNotFound {
			err = cserr.ErrDbNotFound
			return
		}
		log.Errorf(err, "find user device error !")
		err = cserr.ErrDb
		return
	}
	err = nil
	return
}


func GetUinByLoginPhone(tel string) (string, error) {
	db,tb := GetAccountDbAndTbNameByTelephone(tel, AccountDbUserLoginPhoneTbNameBase)
	telphone := model.DbUserLoginPhone {
		Phone: tel,
	}
	if err := db.Table(tb).First(&telphone).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", cserr.ErrAccountUserNotExist
		}
		log.Errorf(err, "find tel user db fail with tel:%s", tel)
		return "", cserr.ErrDb
	}
	return telphone.Uin,nil
}

func GetUinByLoginUserName(uname string) (string, error) {
	db,tb := GetAccountDbAndTbNameByKey(uname, AccountDbUserLoginNameTbNameBase)
	un := model.DbUserLoginName {
		UserName: uname,
	}
	if err := db.Table(tb).First(&un, " user_name = ?", uname).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", cserr.ErrAccountUserNotExist
		}
		log.Errorf(err, "find tel user db fail with uname:%s", uname)
		return "", cserr.ErrDb
	}
	return un.Uin,nil
}

func GetUinByLoginEmail(email string) (string, error) {
	db,tb := GetAccountDbAndTbNameByKey(email, AccountDbUserLoginEmailTbNameBase)
	un := model.DbUserLoginEmail {
		Email: email,
	}
	if err := db.Table(tb).First(&un, " email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", cserr.ErrAccountUserNotExist
		}
		log.Errorf(err, "find tel user db fail with email:%s", email)
		return "", cserr.ErrDb
	}
	return un.Uin,nil
}


func DbUserLoginDeviceGetList(uin string, osType comm.DeviceOsType) (list []model.DbUserLoginDevice) {
	db, tb := GetAccountDbAndTbNameByUin(uin, AccountDbUserLoginDeviceTbNameBase)
	if osType == comm.DeviceOsType_OS_TYPE_UNKNOWN {
		if e := db.Table(tb).Find(&list, "uin = ?", uin).Error; e != nil {
			log.Errorf(e, "find login device error")
			return nil
		}
	} else {
		if e := db.Table(tb).Find(&list, "uin = ? AND os_type = ?", uin, osType).Error; e != nil {
			log.Errorf(e, "find login device error")
			return nil
		}
	}
	return
}

func GetAccountTicketSecretByUin(Uin string) string {
	dbUser,e := DbUserGetByUin(Uin)
	if util.CheckError(e, "db user:%s get fail", Uin) {
		return ""
	}
	return dbUser.Uin

}



func DbUserGetByUserName(uname string) (*model.DbUser, cserr.ICSErrCodeError) {
	db := GetAccountDbByKey([]byte(uname))
	tbname := GetAccountTbNameByKey([]byte(uname), AccountDbUserLoginNameTbNameBase)
	dbUserName := model.DbUserLoginName{}
	e := db.Table(tbname).Where("user_name = ?", uname).First(&dbUserName).Error
	if e == gorm.ErrRecordNotFound {
		return nil,nil
	} else {
		if e != nil {
			log.Errorf(e, "db first error uname:%s", uname)
			return nil,cserr.ErrDb
		}
	}
	uin := dbUserName.Uin
	return DbUserGetByUin(uin)
}


func DbUserGetByTel(tel string) (*model.DbUser, cserr.ICSErrCodeError) {
	db := GetAccountDbByKey([]byte(tel))
	tbname := GetAccountTbNameByKey([]byte(tel), AccountDbUserLoginPhoneTbNameBase)
	dbUserPhone := model.DbUserLoginPhone{
	}
	e := db.Table(tbname).Where("phone = ?", tel).First(&dbUserPhone).Error
	if e == gorm.ErrRecordNotFound {
		return nil,nil
	} else {
		if e != nil {
			log.Errorf(e, "db first error tel:%s", tel)
			return nil,cserr.ErrDb
		}
	}
	uin := dbUserPhone.Uin
	return DbUserGetByUin(uin)
}


func DbUserGetByUin(uin string) (*model.DbUser, cserr.ICSErrCodeError) {
	dbu := model.DbUser{}
	db,tb,tbid := GetAccountDbUserDbAndTbNameByUin(uin)
	if err := db.Table(tb).First(&dbu, "id = ?", tbid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, cserr.ErrAccountUserNotExist
		}
		log.Errorf(err, "find user db fail with uin:%s", uin)
		return nil, cserr.ErrDb
	}
	return &dbu, nil
}





func DbUserTokenGet(uin, device_id string) *model.DbUserTokenState {
	db, tbname := GetAccountDbAndTbNameByUin(uin, AccountDbUserTokenStateTbNameBase)
	tokenState := model.DbUserTokenState{}
	if err := db.Table(tbname).First(&tokenState,"uin = ? AND device_id = ?", uin, device_id).Error; err != nil {
		log.Errorf(err,"token state get error ")
		return nil
	}
	return &tokenState
}


