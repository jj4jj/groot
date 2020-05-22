package model

import (
	"github.com/jinzhu/gorm"
	"groot/proto/comm"
)

//这个存储在消息顺序消息数据库(支持TB级).这里会分库和表
type DbChatMsg struct {
	gorm.Model
	ChatSessionId			string			`gorm:"unique_index"`
	SenderUin				string			`gorm:"index"`
	MsgType 				comm.MsgType 	`gorm:"index"`
	Content 				[]byte
}



