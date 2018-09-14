package redis

import (
	"github.com/go-redis/redis"
	. "tools"
)

var redisClientList = make(map[string]*redis.Client)

func GetLink(key string) *redis.Client {
	redisClient, _ := redisClientList[key]
	redisClient.Ping()
	return redisClient
}

func InitRedis(redisListConfig map[string]interface{}) {
	for key, data := range redisListConfig {
		redisConfig := data.(map[string]interface{})
		redisClient := redis.NewClient(&redis.Options{
			Addr:     redisConfig["host"].(string) + ":" + NumToString(redisConfig["port"]),
			Password: redisConfig["auth_pass"].(string),
			DB:       int(redisConfig["db"].(float64)),
		})

		pong, err := redisClient.Ping().Result()
		if err == nil {
			INFO("Redis_" + key + "连接成功..." + pong)
			redisClientList[key] = redisClient
		} else {
			ERR("Redis_"+key+"连接失败", err)
		}
	}
}
