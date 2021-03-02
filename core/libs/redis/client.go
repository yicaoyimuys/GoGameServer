package redis

import (
	"GoGameServer/core/libs/dict"
	"time"

	"github.com/go-redis/redis"
)

type Client struct {
	redisClient *redis.Client
	prefix      string
}

func NewClient(redisConfig map[string]interface{}) (*Client, error) {
	prefix := dict.GetString(redisConfig, "prefix")
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

	return &Client{
		redisClient: client,
		prefix:      prefix,
	}, nil
}

func (this *Client) GetKey(key string) string {
	return this.prefix + "." + key
}

func (this *Client) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	key = this.GetKey(key)
	return this.redisClient.Set(key, value, expiration)
}

func (this *Client) Get(key string) *redis.StringCmd {
	key = this.GetKey(key)
	return this.redisClient.Get(key)
}

func (this *Client) HSet(key, field string, value interface{}) *redis.BoolCmd {
	key = this.GetKey(key)
	return this.redisClient.HSet(key, field, value)
}

func (this *Client) HGet(key, field string) *redis.StringCmd {
	key = this.GetKey(key)
	return this.redisClient.HGet(key, field)
}

func (this *Client) HGetAll(key string) *redis.StringStringMapCmd {
	key = this.GetKey(key)
	return this.redisClient.HGetAll(key)
}
