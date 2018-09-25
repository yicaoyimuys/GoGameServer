package redisInstances

import (
	"core"
	"core/libs/redis"
)

func Global() *redis.Client {
	return core.Service.GetRedisClient("global")
}

func User() *redis.Client {
	return core.Service.GetRedisClient("user")
}
