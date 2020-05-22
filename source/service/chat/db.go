package chat

import (
	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"groot/comm/conf"
	"groot/proto/comm"
	"groot/proto/cserr"
	"groot/proto/csmsg"
	"groot/service/chat/model"
	sfw_db "groot/sfw/db"
	"groot/sfw/util"
)


var (
	ChatComplainDb          * gorm.DB
	ChatComplainDbName      string
	ChatComplainTbName      string
	ChatDbNameBase          string
	DbChatSettingTbNameBase string
	DbChatSessionTbNameBase string
)

func init(){
	ChatComplainDbName = "db_misc_chat_complain"
	ChatDbNameBase = "db_chat"
	ChatComplainTbName = sfw_db.GetModelTableName(&model.DbChatComplain{})
	DbChatSettingTbNameBase = sfw_db.GetModelTableNameBase(&model.DbChatSetting{})
	DbChatSessionTbNameBase = sfw_db.GetModelTableNameBase(&model.DbChatSession{})
}

func GetChatDbAndTbNameByKey(tbNameBase string, key []byte)(*gorm.DB,string) {
	config := conf.GetAppConfig()
	dbname := sfw_db.GetDbNameByKey(key, config.Chat.UseDb.SplitDbNum, ChatDbNameBase)
	tbname := sfw_db.GetTbNameByKey(key, config.Chat.UseDb.SplitDbNum,
		config.Chat.UseDb.SplitTbNum, tbNameBase)
	return sfw_db.GetDb(dbname),tbname
}
func GetChatSettingDbAndTbName(UinMine string) (*gorm.DB,string) {
	return GetChatDbAndTbNameByKey(DbChatSettingTbNameBase, []byte(UinMine))
}

func DbChatSessionUpdateLastMsgUniqId(sessionId string, lastMsgUniqId uint64) cserr.ICSErrCodeError {
	chatSession := model.DbChatSession{
		LastMsgUniqId:lastMsgUniqId,
	}
	db,tbname := GetChatDbAndTbNameByKey(DbChatSessionTbNameBase, []byte(sessionId))
	err := db.Table(tbname).FirstOrCreate(&chatSession, "chat_session_id = ?", sessionId).Error
	if err != nil {
		log.Errorf(err, "chat session get fail", sessionId)
		return cserr.ErrDb
	}

	err = db.Table(tbname).Where("chat_session_id = ? AND last_session_msg_seq = %d",sessionId,
		chatSession.LastSessionMsgSeq).Updates(model.DbChatSession{
			LastMsgUniqId: lastMsgUniqId,
			LastSessionMsgSeq: chatSession.LastSessionMsgSeq+1,
		}).Error
	if err != nil {
		log.Errorf(err, "chat session update seq err")
		return cserr.ErrDbUpdateFail
	}

	return nil
}

func DbChatSessionGetLastMsgUniqIdWithCheckSeq(sessionId string, checkSeq uint64) (uint64,cserr.ICSErrCodeError) {
	chatSession := model.DbChatSession{}
	db,tbname := GetChatDbAndTbNameByKey(DbChatSessionTbNameBase, []byte(sessionId))
	err := db.Table(tbname).First(&chatSession, "chat_session_id = ?", sessionId).Error
	if err == gorm.ErrRecordNotFound {
		return 0,nil
	}
	if err != nil {
		return 0,cserr.ErrDbGetFail
	}
	if chatSession.LastSessionMsgSeq != checkSeq {
		log.Debugf("chat session check msg seq")
		return 0, cserr.ErrCheckSeq
	}
	return chatSession.LastMsgUniqId,nil
}




func DbChatSettingUpdate(UinMine string, targetType comm.MsgTargetType, targetId string,
	mod *model.DbChatSetting) cserr.ICSErrCodeError {

	db, tbname := GetChatSettingDbAndTbName(UinMine)

	err := db.Table(tbname).Where("uin = ? AND target_type = ? AND target_id = ?", UinMine, targetType,
		targetId).Updates(*mod).Error

	if util.CheckError(err, "db update error with uin=%s", UinMine) {
		return cserr.ErrDb
	}
	return nil
}

func DbChatCompainInsert(UinMine string, targetType comm.MsgTargetType, targetId string,
				cmpType csmsg.ComplainType, content string)  cserr.ICSErrCodeError {
	var count uint
	err := ChatComplainDb.Table(ChatComplainTbName).Where("uin = %s", UinMine).Count(&count).Error
	if util.CheckError(err, "count uin:%s error", UinMine) {
		return cserr.ErrDb
	}
	if count > 20 {
		log.Warnf("uin:%s insert complain too much !", UinMine)
		return cserr.ErrChatComplainTooMuch
	}

	err = ChatComplainDb.Table(ChatComplainTbName).Create(&model.DbChatComplain{
		Uin:UinMine,
		TargetType: targetType,
		TargetId: targetId,
		ComplainType: cmpType,
		Content: content,
	}).Error

	if util.CheckError(err, "create complain uin:%s error", UinMine) {
		return cserr.ErrDb
	}
	return nil
}


func initChatDbCtx() error {

	ChatComplainDb = sfw_db.GetDb(ChatComplainDbName)

	return nil
}