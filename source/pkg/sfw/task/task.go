package task

import (
	"github.com/gogo/protobuf/proto"
	"github.com/lexkong/log"
	"groot/comm/constk"
	"groot/proto/ssmsg"
	"groot/pkg/util"
)

type TaskParamType proto.Message

func Push(topic string, msg []byte) error {
	return MsgQueBroker.Push(topic, msg)
}

func Pull(topic string) <-chan []byte {
	return MsgQueBroker.Pull(topic)
}

func AddTask(topic string, tp TaskParamType) error {
	b, e := proto.Marshal(tp)
	if e != nil {
		log.Errorf(e, "task topic:%s marshal error")
		return e
	}
	log.Debugf("add task topic:%s param len:%d", topic, len(b))
	e = MsgQueBroker.Push(topic, b)
	if e != nil {
		log.Errorf(e, "broker push error ")
		return e
	}
	return nil
}

func TcsCall(service,method string, param proto.Message){
	var bParam []byte
	if param != nil {
		b, e:= proto.Marshal(param)
		if e != nil {
			log.Errorf(e, "call service:%s.%s param error", service, method)
			return
		}
		bParam = b
	}
	msg := & ssmsg.TaskCallServiceMsg {
		Service: service,
		Method: method,
		Param: bParam,
	}
	util.CheckError(AddTask(util.GetBrokerTopicName(constk.BRK_TOPIC_CALL_SERVICE, service), msg),
	"add service task:%s.%s fail", service, method)
}
