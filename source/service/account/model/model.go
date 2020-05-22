package model

import (
	"github.com/jinzhu/gorm"
	"groot/proto/comm"
	"time"
)

type (
	DbUserLoginEmail struct {
		Email    string `gorm:"primary_key"`
		Uin 	 string	`gorm:"unique_index"`
	}
	DbUserLoginPhone struct {
		Phone    string `gorm:"primary_key"`
		Uin 	 string	`gorm:"unique_index"`
	}
	DbUserLoginName struct {
		UserName                    string `gorm:"primary_key"`
		Uin 	 					string	`gorm:"unique_index"`
		Password                    string
	}
	DbUserLoginDevice struct {
		gorm.Model
		//foreign key
		Uin 		string `gorm:"uniq_index:index_user_device"`
		//设备标识
		DeviceID 	string `gorm:"uniq_index:index_user_device"`
		OsType              comm.DeviceOsType	`gorm:"index"`
		DeviceType  string
		//
		IsOnline            bool   
		DeviceName          string 
		LoginAuthType       comm.UserAuthType //默认手机登陆
		LoginImei           string
		LoginDeviceInfo     string 
		LoginIP             string 
		LoginRealCountry    string 
		BundleID            string 
		Language            string 
		TimeZone            string 
		AutoAuthKey         string 
		ClientDBEncryptKey  string 
		ClientDBEncryptInfo string 
		LoginClientVersion  int64  
		LoginTimestamp      int64  
		LogoutTimestamp     int64  
		DevicePushToken     string 
		UnreadCount         uint32 
	}
	DbUser struct {
		ID        uint64 `gorm:"primary_key"`
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time `sql:"index"`
		//////////////////////////////////
		Uin 			string	`gorm:"unique_index"`
		LoginTel		string	`gorm:"unique_index"`
		LoginUserName	string	`gorm:"unique_index"`
		LoginEmail		string	`gorm:"unique_index"`
		/////////////////////////////////////////////
		///security
		ModifyUserNameLastTimeStamp int64 `gorm:"default:0"`
		ModifyPasswdLastTimeStamp int64 `gorm:"default:0"`
		UserType                    comm.UserType	`gorm:"index"`
		RegisterInitTimeStamp       int64
		TicketSecret                string
	}
	DbUserTokenState struct {
		Uin 		string `gorm:"unique_index:index_user_device"`
		DeviceID 	string `gorm:"unique_index:index_user_device" json:"device_id" bson:"device_id"`
		TokenSecKey         string
		JwtExpiredTime      int64
		JwtIssueTime        int64
	}
)

