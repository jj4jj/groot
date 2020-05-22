package model

import (
	"fmt"
	"github.com/lexkong/log"
	"groot/sfw/util"
	"math"
)


//!!!!!!!never modify begin!!!!!!!!!!!!!!!!!!!!!
const DbUserID2UinNumberBase uint64 = 1000000
const MAX_USER_DB_NUM uint64 = 100
const MAX_USER_TB_NUM uint64 = 100
const MAX_ROUTE_RESERVE_DIGIST_NUM = 5
var UserUinReservedNum uint64
func init () {
	UserUinReservedNum = uint64(math.Pow10(MAX_ROUTE_RESERVE_DIGIST_NUM))
}
func GetDbUserRouteFromUin(uin string) (dbUserId uint64, dbidx, tbidx uint32) {
	var uinBin uint64 = 0
	_, err := fmt.Sscanf(uin, "u%x", &uinBin)
	if err != nil {
		log.Errorf(err, "get db user id from:%s error", uin)
		dbUserId = 0
		dbidx = uint32(MAX_USER_DB_NUM)
		tbidx = uint32(MAX_USER_TB_NUM)
		return
	}

	uinBin -= DbUserID2UinNumberBase
	uinBin -= UserUinReservedNum

	//raw
	routeNum := uinBin % (MAX_USER_DB_NUM*MAX_USER_TB_NUM)
	tbidx = uint32(routeNum % MAX_USER_TB_NUM)
	dbidx = uint32(routeNum / MAX_USER_TB_NUM)
	dbUserId = uinBin / (MAX_USER_DB_NUM*MAX_USER_TB_NUM)

	return
}
func GetUinFromDbUserRoute(dbUserID uint64, dbidx, tbidx uint32) string {
	if dbidx >= uint32(MAX_USER_DB_NUM) || tbidx >= uint32(MAX_USER_TB_NUM) {
		log.Warnf("db user dbidx:%d tbidx:%d is invalid", dbidx, tbidx)
		panic("split db user factor dbidx and tbidx is error")
	}

	var uinBin uint64 = dbUserID
	uinBin *= uint64(MAX_USER_DB_NUM*MAX_USER_TB_NUM)
	uinBin += uint64(dbidx) * MAX_USER_TB_NUM
	uinBin += uint64(tbidx)
	uinBin += UserUinReservedNum
	uinBin += DbUserID2UinNumberBase

	s := fmt.Sprintf("%08x", uinBin)
	return "u" + s
}
//!!!!!!!never modify end !!!!!!!!!!!!!!!!!!!!!
func RandomUserNameByTelephone(telephone string) string {
	//telephone for route todo
	return "uid_" + util.RandomString(16)
}



