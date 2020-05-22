package model

import (
	"github.com/jinzhu/gorm"
	"groot/proto/comm"
	"groot/proto/csmsg"
)


type DbChatComplain struct {
	gorm.Model
	Uin					string		`gorm:"index"`
	TargetType 	comm.MsgTargetType	`gorm:"index"`
	TargetId	string 				`gorm:"index"`
	ComplainType csmsg.ComplainType	`gorm:"index"`
	Content 	string
	Progress 	int32
	ProcessNote 		string
}

//共享会话统计
type DbChatSession struct {
	ChatSessionId 		string	`gorm:"primary_key"`
	LastMsgUniqId		uint64	//最新的全局消息序号
	LastSessionMsgSeq 	uint64
}


//好友会话的设置.普通存储
type DbChatSetting struct {
	UinMine					string	`gorm:"unique_index:chat_setting_key"`
	TargetType 	comm.MsgTargetType	`gorm:"unique_index:chat_setting_key"`
	TargetId	string 				`gorm:"unique_index:chat_setting_key"`
	//
	FriendChatStoreName 	string
	FriendChatSessionId 	string
	ReadMsgId				uint64	//已读的最近一个消息序号
	IsQuiet                bool   `json:"is_quiet" bson:"is_quiet"`             //是否静音
	ChatBackGround         string `json:"chat_back_ground" bson:"chat_back_ground"`
	LastUpdateTimeStamp    int64  `json:"last_update_time_stamp" bson:"last_update_time_stamp"`
	LastClearChatTimeStamp int64  `json:"last_clear_chat_time_stamp" bson:"last_clear_chat_time_stamp"`
	RequestFriendMessage string `json:"request_friend_message" bson:"request_friend_message"`
}


