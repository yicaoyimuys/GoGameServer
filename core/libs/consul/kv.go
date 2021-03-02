package consul

import (
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	//"GoGameServer/core/libs/logger"
)

var (
	kv          *api.KV
	useCache    bool
	caches      map[string]cacheValue
	cachesMutex sync.Mutex
)

type cacheValue struct {
	value string
	time  int64
}

func InitKV(cache bool) error {
	if kv != nil {
		return nil
	}

	useCache = cache
	if useCache {
		caches = make(map[string]cacheValue)
	}

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}

	kv = client.KV()
	return nil
}

func kv_getCache(key string) string {
	cachesMutex.Lock()
	defer cachesMutex.Unlock()

	if v, ok := caches[key]; ok {
		now := time.Now().Unix()
		//缓存10秒
		if now-v.time < 10 {
			return v.value
		}
	}
	return ""
}

func kv_setCache(key string, value string) {
	cachesMutex.Lock()
	defer cachesMutex.Unlock()

	caches[key] = cacheValue{value, time.Now().Unix()}
}

func KV_Get(key string) string {
	var value = ""
	if useCache {
		value = kv_getCache(key)
	}

	if len(value) == 0 {
		pair, _, err := kv.Get(key, nil)
		if err == nil && pair != nil {
			value = string(pair.Value)
			kv_setCache(key, value)
		} else {
			//logger.Debug("KV_Get", err, key+"不存在")
		}
	}
	return value
}

func KV_Set(key string, value string) error {
	pair := &api.KVPair{
		Key:   key,
		Value: []byte(value),
	}
	_, err := kv.Put(pair, nil)
	return err
}
