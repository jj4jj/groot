package util

import (
	"fmt"
	"github.com/lexkong/log"
	"groot/proto/cserr"
)

func ErrOK(err error) bool {
	return err == nil || err == cserr.ErrSuccess
}

func ErrStr(err error) string {
	if err == nil {
		return "<ok:nil>"
	} else if ic,ok := err.(cserr.ICSErrCodeError); ok {
		return ic.Error()
	} else {
		return "<unknown:-1>"
	}
}
func ErrCode(err error) int32 {
	if err == nil {
		return 0
	} else if ic,ok := err.(cserr.ICSErrCodeError); ok {
		return int32(ic.Errno())
	} else {
		return -1
	}
}

func CheckError(err error, sfmt string, args ...interface{}) bool {
	if !ErrOK(err) {
		if sfmt != "" {
			log.Errorf(err, sfmt, args...)
		}
		return true
	}
	return false
}

func CatchExcpetion(fn func(interface{}), sfmt string, args ...interface{}){
	if err := recover(); err != nil {
		if sfmt != "" {
			info := fmt.Sprintf(sfmt, args...)
			log.Warnf("call func error:%v info:[%s]", err, info)
		}
		if fn != nil{
			fn(err)
		}
	}
}