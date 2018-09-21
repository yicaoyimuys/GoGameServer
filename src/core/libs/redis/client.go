package redis

import (
	"core/libs/dict"
	"github.com/go-redis/redis"
)

type Client struct {
	*redis.Client
}

func NewClient(redisConfig map[string]interface{}) (*Client, error) {
	host := dict.GetString(redisConfig, "host")
	port := dict.GetString(redisConfig, "port")
	pass := dict.GetString(redisConfig, "auth_pass")
	db := dict.GetInt(redisConfig, "db")

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: pass,
		DB:       db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &Client{client}, nil
}
