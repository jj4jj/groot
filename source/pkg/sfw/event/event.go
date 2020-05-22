package sfw_event

import (
	"github.com/gogo/protobuf/proto"
	"github.com/lexkong/log"
	"groot/comm/constk"
	"groot/proto/ssmsg"
	"groot/pkg/broker"
	"groot/pkg/util"
	"reflect"
	"sync"
)

type (
	EventHandler   interface{} //real proto type:func(evt ssmsg.SvcEventType, param EventParamType)
	EvtTopicCtx    struct {
		evt      	ssmsg.SvcEventType
		topic    	string
		paramType  	reflect.Type
		param		[]reflect.Value //EventParamType
		ch       	<-chan []byte
		handlers []reflect.Value
	}
	EvtTopicFrame struct {
		ctx    *EvtTopicCtx
		buffer []byte
	}
)

var (
	eventBroker    broker.Broker
	topicListenerLock	sync.RWMutex
	topicListeners map[string]*EvtTopicCtx
)

func init() {
	topicListeners = make(map[string]*EvtTopicCtx)
}

func SetEventBroker(name string) {
	eventBroker = broker.GetBroker(name)
}

//go routine num
func StartEventDispatcher(grNum int) {
	if grNum <= 0 {
		grNum = 1
	}
	log.Infof("start event dispatcher num:%d with broker name:%s", grNum, broker.GetBrokerName(eventBroker))
	muxChannel := make(chan EvtTopicFrame, 128)
	for topic, ctx := range topicListeners {
		ctx.ch = eventBroker.Subscribe(topic)
		log.Debugf("event broker listen topic:%s for event:%d ...", topic, ctx.evt)
		for j:=0;j<grNum;j++ {
			ctx.param = append(ctx.param, reflect.New(ctx.paramType))
		}
		go func(c *EvtTopicCtx) {
			for b := range c.ch {
				muxChannel <- EvtTopicFrame{
					ctx:    c,
					buffer: b,
				}
			}
		}(ctx)
	}
	var wg sync.WaitGroup
	wg.Add(grNum)
	//merge channels
	for i := 0; i < grNum; i++ {
		go func(index int) {
			defer wg.Done()
			for frame := range muxChannel {
				log.Debugf("event recv ev ctx frame topic:%s evt:%d", frame.ctx.topic, frame.ctx.evt)
				evparm := frame.ctx.param[index]
				if err := proto.Unmarshal(frame.buffer, evparm.Interface().(proto.Message)); err != nil {
					log.Errorf(err, "evt topic:%s recv msg len:%d parse error ", frame.ctx.topic, len(frame.buffer))
					continue
				}
				in := []reflect.Value{reflect.ValueOf(frame.ctx.evt), reflect.ValueOf(evparm)}
				for _, handler := range frame.ctx.handlers {
					handler.Call(in)
				}
			}
		}(i)
	}
	wg.Wait()
}

func FireEvent(evt ssmsg.SvcEventType, evparam proto.Message) {
	bf, err := proto.Marshal(evparam)
	if err != nil {
		log.Errorf(err, "fire event:%d marshal error", evt)
		return
	}
	//get topic
	topic := util.GetBrokerTopicName(constk.BRK_TOPIC_SERVICE_EVENT, evt)
	err = eventBroker.Publish(topic, bf)
	if err != nil {
		log.Errorf(err, "publish topic:%s error !", topic)
	}
	log.Debugf("fire event topic:%s param len:%d success", topic, len(bf))
}

func checkEventHandlerFunc(handler interface{}) bool {
	hv := reflect.ValueOf(handler)
	if hv.Type().Kind() != reflect.Func {
		return false
	}

	if hv.Type().NumIn() != 2 {
		return false
	}

	if hv.Type().In(0) != reflect.TypeOf(ssmsg.SvcEventType_SVC_EVT_NONE) {
		return false
	}

	protoMessageInterfaceType := reflect.TypeOf(new(proto.Message)).Elem()
	if !hv.Type().In(1).Implements(protoMessageInterfaceType) {
		return false
	}
	return true
}

//handler will run another goroutine (make sure concurrency processing)
func ListenEvent(evt ssmsg.SvcEventType, handler EventHandler) {
	if handler == nil {
		return
	}
	if !checkEventHandlerFunc(handler) {
		log.Errorf(nil, "evt:%d handler func is invalid:f(ssmsg.SvcEventType,proto.Message)", evt)
		return
	}
	topic := util.GetBrokerTopicName(constk.BRK_TOPIC_SERVICE_EVENT, evt)
	log.Debugf("event listen type:%d topic:%s", evt, topic)
	hv := reflect.ValueOf(handler)
	topicListenerLock.Lock()
	defer topicListenerLock.Unlock()
	////////////////////////////////
	ctx, ok := topicListeners[topic]
	if ok {
		if hv.Type().In(1).Elem() == ctx.paramType  {
			log.Errorf(nil, "add handler but param 2 type not match")
			return
		}
		ctx.handlers = append(ctx.handlers, hv)
	} else {
		topicListeners[topic] = &EvtTopicCtx{
			ch:       nil,
			evt:      evt,
			paramType: hv.Type().In(1).Elem(),
			topic:    topic,
			handlers: []reflect.Value{hv},
		}
	}
}
