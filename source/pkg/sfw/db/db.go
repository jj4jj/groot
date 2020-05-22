package sfw_db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/lexkong/log"
	"groot/comm/conf"

	//_ "github.com/mattn/go-sqlite3"
)

var (
	dbBackends = make(map[string]*gorm.DB)
	defaultDb  *gorm.DB
	dbUseMap = make(map[string]string)
)

func GetWithUseMap(name string) *gorm.DB {
	connx,ok := dbUseMap[name]
	if ok  {
		return dbBackends[connx]
	} else {
		return nil
	}
}

func GetDb(name string) *gorm.DB {
	config := conf.GetAppConfig()
	if config.RunEnv == "dev" {
		return defaultDb
	}
	db, ok := dbBackends[name]
	if ok {
		return db
	} else {
		if mdb := GetWithUseMap(name); mdb != nil  {
			return mdb
		}
		log.Warnf("get db name:%s is not exist", name)
	}
	return nil
}

func AddDbUseMap(dbName,connxDbName string) {
	dbUseMap[dbName] = connxDbName
}

func UpdateDbUseMap(mp map[string]string) {
	for dbName,connxName := range mp {
		dbUseMap[dbName] = connxName
	}
}

func GetDefaultDb() *gorm.DB {
	return defaultDb
}

//初始化DB连接
func AddDb(name string, driver string, dbenv string, maxIdleCx int, maxOpenCx int) error {
	db, err := gorm.Open(driver, dbenv) // "<user>:<password>/<database>?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return err
	}
	db.DB().SetMaxIdleConns(maxIdleCx)
	db.DB().SetMaxOpenConns(maxOpenCx)
	dbBackends[name] = db
	if defaultDb == nil {
		defaultDb = db
	}

	if conf.GetAppConfig().RunEnv == "dev" {
		db.LogMode(true)
	}

	return nil
}

func CloseAllDb() {
	for _, db := range dbBackends {
		db.Close()
	}
}
