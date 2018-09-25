package redisCaches

import (
	"core"
	. "core/libs"
	"encoding/json"
	"servives/public/dbModels"
	"time"
)

const (
	DB_User_Key = "DB_User_"
)

//设置DBUser缓存
func SetUser(dbUser *dbModels.User) error {
	redisClient := core.Service.GetRedisClient("user")

	userKey := DB_User_Key + NumToString(dbUser.Id)
	userData, _ := json.Marshal(dbUser)
	return redisClient.Set(userKey, userData, time.Hour*24).Err()
}

//获取DBUser缓存
func GetUser(userId uint64) *dbModels.User {
	redisClient := core.Service.GetRedisClient("user")

	key := DB_User_Key + NumToString(userId)
	val, err := redisClient.Get(key).Result()
	if err != nil {
		return nil
	}

	var dbUser dbModels.User
	err = json.Unmarshal([]byte(val), &dbUser)
	return &dbUser
}
