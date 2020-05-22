package crud

import (
	"github.com/lexkong/log"
	"groot/proto/cserr"
	"groot/service/account/model"
	sfw_db "groot/sfw/db"
	"groot/sfw/util"
	"time"
)

func DbUserLoginDeviceUpdate(uin,deviceId string, fieldsMap sfw_db.UpdateFieldsMap) {
	db, tb := GetAccountDbAndTbNameByUin(uin, AccountDbUserLoginDeviceTbNameBase)
	e := db.Table(tb).Where("uin = ?  AND device_id = ?", uin, deviceId).UpdateColumns(fieldsMap).Error
	if e != nil {
		log.Errorf(e, "db update fail")
	}
}


func DbUserTokenUpdate(uin,device_id string, tokenSecKey string, timeoutSecond int64) cserr.ICSErrCodeError {
	db, tbname := GetAccountDbAndTbNameByUin(uin, AccountDbUserTokenStateTbNameBase)
	tokenState := model.DbUserTokenState{}
	timeNow := time.Now().Unix()
	if db.Table(tbname).First(&tokenState, "uin = ? AND device_id = ?",  uin, device_id).RecordNotFound() {
		//not exist . create
		tokenState.Uin = uin
		tokenState.DeviceID = device_id
		tokenState.TokenSecKey = tokenSecKey
		tokenState.JwtExpiredTime = timeNow + timeoutSecond
		tokenState.JwtIssueTime = timeNow
		if err := db.Table(tbname).Create(&tokenState).Error; err != nil {
			log.Errorf(err,"create token state error !")
			return cserr.ErrDb
		}
		return nil
	} else {
		//update
		if err := db.Table(tbname).Where("uin = ? AND device_id = ?", uin,
			device_id).Updates(model.DbUserTokenState{
			TokenSecKey:    tokenSecKey,
			JwtExpiredTime: timeNow + timeoutSecond,
			JwtIssueTime:   timeNow,
		}).Error; err != nil {
			log.Errorf(err,"update token state error !")
			return cserr.ErrDb
		}
		return nil
	}
}


func DbUserUpdateFieldsByValueMap(uin string, fieldsMap sfw_db.UpdateFieldsMap) cserr.ICSErrCodeError {
	db,tbname,_ := GetAccountDbUserDbAndTbNameByUin(uin)
	err := db.Table(tbname).Where("uin = ?", uin).Updates(fieldsMap).Error
	if util.CheckError(err, "update db user uin:%s error", uin) {
		return cserr.ErrDbUpdateFail
	}
	return nil
}

func DbUserLoginDeviceSave(dbLoginDevice * model.DbUserLoginDevice) error {
	db,tbname := GetAccountDbAndTbNameByUin(dbLoginDevice.Uin, AccountDbUserLoginDeviceTbNameBase)
	err := db.Table(tbname).Save(dbLoginDevice).Error
	util.CheckError(err, "save db user login device uin:%s id:%s error",
		dbLoginDevice.Uin, dbLoginDevice.DeviceID)
	return err
}
