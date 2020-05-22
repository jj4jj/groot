package stream

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/lexkong/log"
	"groot/proto/comm"
	"groot/sfw/app"
	"groot/sfw/util"
	"sync"
)

const (
	STREAM_CONNX_ANONYMOUSE ConnxState = iota
	STREAM_CONNX_AUTHED_USER
)
type (
	ConnxState int
	StateBase  struct {
		Uin                 string
		State               ConnxState
		DbUserLoginDeviceId uint64
	}
	IStreamNeedAuth interface {
		NeedAuth()bool
	}
	IStreamUserStatus interface {
		OnUserConnected(Uin  string, dbUserLoginDeviceId uint64)
		OnUserDisconnected(Uin  string, dbUserLoginDeviceId uint64)
	}
	LogicBase struct {
		hostInst    interface{}
		poolRWMutex sync.RWMutex
		pool        map[app.ServiceStream]*StateBase
		online      map[string][]app.ServiceStream
	}
)

func MakeLogicBase(hostInst interface{}) *LogicBase {
	return &LogicBase{
		hostInst: hostInst,
		pool : make(map[app.ServiceStream]*StateBase),
		online: make(map[string][]app.ServiceStream),
	}
}

func (l * LogicBase)GetHost() interface{} {
	return l.hostInst
}

func(l *LogicBase) Multicast(msgList []*comm.CSUserEvent){
	go func() {
		l.poolRWMutex.RLock()
		defer l.poolRWMutex.RUnlock()
		for _,msg := range msgList {
			ls := l.GetStreamListByUin(msg.Uin)
			for _,s := range ls {
				util.CheckError(s.Send(msg), "")
			}
		}
	}()
}

func (l *LogicBase) GetStreamListByUin(Uin string) []app.ServiceStream {
	return l.online[Uin]
}
func (l *LogicBase) UserIsOnline(Uin string) bool {
	return len(l.online[Uin]) > 0
}

func (l *LogicBase) SendUserMsg(Uin string, msg proto.Message) int {
	//get user stream
	var sendNum = 0
	for i,s := range l.online[Uin] {
		if util.CheckError(s.Send(msg),"send user:%s idx:%d error", Uin, i) == false {
			sendNum++
		}
	}
	return sendNum
}

func(l *LogicBase) Broadcast(msg *comm.CSUserEvent){
	log.Debugf("bradcast session param:%s", msg.String())
	go func() {
		l.poolRWMutex.RLock()
		defer l.poolRWMutex.RUnlock()
		for stream := range l.pool {
			util.CheckError(stream.Send(msg), "")
		}
	}()
}


func (l *LogicBase) OnOpen(s app.ServiceStream, c * gin.Context) bool {
	stateBase := StateBase {
		State: STREAM_CONNX_ANONYMOUSE,
		Uin:   "",
	}
	cUin,exist1 := c.Get("auth.UinSigner")
	_,exist2 := c.Get("db.UserTokenState")
	if exist1 && exist2 {
		stateBase.State = STREAM_CONNX_AUTHED_USER
		stateBase.Uin = cUin.(string)
		/*
		dbUserTokenState := tokenState.(*account.DbUserTokenState)
		stateBase.DbUserLoginDeviceId =dbUserTokenState.DbUserLoginDeviceId
		*/
		l.online[stateBase.Uin] = append(l.online[stateBase.Uin], s)
	} else {
		if l.hostInst != nil {
		  	if bi,ok := l.hostInst.(IStreamNeedAuth);ok && bi.NeedAuth() {
				return false
			}
		}
	}
	l.poolRWMutex.Lock()
	l.pool[s] = &stateBase
	l.poolRWMutex.Unlock()
	log.Debugf("stream:%s open with ctx:%v", s, stateBase)

	if l.hostInst != nil {
		if bi,ok := l.hostInst.(app.ServiceStreamLogic);ok {
			if !bi.OnOpen(s, c) {
				return false
			}
		}
		if stateBase.Uin != "" && stateBase.DbUserLoginDeviceId > 0 {
			if bi,ok := l.hostInst.(IStreamUserStatus); ok {
				bi.OnUserConnected(stateBase.Uin, stateBase.DbUserLoginDeviceId)
			}
		}
	}

	return true
}

func (l *LogicBase) OnMsg(s app.ServiceStream, request proto.Message){
	log.Debugf("stream:%v recv msg:%s", s, request.String())
	stateBase := l.pool[s]
	if stateBase == nil {
		log.Errorf(nil,"stream:%s on msg but stream State not exist ", s)
		return
	}
	if l.hostInst != nil {
		if bi, ok := l.hostInst.(app.ServiceStreamLogic); ok {
			bi.OnMsg(s, request)
		}
	}
}
func (l *LogicBase) OnClose(s app.ServiceStream, r app.ServiceStreamCloseReason){
	log.Debugf("stream:%s closed reason:%d", s, r)
	stateBase := l.pool[s]
	if l.hostInst != nil {
		if bi, ok := l.hostInst.(app.ServiceStreamLogic); ok {
			bi.OnClose(s, r)
		}

		if stateBase != nil && stateBase.Uin != "" && stateBase.DbUserLoginDeviceId > 0 {
			if bi,ok := l.hostInst.(IStreamUserStatus); ok {
				bi.OnUserDisconnected(stateBase.Uin, stateBase.DbUserLoginDeviceId)
			}
		}
	}
	if stateBase != nil {
		l.poolRWMutex.Lock()
		delete(l.pool, s)
		l.poolRWMutex.Unlock()
	}
}

