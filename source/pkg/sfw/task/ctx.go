package task

import (
	"errors"
	"groot/pkg/broker"
)

var (
	MsgQueBroker broker.Broker
)

func Init(mqbroker string) error {
	brk2 := broker.GetBroker(mqbroker)
	if brk2 == nil {
		return errors.New("msg broker not found ")
	}

	MsgQueBroker = brk2

	return nil
}
