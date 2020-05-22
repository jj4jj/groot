package crud

import (
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"groot/comm/conf"
	"groot/proto/comm"
	"groot/proto/cserr"
	"groot/service/notify/model"
	"groot/sfw/db"
	"groot/sfw/util"
)

var (
	NotifyDbTbNameBase          string
	DbNotifyUserEventTbNameBase string
)


func init(){
	NotifyDbTbNameBase = "db_notify"
	DbNotifyUserEventTbNameBase = sfw_db.GetModelTableNameBase(&model.DbNotifyUserEvent{})
}

func DbNotifyGetUserEventByDbId(uin string, dbEventId uint64)(*model.DbNotifyUserEvent, cserr.ICSErrCodeError) {
	db, tb := GetNotifyDbAndTbNameByKey(uin, DbNotifyUserEventTbNameBase)
	de := model.DbNotifyUserEvent{}
	e := db.Table(tb).First(&de, "uin = ? AND id = ?", uin, dbEventId).Error
	if e != nil {
		if e == gorm.ErrRecordNotFound {
			return nil,cserr.ErrEventNotExist
		}
		return nil, cserr.ErrDbGetFail
	}
	return &de,nil
}

func DbNotifyUserEventRemove(uin string, userEventIdList []string) cserr.ICSErrCodeError {
	db, tb := GetNotifyDbAndTbNameByKey(uin, DbNotifyUserEventTbNameBase)
	de := model.DbNotifyUserEvent{}
	if len(userEventIdList) > 200 {
		return cserr.ErrBadRequest
	}
	if e := db.Table(tb).Delete(&de, "uin = ? AND user_event_id in (?)", uin, userEventIdList).Error; e != nil {
		log.Errorf(e, "del user:%s event list fail", uin)
		return cserr.ErrDb
	}
	return nil
}

func DbNotifyUserEventGetList(uin string,offset,limit int32, idFrom string)(list []*model.DbNotifyUserEvent, err cserr.ICSErrCodeError){
	db, tb := GetNotifyDbAndTbNameByKey(uin, DbNotifyUserEventTbNameBase)
	if limit <= 0 || limit > 100 {
		err = cserr.ErrBadRequest
		return
	}
	dbId,_,_ := model.GetDbEventRouteFromEventId(idFrom)
	var dbList []model.DbNotifyUserEvent
	if e:=db.Table(tb).Order("ID ASC").Offset(offset).Limit(limit).Find(&dbList,
		"uin = ? AND id > ?", uin, dbId).Error; e != nil {
		log.Errorf(e, "get user:%s event list fail", uin)
		err = cserr.ErrDbGetFail
		return
	}
	for i := range dbList {
		list = append(list, &dbList[i])
	}
	err = nil
	return
}


func GetNotifyDbAndTbNameByKey(key,tbNameBase string)(*gorm.DB, string) {
	bk := []byte(key)
	config := conf.GetAppConfig()
	dbname := sfw_db.GetDbNameByKey(bk, config.Notify.UseDb.SplitDbNum, NotifyDbTbNameBase)
	tbname := sfw_db.GetTbNameByKey(bk, config.Notify.UseDb.SplitDbNum, config.Notify.UseDb.SplitTbNum, tbNameBase)
	return sfw_db.GetDb(dbname), tbname
}


func GetDbNotifyUserEventDbAndTbNameByUserEventId(userEventId string)(*gorm.DB, string) {
	_,dbidx,tbidx := model.GetDbEventRouteFromEventId(userEventId)
	dbname := sfw_db.GetDbNameByIdx(NotifyDbTbNameBase, dbidx)
	tb := sfw_db.GetTbNameByIdx(DbNotifyUserEventTbNameBase, dbidx, tbidx)
	db := sfw_db.GetDb(dbname)
	return db,tb
}

func DbNotifyUserEventGetByUserEventId(userEventId string)(event *model.DbNotifyUserEvent,err cserr.ICSErrCodeError) {
	dbid,dbidx,tbidx := model.GetDbEventRouteFromEventId(userEventId)
	dbname := sfw_db.GetDbNameByIdx(NotifyDbTbNameBase, dbidx)
	tb := sfw_db.GetTbNameByIdx(DbNotifyUserEventTbNameBase, dbidx, tbidx)
	db := sfw_db.GetDb(dbname)
	event = &model.DbNotifyUserEvent{}
	e := db.Table(tb).First(event, "id = ? AND user_event_id = ?", dbid, userEventId).Error
	if e != nil {
		if e == gorm.ErrRecordNotFound {
			err = cserr.ErrEventNotExist
			return
		}
		err= cserr.ErrDb
		return
	}
	err = nil
	return
}

func GetNotifyDbidxAndTbidxByKey(key string)(dbidx,tbidx uint32){
	code := sfw_db.GetDbKeyHashCode([]byte(key))
	config := conf.GetAppConfig()
	return code % config.Notify.UseDb.SplitDbNum, code % config.Notify.UseDb.SplitTbNum
}


func DbUserNotifyEventAdd(Uin string, evType comm.UserEventType, intParam int64, strParam string, evParam proto.Message) (dbEvent *model.DbNotifyUserEvent,err cserr.ICSErrCodeError) {
	dbEvent = &model.DbNotifyUserEvent{
		Uin: Uin,
		EventType: evType,
		StrParam: strParam,
		IntParam: intParam,
	}
	if evParam != nil {
		b, e:= proto.Marshal(evParam)
		if util.CheckError(e, "user:%s event:%d pack error", Uin, evType) {
			err = cserr.ErrProto
			return
		}
		dbEvent.EvtParam = b
	}
	dbidx,tbidx := GetNotifyDbidxAndTbidxByKey(Uin)
	dbname := sfw_db.GetDbNameByIdx(NotifyDbTbNameBase, dbidx)
	tb := sfw_db.GetTbNameByIdx(DbNotifyUserEventTbNameBase, dbidx, tbidx)
	db := sfw_db.GetDb(dbname)
	e := db.Table(tb).Create(dbEvent).Error
	if e != nil {
		log.Errorf(e, "db create user:%s event:%s fail", Uin, evType)
		err = cserr.ErrDb
		return
	}
	userEventId := model.GetEventIdFromDbEventRoute(dbEvent.ID, dbidx, tbidx)
	e = db.Table(tb).Where("id = ?", dbEvent.ID).Update("user_event_id", userEventId).Error
	if e != nil {
		log.Errorf(e, "db update user:%s event:%s event id:%s fail", Uin, evType,userEventId)
		err = cserr.ErrDb
		return
	}
	err = nil
	return
}