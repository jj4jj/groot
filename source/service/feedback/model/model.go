package model

import "github.com/jinzhu/gorm"

type DbFeedback struct {
	gorm.Model
	Uin 		string	`gorm:"index"`
	Content 	string
	DeviceJson	[]byte	//json
	Ext 		[]byte	//json
}