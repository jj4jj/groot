package model

import (
	"groot/proto/comm"
	"time"
)

//1. frien realation ship
//2. friend settings
const (
	FriendAskMaxTimeDuring time.Duration = 7 * 24 * time.Hour 	//7 days
)

//好友关系
type DbFriendship struct {
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Uin				string	`gorm:"unique_index:friend_relation"`
	FriendUin		string	`gorm:"unique_index:friend_relation"`
	FriendStatus	comm.FriendStatus		`gorm:"index"`
	FriendType      comm.FriendType   `gorm:"index" json:"friend_type" bson:"friend_type"`
	FriendSource 	comm.UserSourceType	`gorm:"index"`
}

//好友设置
type DbFriendSettingNote struct {
	NickName             string   `protobuf:"bytes,1,opt,name=nick_name,json=nickName,proto3" json:"nick_name,omitempty"`
	TagList              string `protobuf:"bytes,2,rep,name=tag_list,json=tagList,proto3" json:"tag_list,omitempty"`
	Desc                 string   `protobuf:"bytes,3,opt,name=desc,proto3" json:"desc,omitempty"`
	Telephone            string   `protobuf:"bytes,4,opt,name=telephone,proto3" json:"telephone,omitempty"`
	GroupList			 string
}

type DbFriendSetting struct {
	Uin				string	`gorm:"unique_index:friend_relation_setting"`
	FriendUin		string	`gorm:"unique_index:friend_relation_setting"`
	BecomeFriendTimeStamp	int64            `json:"become_friend_time_stamp" bson:"become_friend_time_stamp"`
	IsStar        	bool   `json:"is_star" bson:"is_star"`
	DbFriendSettingNote
	/*
	FriendUserName        string             `json:"friend_user_name" bson:"friend_user_name"` //如果是单聊，就是对方UserName，如果是群聊，则为群聊的username
	FriendToken   string `json:"friend_token" bson:"friend_token"`       //好友操作token
	FriendSource  int    `json:"friend_source" bson:"friend_source"`     //来源
	CommonChatBar string `json:"common_chat_bar" bson:"common_chat_bar"`
	IsWatchHerTrends       bool   `json:"is_watch_her_trends" bson:"is_watch_her_trends"`           //是否看他的动态
	IsGiveHerWatchTrends   bool   `json:"is_give_her_watch_trends" bson:"is_give_her_watch_trends"` //是否看他的动态
	 */
}










