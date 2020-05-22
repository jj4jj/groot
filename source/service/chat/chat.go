package chat

import (
	"fmt"
	"groot/proto/comm"
	"groot/proto/cserr"
)

func InsertFriendChatMsg(UinMine,UinFriend string, msgType comm.MsgType, msgContent []byte) cserr.ICSErrCodeError {
	return nil
}


func GetFriendChatSessionId(Uin1,Uin2 string) string {
	if Uin1 <= Uin2 {
		return fmt.Sprintf("ch:%s-%s",Uin1,Uin2)
	}
	return fmt.Sprintf("ch:%s-%s",Uin2,Uin1)
}


