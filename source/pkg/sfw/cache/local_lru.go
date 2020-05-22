package cache

import (
	"github.com/garyburd/redigo/redis"
	"github.com/lexkong/log"
	"groot/pkg/util"
	"github.com/hashicorp/golang-lru"
	"strconv"
	"strings"
)


type (
	LocalLruCache struct {
		cache *lru.Cache
	}
)

//max=x
func MakeLocalLruCache(connx string) Cache {
	as := strings.Split(connx, ";")
	maxCap := 1024*1024
	for i := range as {
		vs := strings.Split(as[i], "=")
		if len(vs) == 2 && strings.ToLower(vs[0]) == "max" {
			n,e := strconv.Atoi(vs[1])
			if e == nil {
				maxCap = n
			}
		}
	}

	lruCache,e := lru.New(maxCap)
	if e == nil {
		return &LocalLruCache{
			cache: lruCache,
		}
	} else {
		log.Errorf(e, "create lru cache fail")
	}
	return nil
}

func (r *LocalLruCache)Set(key string, v interface{}, timeoutSeconds int64) {
	//c.Do("")
	r.cache.Add(key, v)
}


func (r *LocalLruCache)Get(key string) interface{} {
	//todo
	if v,ok := r.cache.Get(key); ok {
		return v
	}
	return nil
}


func init() {
	RegisterCacheType("local_lru", MakeLocalLruCache)
}
