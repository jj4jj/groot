package msg

import (
	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"groot/comm/conf"
	"groot/proto/cserr"
	"groot/proto/csmsg"
	"groot/service/msg/model"
	"groot/sfw/app"
	"groot/sfw/db"
)

func initMsgDbCtx() error {
	return nil
}

var (
	MsgDbNameBase 	string
	DbChatMsgTbNameBase		string
)

func init(){
	MsgDbNameBase = "db_msg"
	DbChatMsgTbNameBase = sfw_db.GetModelTableNameBase(&model.DbChatMsg{})
}

func GetChatMsgDbAndTbNameByIdx(dbidx,tbidx uint32) (*gorm.DB, string) {
	dbname := sfw_db.GetDbNameByIdx(MsgDbNameBase, dbidx)
	tbname := sfw_db.GetTbNameByIdx(DbChatMsgTbNameBase, dbidx, tbidx)
	return sfw_db.GetDb(dbname), tbname
}

func InsertChatMsgByHint(sessionId,Uin string, content *csmsg.MsgContent, hintMsgUniqId uint64) (uint64,cserr.ICSErrCodeError) {
	config := conf.GetAppConfig()
	var dbidx,tbidx uint32
	if hintMsgUniqId == 0 {
		tbidx = uint32(app.AppRandom.Intn(model.DbChatMsgTbMaxNum))
		if config.RunEnv != "dev" {
			tbidx = 0
		}
		dbidx = config.Msg.EnableSplitDbBegIdx +
			uint32(app.AppRandom.Intn(int(config.Msg.EnableSplitDbEndIdx - config.Msg.EnableSplitDbBegIdx + 1)))
		log.Debugf("session id:%s uin:%s random msg route:%d.%d ...", sessionId,Uin,dbidx,tbidx)
	} else {
		_,dbidx,tbidx = model.ParseDbRouteByMsgUniqId(hintMsgUniqId)
	}

	db,tb := GetChatMsgDbAndTbNameByIdx(dbidx, tbidx)

	chatMsg := model.DbChatMsg{
		ChatSessionId : sessionId,
		SenderUin: Uin,
		MsgType: content.MsgType,
		Content: content.Content,

	}
	err := db.Table(tb).Create(&chatMsg).Error
	if err != nil {
		log.Errorf(err,"create chat msg eror")
		return 0,cserr.ErrDb
	}

	log.Debugf("db chat msg insert msg id:%d", chatMsg.ID)

	return model.GetMsgUniqIdFromDbRoute(uint32(chatMsg.ID), dbidx,tbidx),nil
}