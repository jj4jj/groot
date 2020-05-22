package feedback

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"groot/proto/cserr"
	"groot/proto/csmsg"
	"groot/service/feedback/model"
	"groot/sfw/db"
	"groot/sfw/util"
)

var (
	FeedbackDb     *gorm.DB
	FeedbackDbName string
	FeedBackTbName string
)

func init(){
	FeedbackDbName = "db_feedback"
	FeedBackTbName = sfw_db.GetModelTableName(&model.DbFeedback{})
}

func InitDbContext() error {
	FeedbackDb = sfw_db.GetDb(FeedbackDbName)
	return nil
}


func DbFeedbackCreate(uin,content,qq,extInfo,tel string, device *csmsg.CSUserLoginDevice) cserr.ICSErrCodeError {
	dj,err := json.Marshal(device)
	if util.CheckError(err, "json pack err") {
		return cserr.ErrProto
	}
	ext := & struct {
		qq  string
		tel string
		info string
	}{
		qq: qq,
		tel: tel,
		info: extInfo,
	}
	ej, err := json.Marshal(&ext)
	if util.CheckError(err, "json pack err"){
		return cserr.ErrProto
	}
	e := FeedbackDb.Table(FeedBackTbName).Create(model.DbFeedback{
		Uin: uin,
		Content:content,
		DeviceJson:dj,
		Ext:ej,
	}).Error
	if util.CheckError(e, "create feed back err") {
		return cserr.ErrDb
	}
	return nil
}