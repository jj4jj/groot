package model

import "github.com/lexkong/log"

//never modify ////////////////////////////////////////////////////
//[dbid].[reserve.2].[db(4).tb(3)]
const DbChatMsgIdReserveBase = 100
const DbChatMsgTbMaxNum = 1000
const DbChatDbMaxNum = 10000
const DbChatReserveIdFlags = 1
func GetMsgUniqIdFromDbRoute(msgDbId,dbidx,tbidx uint32) uint64 {
	var ChatMsgId = uint64(msgDbId)
	ChatMsgId *= DbChatMsgIdReserveBase
	ChatMsgId += DbChatReserveIdFlags
	ChatMsgId *= DbChatDbMaxNum
	ChatMsgId += uint64(dbidx)
	ChatMsgId *= DbChatMsgTbMaxNum
	ChatMsgId += uint64(tbidx)
	return ChatMsgId
}
func ParseDbRouteByMsgUniqId(msgUniqueId uint64) (msgDbId,dbidx,tbidx uint32) {
	tbidx = uint32(msgUniqueId % DbChatMsgTbMaxNum)
	msgUniqueId /= DbChatMsgTbMaxNum
	dbidx = uint32(msgUniqueId % DbChatDbMaxNum)
	msgUniqueId /= DbChatDbMaxNum
	flag := msgUniqueId%DbChatMsgIdReserveBase
	if flag == DbChatReserveIdFlags {
		log.Warnf("error uniq id flag:%x", flag)
	}
	msgUniqueId /= DbChatMsgIdReserveBase
	msgDbId = uint32(msgUniqueId)
	return
}
//never modify ////////////////////////////////////////////////////



