package redisCaches

import (
	. "GoGameServer/core/libs"
	"GoGameServer/servives/public/mysqlModels"
	"GoGameServer/servives/public/redisInstances"
	"GoGameServer/servives/public/redisKeys"
	"encoding/json"
	"time"
)

//设置DBUser缓存
func SetUser(dbUser *mysqlModels.User) error {
	userKey := redisKeys.DbUser + NumToString(dbUser.Id)
	userData, _ := json.Marshal(dbUser)
	return redisInstances.User().Set(userKey, userData, time.Hour*24).Err()
}

//获取DBUser缓存
func GetUser(userId uint64) *mysqlModels.User {
	key := redisKeys.DbUser + NumToString(userId)
	val, err := redisInstances.User().Get(key).Result()
	if err != nil {
		return nil
	}

	var dbUser mysqlModels.User
	err = json.Unmarshal([]byte(val), &dbUser)
	return &dbUser
}
