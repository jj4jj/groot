package friend

import (
	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"groot/proto/comm"
	"groot/proto/cserr"
	"groot/service/account/crud"
	"groot/service/chat"
	"groot/service/friend/model"
	"groot/sfw/util"
	"time"
)

func DbFriendAsk(UinMine,UinFriend string, content string, source comm.UserSourceType) cserr.ICSErrCodeError {
	systemUin := crud.GetSystemUin()
	if UinFriend == systemUin {
		return cserr.ErrBadRequest
	}

	dbmine, tbmine := GetFriendDbAndTbNameByUin(UinMine, DbFriendShipTbNameBase)

	friendship := model.DbFriendship{
		Uin:          UinMine,
		FriendUin:    UinFriend,
		FriendStatus: comm.FriendStatus_FRIEND_STATUS_STRANGER,
		FriendSource: source,
	}
	err := dbmine.Table(tbmine).FirstOrCreate(&friendship, "uin = ? friend_uin = ?", UinMine, UinFriend).Error
	if err != nil {
		log.Errorf(err, "friendship db fc error")
		return cserr.ErrDb
	}

	//ok friend ship got
	if friendship.FriendStatus == comm.FriendStatus_FRIEND_STATUS_NORMAL {
		return cserr.ErrFriendExistAlready
	}
	if friendship.FriendStatus == comm.FriendStatus_FRIEND_STATUS_WAIT_ANSWER {
		log.Debugf("UinSigner:%s Add Friend UinSigner:%s is in black equal stranger", UinMine, UinFriend)
	}

	//
	if util.CheckError(chat.InsertFriendChatMsg(UinMine, UinFriend, comm.MsgType_MSG_TYPE_FRIEND_ASK, []byte(content)),
		"insert hello msg") {
		return cserr.ErrFriendInsertMsg
	}

	//set state
	if friendship.FriendStatus != comm.FriendStatus_FRIEND_STATUS_WAIT_ANSWER {
		friendship.FriendStatus = comm.FriendStatus_FRIEND_STATUS_WAIT_ANSWER
		dbmine.Table(tbmine).Where("uin = ? AND friend_uin = ?", UinMine,
			UinFriend).Update(model.DbFriendship{
			FriendStatus: comm.FriendStatus_FRIEND_STATUS_WAIT_ANSWER,})
	}

	//check friend status
	friendState := DbFriendGetStatus(UinFriend, UinMine)
	if friendState == comm.FriendStatus_FRIEND_STATUS_BLACK {
		log.Debugf("uin:%s ask add friend:%s but in black", UinMine, UinFriend)
		return cserr.ErrFriendStatus
	}
	if friendState == comm.FriendStatus_FRIEND_STATUS_NORMAL {
		util.CheckError(dbmine.Table(tbmine).Where("uin = ? AND friend_uin = ?", UinMine,
			UinFriend).Update("FriendStatus", comm.FriendStatus_FRIEND_STATUS_NORMAL).Error,"recv friend")
		//event notify user add friend success
		NotifyUserFriendNew(UinMine, UinFriend)
	} else {
		util.CheckError(dbmine.Table(tbmine).Where("uin = ? AND friend_uin = ?", UinMine,
			UinFriend).Update("FriendStatus", comm.FriendStatus_FRIEND_STATUS_WAIT_ANSWER).Error,"update state")
		//event notify friend user new friend apply
		NotifyUserFriendAsk(UinMine, UinFriend, content)
	}

	//
	if util.CheckError(chat.InsertFriendChatMsg(UinMine, UinFriend, comm.MsgType_MSG_TYPE_FRIEND_ASK,
		[]byte(content)),"insert hello msg") {
		return cserr.ErrFriendInsertMsg
	}

	return nil
}

func DbFriendGetStatus(UinMine, UinFriend string) comm.FriendStatus {
	systemUin := crud.GetSystemUin()
	if UinMine == systemUin || UinFriend == systemUin {
		return comm.FriendStatus_FRIEND_STATUS_NORMAL
	}
	dbmine, tbmine := GetFriendDbAndTbNameByUin(UinMine, DbFriendShipTbNameBase)
	friendShip := &model.DbFriendship{}
	err := dbmine.Table(tbmine).First(friendShip,"uin = ? AND friend_uin = ?", UinMine,
		UinFriend).Error
	if gorm.IsRecordNotFoundError(err) {
		return comm.FriendStatus_FRIEND_STATUS_STRANGER
	}
	if util.CheckError(err, "db get friend status error"){
		return comm.FriendStatus_FRIEND_STATUS_NONE
	}
	return friendShip.FriendStatus
}

func DbFriendRemove(UinMine,UinFriend string) cserr.ICSErrCodeError {
	if crud.GetSystemUin() == UinFriend {
		return cserr.ErrBadRequest
	}
	dbmine, tbmine := GetFriendDbAndTbNameByUin(UinMine, DbFriendShipTbNameBase)
	ship := model.DbFriendship{}
	err := dbmine.Unscoped().Table(tbmine).Delete(&ship, "uin = ? AND friend_uin = ?", UinMine, UinFriend).Error
	if util.CheckError(err , "delete friend uin:%s ->  friend uin:%s fail", UinMine, UinFriend) {
		return cserr.ErrDb
	}
	return nil
}

func DbFriendCheckAddInitSystemFriend(uin string, db * gorm.DB, tbname string) {
	//第一次拉去好友初始化系统信息.
	dbship := model.DbFriendship{}
	systemUin := crud.GetSystemUin()
	if db.Table(tbname).First(&dbship, "uin = ? AND friend_uin = ?", uin,
		systemUin).RecordNotFound() {
		dbship.Uin = uin
		dbship.FriendUin = systemUin
		dbship.FriendSource = comm.UserSourceType_USER_SOURCE_INIT
		dbship.FriendStatus = comm.FriendStatus_FRIEND_STATUS_NORMAL
		dbship.FriendType = comm.FriendType_FRIEND_TYPE_SYSTEM
		if e:=db.Table(tbname).Create(&dbship).Error; e != nil {
			log.Errorf(e, "create friend ship error uin:%s fail",uin)
		}
	}
}

func DbFriendGetList(Uin string, offset,limit int32, status comm.FriendStatus) (list []model.DbFriendship,
	err	cserr.ICSErrCodeError ){
	if limit < 0 || limit > 100 || offset < 0 {
		err = cserr.ErrBadRequest
		return
	}
	dbmine, tbmine := GetFriendDbAndTbNameByUin(Uin, DbFriendShipTbNameBase)
	var dberr error
	if status == comm.FriendStatus_FRIEND_STATUS_NONE {
		DbFriendCheckAddInitSystemFriend(Uin,dbmine, tbmine)
		dberr = dbmine.Table(tbmine).Where("uin = %s", Uin).Order("friend_uin DESC").
			Offset(offset).Limit(limit).Find(&list).Error
	} else {
		dberr = dbmine.Table(tbmine).Where("uin = %s AND status = %d", Uin, status).Order("friend_uin DESC").
			Offset(offset).Limit(limit).Find(&list).Error
	}
	if util.CheckError(dberr, "get friend list db err") {
		err = cserr.ErrDb
	} else {
		err = nil
	}
	return
}

func DbFriendSetBlock(UinMine,UinFriend string, block bool) cserr.ICSErrCodeError {
	dbmine, tbmine := GetFriendDbAndTbNameByUin(UinMine, DbFriendShipTbNameBase)
	if !block {
		blockVerify := comm.FriendStatus_FRIEND_STATUS_BLACK
		blockSet := comm.FriendStatus_FRIEND_STATUS_NORMAL
		if util.CheckError(dbmine.Table(tbmine).Where("uin = ? AND friend_uid = ? AND friend_status = ?", UinMine,
			UinFriend, blockVerify).Update(model.DbFriendship{
			FriendStatus: blockSet,
		}).Error, "updatge db error") {
			return cserr.ErrDb
		}

	} else {
		if util.CheckError(dbmine.Table(tbmine).Where("uin = ? AND friend_uid = ?", UinMine,
			UinFriend).Update(model.DbFriendship{
			FriendStatus: comm.FriendStatus_FRIEND_STATUS_BLACK,
		}).Error, "update db err") {
			return cserr.ErrDb
		}
	}
	return nil
}

func GetFriendUinList(uinMine string, status ...comm.FriendStatus) []string {
	var uinList []string
	dbmine, tbmine := GetFriendDbAndTbNameByUin(uinMine, DbFriendShipTbNameBase)
	var dbList []model.DbFriendship
	e := dbmine.Table(tbmine).Select("friend_uin").Where("uin = %s AND friend_status in (?)",
		uinMine, status).Find(&dbList).Error
	if e != nil {
		if e == gorm.ErrRecordNotFound {
			return nil
		}
		log.Errorf(e, "get friend list error")
		return nil
	}
	for i := range dbList {
		uinList = append(uinList, dbList[i].FriendUin)
	}
	return uinList
}

func CheckFriendStatusEachOther(UinMine string, UinCheck string, stateMust comm.FriendStatus) bool {
	return util.CheckError(FriendCheckStatus(UinMine, []string{UinCheck}, stateMust),"") &&
			util.CheckError(FriendCheckStatus(UinCheck, []string{UinMine}, stateMust),"")
}

func FriendCheckStatus(UinMine string, UinCheckList []string, stateMust comm.FriendStatus) cserr.ICSErrCodeError {
	dbmine, tbmine := GetFriendDbAndTbNameByUin(UinMine, DbFriendShipTbNameBase)
	var cnt = 0
	e := dbmine.Table(tbmine).Where("uin = %s AND friend_uin in (?) AND friend_status = ?", UinMine,
		UinCheckList, stateMust).Count(&cnt).Error
	if e != nil {
		log.Errorf(e, "get friend list error")
		return cserr.ErrDb
	}
	if cnt != len(UinCheckList) {
		log.Warnf("uin:%s friend state:%d list get num:%d check num:%d", UinMine, stateMust,
			cnt, len(UinCheckList))
		return cserr.ErrFriendStatus
	}
	return nil
}

func DbFriendAnswer(UinMine,UinFriend string, content string, agree bool) cserr.ICSErrCodeError {
	dbmine, tbmine := GetFriendDbAndTbNameByUin(UinMine, DbFriendShipTbNameBase)
	dbfriend, tbfriend := GetFriendDbAndTbNameByUin(UinFriend, DbFriendShipTbNameBase)

	//find ask friend info
	askFriendship := model.DbFriendship{}
	err := dbfriend.Table(tbfriend).First(&askFriendship, "uin = ? friend_uin = ?", UinFriend,UinMine).Error
	if err != nil {
		log.Errorf(err, "friendship db fc error")
		return cserr.ErrDb
	}

	if askFriendship.FriendStatus != comm.FriendStatus_FRIEND_STATUS_WAIT_ANSWER {
		return cserr.ErrFriendStatus
	}

	timeNow := time.Now()
	if askFriendship.UpdatedAt.Add(model.FriendAskMaxTimeDuring).After(timeNow) {
		//ask timeout
		return cserr.ErrFriendStatus
	}

	if agree {
		//agree
		agreeFriendship := model.DbFriendship{
			Uin:          UinMine,
			FriendUin:    UinFriend,
			FriendStatus: comm.FriendStatus_FRIEND_STATUS_NORMAL,
		}
		err = dbmine.Table(tbmine).FirstOrCreate(&agreeFriendship,"uin = ? AND friend_uin = ?",
			UinMine, UinFriend).Error
		if util.CheckError(err, "friend db err") {
			return cserr.ErrDb
		}
		if agreeFriendship.FriendStatus != comm.FriendStatus_FRIEND_STATUS_NORMAL {
			err = dbmine.Table(tbmine).Where("uin = ? AND friend_uin = ?",
				UinMine, UinFriend).Update("friend_status", comm.FriendStatus_FRIEND_STATUS_NORMAL).Error
			util.CheckError(err, "db update fail")
		}
		NotifyUserFriendNew(UinMine, UinFriend)
		NotifyUserFriendNew(UinFriend, UinMine)
	} else {
		//reject
		NotifyUserFriendReject(UinFriend, UinMine)
	}
	return nil
}


