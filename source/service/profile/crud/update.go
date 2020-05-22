package crud

import (
	"github.com/lexkong/log"
	"groot/proto/cserr"
	"groot/sfw/db"
)

func DbUserBaseInfoUpdateFieldByName(uin string, attrName string, val interface{}) cserr.ICSErrCodeError {
	db,tb := GetProfileDbAndTbNameByUin(uin, DbUserBaseInfoTbNameBase)
	e := db.Table(tb).Where("uin = ?", uin).Update(attrName, val).Error
	if e != nil {
		log.Errorf(e, "update fail")
		return cserr.ErrDbUpdateFail
	}
	return nil
}

func DbUserBaseInfoUpdateFieldsByValueMap(uin string, vals sfw_db.UpdateFieldsMap) cserr.ICSErrCodeError {
	db,tb := GetProfileDbAndTbNameByUin(uin, DbUserBaseInfoTbNameBase)
	e := db.Table(tb).Where("uin = ?", uin).Updates(map[string]interface{}(vals)).Error
	if e != nil {
		log.Errorf(e, "update fail")
		return cserr.ErrDbUpdateFail
	}
	return nil
}


func DbUserSettingUpdateFieldsByValueMap(uin string, vals sfw_db.UpdateFieldsMap) cserr.ICSErrCodeError {
	db,tb := GetProfileDbAndTbNameByUin(uin, DbUserSettingTbNameBase)
	e := db.Table(tb).Where("uin = ?", uin).Updates(map[string]interface{}(vals)).Error
	if e != nil {
		log.Errorf(e, "update fail")
		return cserr.ErrDbUpdateFail
	}
	return nil
}

