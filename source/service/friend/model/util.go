package model

import (
	"groot/proto/csmsg"
	"strings"
)

func DbFriendUserSettingToCsUserNote(setting *DbFriendSetting) *csmsg.FriendUserNote {
	ret := csmsg.FriendUserNote{}
	ret.GroupList = strings.Split(setting.GroupList, ",")
	ret.TagList = strings.Split(setting.TagList, ",")
	ret.Telephone = setting.Telephone
	ret.Desc = setting.Desc
	ret.NickName = setting.NickName
	return &ret
}

func DbFriendUserSettingToCsUserSetting(setting *DbFriendSetting) *csmsg.FriendUserSetting {
	ret := csmsg.FriendUserSetting{}
	ret.IsStar = setting.IsStar
	return &ret
}


