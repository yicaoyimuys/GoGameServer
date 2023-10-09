package service

import (
	"github.com/yicaoyimuys/GoGameServer/core/config"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/redis"
	"go.uber.org/zap"
)

func (this *Service) StartRedis() {
	this.redisClients = make(map[string]*redis.Client)

	redisConfigs := config.GetRedisConfig()
	for aliasName, redisConfig := range redisConfigs {
		client, err := redis.NewClient(redisConfig)
		CheckError(err)

		if client != nil {
			this.redisClients[aliasName] = client
			INFO("Redis连接成功", zap.String("AliasName", aliasName))
		}
	}
}

func (this *Service) GetRedisClient(redisAliasName string) *redis.Client {
	client, _ := this.redisClients[redisAliasName]
	return client
}
