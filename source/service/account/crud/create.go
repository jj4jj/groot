package crud

import (
	"github.com/lexkong/log"
	"groot/comm/conf"
	"groot/proto/comm"
	"groot/proto/cserr"
	"groot/service/account/model"
	sfw_db "groot/sfw/db"
	"groot/sfw/util"
)

func DbUserCreateByTelephone(tel string) (* model.DbUser, error) {

	config := conf.GetAppConfig()

	dbnum := config.Account.UseDb.SplitDbNum
	tbnum := config.Account.UseDb.SplitTbNum

	dbKeyCode := sfw_db.GetDbKeyHashCode([]byte(tel))
	dbidx := (dbKeyCode /tbnum) % dbnum
	tbidx := dbKeyCode % tbnum

	dbname := sfw_db.GetDbNameByIdx(AccountDbNameBase, dbidx)
	dbUserTbName := sfw_db.GetTbNameByIdx(AccountDbUserTbNameBase, dbidx, tbidx)
	db := sfw_db.GetDb(dbname)

	loginPhoneTbname := sfw_db.GetTbNameByIdx(AccountDbUserLoginPhoneTbNameBase, dbidx, tbidx)
	dbPhone := model.DbUserLoginPhone {
		Phone: tel,
	}
	//get or create phone
	err := db.Table(loginPhoneTbname).Where("phone = ?", tel).FirstOrCreate(&dbPhone).Error
	if util.CheckError(err, "create login phone : %s error", tel) {
		return nil,cserr.ErrDb
	}

	//gen user name
	userName := model.RandomUserNameByTelephone(tel)
	//create user
	user := model.DbUser{
		LoginTel: tel,
		TicketSecret: util.RandomString(16),
		UserType: comm.UserType_USER_TYPE_NORMAL,
	}
	if util.CheckError(db.Table(dbUserTbName).Create(&user).Error,"create user fail by tel:%s", tel) {
		return nil, cserr.ErrDb
	}

	//get uin
	uin := model.GetUinFromDbUserRoute(user.ID, dbidx, tbidx)
	user.Uin = uin

	//update uin for phone
	if util.CheckError(db.Table(loginPhoneTbname).Where("phone = ?",
		tel).Update("uin", uin).Error, "update phone uin:%s fail", uin) {
		return nil,cserr.ErrDb
	}

	dbUserLoginUserName := model.DbUserLoginName{
		UserName: userName,
		Uin: uin,
		Password: util.RandomString(8),
	}
	userNameDb := GetAccountDbByKey([]byte(userName))
	userNameTb := GetAccountTbNameByKey([]byte(userName), AccountDbUserLoginNameTbNameBase)
	err = userNameDb.Table(userNameTb).Create(&dbUserLoginUserName).Error
	if util.CheckError(err, "create user login pass by uin:%s error", uin) {
		return nil,cserr.ErrDb
	}

	//update uin,username dbuser
	if util.CheckError(db.Table(dbUserTbName).Where("id = ?",
		user.ID).Updates(model.DbUser{Uin:uin,LoginUserName:userName}).Error, "update user uin:%s fail", uin) {

		log.Errorf(nil, "update db user id:%d uin:%s user name:%s fail", user.ID, uin, userName)
		return nil,cserr.ErrDb
	}


	log.Infof("create user:%s by tel:%s success", userName, tel)
	return &user, nil
}
