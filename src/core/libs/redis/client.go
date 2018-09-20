package redis

import (
	"core/libs/common"
	"core/libs/logger"
	"github.com/go-redis/redis"
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
			Addr:     redisConfig["host"].(string) + ":" + common.NumToString(redisConfig["port"]),
			Password: redisConfig["auth_pass"].(string),
			DB:       int(redisConfig["db"].(float64)),
		})

		pong, err := redisClient.Ping().Result()
		if err == nil {
			logger.Info("Redis_" + key + "连接成功..." + pong)
			redisClientList[key] = redisClient
		} else {
			logger.Error("Redis_"+key+"连接失败", err)
		}
	}
}
