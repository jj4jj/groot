package friend

import (
	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"groot/comm/conf"
	"groot/proto/cserr"
	"groot/proto/csmsg"
	"groot/service/friend/model"
	"groot/sfw/db"
	"strings"
)

var (
	FriendDbNameBase 				string
	DbFriendSettingTbNameBase		string
	DbFriendShipTbNameBase			string
)
func init(){
	FriendDbNameBase = "db_friend"
	DbFriendSettingTbNameBase = sfw_db.GetModelTableNameBase(&model.DbFriendSetting{})
	DbFriendShipTbNameBase = sfw_db.GetModelTableNameBase(&model.DbFriendship{})
}

func GetFriendDbByDbIdx(dbidx uint32) *gorm.DB {
	dbname := sfw_db.GetDbNameByIdx(FriendDbNameBase, dbidx)
	return sfw_db.GetDb(dbname)
}

func GetFriendDbAndTbNameByUin(uin, tbname_base string) (*gorm.DB, string) {
	config := conf.GetAppConfig()
	dbname := sfw_db.GetDbNameByKey([]byte(uin), config.Friend.UseDb.SplitDbNum, FriendDbNameBase)
	tbname := sfw_db.GetTbNameByKey([]byte(uin), config.Friend.UseDb.SplitDbNum, config.Friend.UseDb.SplitTbNum, tbname_base)
	return sfw_db.GetDb(dbname), tbname
}

func DbFriendBatchGetSetting(UinMine string, UinFriendList []string, mapResult map[string]*model.DbFriendSetting){
	db,tbname := GetFriendDbAndTbNameByUin(UinMine, DbFriendSettingTbNameBase)
	var settingList []model.DbFriendSetting
	e := db.Table(tbname).Where("uin in (?)", UinFriendList).Find(&settingList).Error
	if e != nil {
		return
	}
	for i,_ := range settingList {
		mapResult[settingList[i].FriendUin] = &settingList[i]
	}
}

func DbFriendSettingUpdate(UinMine,UinFriend string, note *csmsg.FriendUserNote, setting *csmsg.FriendUserSetting) cserr.ICSErrCodeError {
	db,tbname := GetFriendDbAndTbNameByUin(UinMine, DbFriendSettingTbNameBase)
	dbSettingInit := model.DbFriendSetting{
		Uin: UinMine,
		FriendUin: UinFriend,
	}
	if e:= db.Table(tbname).FirstOrCreate(&dbSettingInit, "uin = ? AND friend_uin = ?", UinMine, UinFriend).Error;
		e != nil {
		log.Errorf(e, "db first or create error")
		return cserr.ErrDb
	}
	var updateAttrs model.DbFriendSetting
	if note != nil {
		updateAttrs.NickName = note.NickName
		updateAttrs.TagList= strings.Join(note.TagList,",")
		updateAttrs.Desc= note.Desc
		updateAttrs.Telephone= note.Telephone
		updateAttrs.GroupList= strings.Join(note.GroupList, ",")
	}
	if setting != nil {
		updateAttrs.IsStar = setting.IsStar
	}
	err := db.Table(tbname).Where("uin = ? AND friend_uin = ?", UinMine,
		UinFriend).Updates(updateAttrs).Error
	if err != nil {
		log.Errorf(err, "get friend setting fail")
		return cserr.ErrDbUpdateFail
	}

	return nil
}

func DbFriendshipGetByUin(UinMine,UinFriend string) *model.DbFriendship {
	dbmine, tbmine := GetFriendDbAndTbNameByUin(UinMine, DbFriendShipTbNameBase)
	ship := & model.DbFriendship{}
	e := dbmine.Table(tbmine).First(ship,"uin = %s AND friend_uin = ?", UinMine, UinFriend).Error
	if e != nil {
		if e == gorm.ErrRecordNotFound {
			return nil
		}
		log.Errorf(e, "db error")
		return nil
	}
	return ship
}

func initFriendDbCtx() error {
	return nil
}
