package cache

import (
	"github.com/garyburd/redigo/redis"
	"github.com/lexkong/log"
	"groot/pkg/util"
)

type (
	RedisCache struct {
		ConnxEnv string
		Pool     *redis.Pool
	}
)

func MakeRedisCache(connx string) Cache {
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
	redisBroker := RedisCache{
		ConnxEnv: connx,
		Pool:     &pool,
	}
	return &redisBroker
}

func (r *RedisCache)Set(key string, v interface{}, timeoutSeconds int64) {
	c := r.Pool.Get()
	if util.CheckError(c.Err(), "get redis cache fail for set key", key) {
		return
	}
	//c.Do("")
}


func (r *RedisCache)Get(key string) interface{} {
	c := r.Pool.Get()
	if util.CheckError(c.Err(), "get redis cache fail for get key:%s", key) {
		return nil
	}
	//todo
	return nil
}


func init() {
	RegisterCacheType("redis", MakeRedisCache)
}
