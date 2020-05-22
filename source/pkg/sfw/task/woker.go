package task

import (
	"github.com/gogo/protobuf/proto"
	"github.com/lexkong/log"
	"reflect"
	"sync"
)


func RunWorker(topic string, workerNum int, handler interface{}) {
	if workerNum <= 0 {
		return
	}
	var wg sync.WaitGroup
	ch := Pull(topic)
	wg.Add(workerNum)
	///
	caller := reflect.ValueOf(handler)
	handlerType := caller.Type()
	protoMessageInterfaceType := reflect.ValueOf(new(proto.Message)).Type().Elem()
	if !(handlerType.Kind() == reflect.Func && handlerType.NumIn() == 1 && handlerType.In(0).Implements(protoMessageInterfaceType)) {
		log.Errorf(nil, "worker handler prototype error need [f(proto.Message)]")
		return
	}
	valueType := handlerType.In(0).Elem()
	for i := 0; i < workerNum; i++ {
		go func() {
			defer wg.Done()
			msg := reflect.New(valueType)
			protoMsg := msg.Interface().(proto.Message)
			for v := range ch {
				if e := proto.Unmarshal(v, protoMsg); e != nil {
					log.Errorf(e, "worker:%s param len:%d unmarshal error !", topic, len(v))
				} else {
					//type WorkerHandlerFunc func(message proto.Message)
					caller.Call([]reflect.Value{msg})
				}
			}
		}()
	}
	wg.Wait()
}
