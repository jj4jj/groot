package model

import (
	"encoding/json"
	"github.com/lexkong/log"
	"groot/proto/comm"
	"groot/proto/csmsg"
	"time"
)

type (
	DbUserSetting struct {
		Uin 				string `gorm:"primary_key"`
		EnableSearchByTel 	 int32
		JsonSetting		 	[]byte `gorm:"size:16384"`
		CustomSetting		 []byte	`gorm:"size:32768"` 	//csmsg.JsonUserProfileSetting
	}
	DbUserBaseInfo struct {
		Uin 			string `gorm:"primary_key"`
		NickName 		string	`gorm:"index"`
		HeadIcon		string
		SexType 		comm.UserSexType	`gorm:"index"`
		Desc 			string
		LocationDesc	string
		Longitude		int64	`gorm:"index"`
		Latitude		int64	`gorm:"index"`
		BirthDay 		time.Time	`gorm:"index"`
	}
	DbUserProfileInfo struct {
		Uin 			string `gorm:"primary_key"`
		College			string
		Job 			string
		Education		string
	}
	DbUserRegisterInfo struct {
		Uin 		  string `gorm:"primary_key"`
		//注册时其他相关信息
		Ticket        string 
		Ip            string
		MacAddr       string 
		Imei          string
		Scene         int32                 //是哪个渠道注册的
		OsType 		  comm.DeviceOsType
		AuthType      comm.UserAuthType                   //默认手机号注册
		TimeZone      string        //注册时的时区
		Language      string          //注册时的语言
		RealCountry   string  //注册时国家
		DeviceInfo    string   //注册时手机信息
		BundleId      string        //注册时bundleid
		ClientVersion int64  
		TimeStamp     int64  `gorm:"default:0" json:"register_time_stamp" bson:"register_time_stamp"`
	}
)


func (s *DbUserSetting) GetJsonUserSetting() *csmsg.JsonUserProfileSetting {
	js := csmsg.JsonUserProfileSetting{}
	if e:= json.Unmarshal(s.JsonSetting, &js); e != nil {
		log.Errorf(e, "json unmarshal setting fail")
		return nil
	}
	return &js
}
