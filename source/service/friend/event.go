package friend

import (
	"github.com/lexkong/log"
	"groot/proto/comm"
	"groot/proto/csmsg"
	"groot/service/notify"
	"groot/service/profile/crud"
	"groot/service/profile/model"
)


func GetFriendBaseByUin(uinMine,uinFriend string) *comm.FriendBase {
	ship := DbFriendshipGetByUin(uinMine, uinFriend)
	if ship == nil {
		return nil
	}
	return &comm.FriendBase{
		Uin: uinFriend,
		Status: ship.FriendStatus,
		FriendType: ship.FriendType,
		Source: ship.FriendSource,
	}
}


func NotifyUserFriendReject(uinMine,uinFriend string) {
	notify.NotifyUserEventNewWithSaveDb(uinFriend, comm.UserEventType_USR_EVT_FRIEND_ANS_REJECT, 0, uinMine, nil)
}

func NotifyUserFriendNew(uinMine,uinFriend string) {
	friendBase := GetFriendBaseByUin(uinMine,uinFriend)
	if friendBase != nil {
		notify.NotifyUserEventNewWithSaveDb(uinMine, comm.UserEventType_USR_EVT_FRIEND_NEW, 0, "",
			&csmsg.NotifyUserEventFriendNew{
				FriendBase: friendBase,
			})
	} else {
		log.Errorf(nil, "found uin:%s friend:%s base fail", uinMine,uinFriend)
	}
}

func NotifyUserFriendAsk(uinMine,uinFriend string, content string) {
	base := crud.DbUserBaseInfoGetByUin(uinMine)
	if base != nil {
		notify.NotifyUserEventNewWithSaveDb(uinFriend, comm.UserEventType_USR_EVT_FRIEND_ASK, 0, "",
			&csmsg.NotifyUserEventFriendAsk{
				TicketAsk: "",
				UinAsk: uinMine,
				Content: content,
				UserBase: model.DbUserBaseToCsUserBase(base),
			})
	} else {
		log.Errorf(nil, "found uin:%s friend:%s base fail", uinMine,uinFriend)
	}
}
