package cache

import "errors"

type (
	Cache interface {
		Get(key string)interface{}
		Set(key string, v interface{}, timeoutSeconds int64)
	}
	CreateFunc func(string) Cache
)

var (
	mpCaches         map[string]Cache     = make(map[string]Cache)
	mpCacheTypeFunc  map[string]CreateFunc = make(map[string]CreateFunc)
	defaultCache     Cache
	defaultCacheName string
)

func RegisterCacheType(name string, fun CreateFunc) error {
	if mpCacheTypeFunc[name] != nil {
		return errors.New("cache type register already ")
	}
	mpCacheTypeFunc[name] = fun
	return nil
}

func AddCache(name string, backend string, connx string) error {
	bf := mpCacheTypeFunc[backend]
	if bf == nil {
		return errors.New("cache type is not registerd ")
	}
	b := bf(connx)
	if b == nil {
		return errors.New("cache create error !")
	}
	mpCaches[name] = b
	if defaultCache == nil {
		defaultCache = b
		defaultCacheName = name
	}
	return nil
}

func GetCache(name string) Cache {
	b, ok := mpCaches[name]
	if ok {
		return b
	}
	return nil
}
func GetCacheName(cache Cache) string {
	for name, v := range mpCaches {
		if v == cache {
			return name
		}
	}
	return "<nil>"
}

func GetDefaultCache() Cache {
	return defaultCache
}

func GetDefaultCacheName() string {
	return defaultCacheName
}
