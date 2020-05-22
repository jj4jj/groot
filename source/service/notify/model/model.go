package model

import (
	"groot/proto/comm"
	"time"
)

//user event route key is uin
type DbNotifyUserEvent struct {
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	/////////////////////////////////////
	Uin           string             `gorm:"index"`
	UserEventId   string             `gorm:"unique_index"`
	EventType     comm.UserEventType `gorm:"index"`
	IntParam      int64
	StrParam      string
	EvtParam      []byte 	`gorm:"size:8192"`
	EvtState      int64	`gorm:"index"`
	EvtStateParam []byte `gorm:"size:2048"`
}