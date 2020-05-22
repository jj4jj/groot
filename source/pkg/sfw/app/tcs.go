package app

import (
	"github.com/gogo/protobuf/proto"
	"github.com/lexkong/log"
	"groot/proto/ssmsg"
	"reflect"
)

type (
	ServiceMethodHandler interface{} //func(proto.Message)
	ServiceTcsMethodCtx struct {
		name       string
		handler    ServiceMethodHandler
		reqValType reflect.Type
		caller     reflect.Value
	}
)


func checkServiceTcsMethodHandler(handler ServiceMethodHandler) bool {
	ht := reflect.ValueOf(handler).Type()
	protoMessageInterfaceType := reflect.ValueOf(new(proto.Message)).Type().Elem()
	if !(ht.Kind() == reflect.Func && ht.NumIn() == 1 && ht.In(0).Implements(protoMessageInterfaceType)) {
		log.Errorf(nil, "tcs method handler prototype error need [f(proto.Message)]")
		return false
	}
	return true
}

func createTcsWorkerHandler(mp map[string]*ServiceTcsMethodCtx) func(*ssmsg.TaskCallServiceMsg) {
	return func(msg *ssmsg.TaskCallServiceMsg){
		c,ok := mp[msg.Method]
		if c == nil || !ok {
			log.Errorf(nil, "not found tcs:%s.%s method ctx", msg.Service, msg.Method)
			return
		}
		reqParam := reflect.New(c.reqValType)
		if e := proto.Unmarshal(msg.Param, reqParam.Interface().(proto.Message)); e != nil {
			log.Errorf(nil, "tcs:%s.%s unpack msg fail", msg.Service, msg.Method)
			return
		}
		in := []reflect.Value{reqParam}
		c.caller.Call(in)
	}
}



