package broker

import (
	"github.com/garyburd/redigo/redis"
	"github.com/lexkong/log"
	"reflect"
)

type (
	RedisBroker struct {
		ConnxEnv string
		Pool     *redis.Pool
	}
)

func MakeRedisBroker(connx string) Broker {
	pool := redis.Pool{
		MaxIdle:     8,
		IdleTimeout: 100,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", connx)
			if err != nil {
				log.Errorf(err, "redis dial error !")
				return nil, err
			}
			return c, err
		},
	}
	redisBroker := RedisBroker{
		ConnxEnv: connx,
		Pool:     &pool,
	}
	return &redisBroker
}

func (b *RedisBroker) Publish(topic string, buff []byte) error {
	conn := b.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("Publish", topic, buff)
	if err != nil {
		log.Errorf(err, "redis publish message error !")
		return err
	}
	return nil
}

func (b *RedisBroker) Subscribe(topic string) <-chan []byte {
	pipe := make(chan []byte, 128)
	go func() {
		conn := b.Pool.Get()
		defer conn.Close()
		defer close(pipe)
		psc := redis.PubSubConn{conn}
		_ = psc.Subscribe(topic)
		log.Debugf("redis subscribe topic:%s", topic)
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				//log.Debugf("channel:%s received message:%v", v.Channel, v.Data)
				pipe <- v.Data
			case redis.Subscription:
				log.Infof("subscription channel:%s kind:%s count:%d", v.Channel, v.Kind, v.Count)
			case error:
				log.Errorf(v, "subscription topic:%s get error", topic)
				return
			}
		}
	}()
	return pipe
}

func (b *RedisBroker) Push(topic string, msgs ...[]byte) error {
	if len(msgs) == 0 {
		return nil
	}
	conn := b.Pool.Get()
	defer conn.Close()
	var args []interface{}
	args = append(args, topic)
	for i, _ := range msgs {
		args = append(args, msgs[i])
	}
	_, err := conn.Do("RPUSH", args...)
	if err != nil {
		log.Errorf(err, "redis push message of topic:%s error !", topic)
		return err
	}

	return nil
}
func (b *RedisBroker) Pull(topic string) <-chan []byte {
	//todo wiht lpop
	pipe := make(chan []byte, 128)
	go func() {
		defer close(pipe)
		conn := b.Pool.Get()
		defer conn.Close()
		for {
			//blpop 1
			reply, err := conn.Do("lpop", topic)
			if err != nil {
				log.Errorf(err, "pop topic:%s error", topic)
			} else {
				if reply != nil {
					ret, ok := reply.([]byte)
					if !ok {
						log.Warnf("pop topic:%s replay res conv(type:%v) error !", topic, reflect.TypeOf(reply))
					} else {
						log.Debugf("pop topic:%s replay res bytes len:%d ", topic, len(ret))
						pipe <- ret
					}
				}
			}
		}
	}()
	return pipe
}

func init() {
	RegisterBrokerType("redis", MakeRedisBroker)
}
