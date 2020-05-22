package broker

import "errors"

type (
	Broker interface {
		Publish(topic string, msg []byte) error
		Subscribe(topic string) <-chan []byte
		Push(topic string, msg ...[]byte) error
		Pull(topic string) <-chan []byte
	}
	CreateFunc func(string) Broker
)

var (
	mpBrokers         map[string]Broker     = make(map[string]Broker)
	mpBrokerTypeFunc  map[string]CreateFunc = make(map[string]CreateFunc)
	defaultBroker     Broker
	defaultBrokerName string
)

func RegisterBrokerType(name string, fun CreateFunc) error {
	if mpBrokerTypeFunc[name] != nil {
		return errors.New("broker type register already ")
	}
	mpBrokerTypeFunc[name] = fun
	return nil
}

func AddBroker(name string, backend string, connx string) error {
	bf := mpBrokerTypeFunc[backend]
	if bf == nil {
		return errors.New("broker type is not registerd ")
	}
	b := bf(connx)
	if b == nil {
		return errors.New("broker create error !")
	}
	mpBrokers[name] = b
	if defaultBroker == nil {
		defaultBroker = b
		defaultBrokerName = name
	}
	return nil
}

func GetBroker(name string) Broker {
	b, ok := mpBrokers[name]
	if ok {
		return b
	}
	return nil
}
func GetBrokerName(broker Broker) string {
	for name, v := range mpBrokers {
		if v == broker {
			return name
		}
	}
	return "<nil>"
}

func GetDefaultBroker() Broker {
	return defaultBroker
}

func GetDefaultBrokerName() string {
	return defaultBrokerName
}
