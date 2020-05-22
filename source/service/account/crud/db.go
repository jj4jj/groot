package crud

import (
	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"groot/comm/conf"
	"groot/proto/comm"
	"groot/service/account/model"
	"groot/sfw/db"
	sdb "groot/sfw/db"
)


var (
	AccountDbNameBase                  string
	AccountDbUserLoginDeviceTbNameBase string
	AccountDbUserLoginEmailTbNameBase  string
	AccountDbUserTbNameBase            string
	AccountDbUserLoginNameTbNameBase   string
	AccountDbUserLoginPhoneTbNameBase  string
	AccountDbUserTokenStateTbNameBase  string
)


func init(){
	AccountDbNameBase = "db_account"
	AccountDbUserLoginDeviceTbNameBase = sfw_db.GetModelTableNameBase(&model.DbUserLoginDevice{}) //"db_user_login_device"
	AccountDbUserLoginEmailTbNameBase = sfw_db.GetModelTableNameBase(&model.DbUserLoginEmail{})
	AccountDbUserTbNameBase = sfw_db.GetModelTableNameBase(&model.DbUser{})
	AccountDbUserLoginNameTbNameBase = sfw_db.GetModelTableNameBase(&model.DbUserLoginName{})
	AccountDbUserLoginPhoneTbNameBase = sfw_db.GetModelTableNameBase(&model.DbUserLoginPhone{})
	AccountDbUserTokenStateTbNameBase = sfw_db.GetModelTableNameBase(&model.DbUserTokenState{})
}


func GetAccountDbByKey(key []byte) *gorm.DB {
	config := conf.GetAppConfig()
	dbname := sdb.GetDbNameByKey(key, config.Account.UseDb.SplitDbNum, AccountDbNameBase)
	return sdb.GetDb(dbname)
}


func GetAccountTbNameByKey(key []byte, tbname_base string) string {
	config := conf.GetAppConfig()
	return sdb.GetTbNameByKey(key, config.Account.UseDb.SplitTbNum, config.Account.UseDb.SplitDbNum, tbname_base)
}


func GetAccountDbAndTbNameByTelephone(tel,tbname_base string) (*gorm.DB, string) {
	btel := []byte(tel)
	db := GetAccountDbByKey(btel)
	tbname := GetAccountTbNameByKey(btel, tbname_base)
	return db, tbname
}


func GetAccountDbAndTbNameByKey(key,tbname_base string) (*gorm.DB, string) {
	bk := []byte(key)
	db := GetAccountDbByKey(bk)
	tbname := GetAccountTbNameByKey(bk, tbname_base)
	return db, tbname
}

func GetAccountDbUserDbAndTbNameByUin(uin string) (db *gorm.DB, tbname string, db_user_id uint64) {
	dbid,dbidx,tbidx := model.GetDbUserRouteFromUin(uin)
	dbname := sfw_db.GetDbNameByIdx(AccountDbNameBase, dbidx)
	tbname = sfw_db.GetTbNameByIdx(AccountDbUserTbNameBase, dbidx, tbidx)
	db = sfw_db.GetDb(dbname)
	db_user_id = dbid
	return
}

func GetAccountDbAndTbNameByUin(uin,tbname_base string) (db *gorm.DB, tbname string) {
	_,dbidx,tbidx := model.GetDbUserRouteFromUin(uin)
	dbname := sfw_db.GetDbNameByIdx(AccountDbNameBase, dbidx)
	tbname = sfw_db.GetTbNameByIdx(tbname_base, dbidx, tbidx)
	db = sfw_db.GetDb(dbname)
	return
}

func GetSystemUin() string {
	uin := model.GetUinFromDbUserRoute(10000,0,0)
	return uin
}

func CreateSystemDbUser() error {
	uin := GetSystemUin()
	db,tbname,db_user_id := GetAccountDbUserDbAndTbNameByUin(uin)
	dbUser := model.DbUser{
		ID: db_user_id,
		Uin: uin,
		LoginTel: "system",
		LoginUserName: "system",
		LoginEmail: "system",
		UserType: comm.UserType_USER_TYPE_SYSTEM,
	}
	if e := db.Table(tbname).FirstOrCreate(&dbUser, "id = ? AND uin = ?", db_user_id, uin).Error; e != nil {
		log.Errorf(e, "get system db uin error !")
		return e
	}
	return nil
}

func InitAccountDbCtx() error {
	config := conf.GetAppConfig()
	sfw_db.UpdateDbUseMap(config.Account.UseDb.DbUseMap)
	return nil
}









