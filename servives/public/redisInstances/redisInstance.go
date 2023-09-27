package redisInstances

import (
	"github.com/yicaoyimuys/GoGameServer/core"
	"github.com/yicaoyimuys/GoGameServer/core/libs/redis"
)

func Global() *redis.Client {
	return core.Service.GetRedisClient("global")
}

func User() *redis.Client {
	return core.Service.GetRedisClient("user")
}
