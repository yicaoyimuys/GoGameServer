package redisCaches

import (
	. "core/libs"
	"encoding/json"
	"servives/public/dbModels"
	"servives/public/redisInstances"
	"servives/public/redisKeys"
	"time"
)

//设置DBUser缓存
func SetUser(dbUser *dbModels.User) error {
	userKey := redisKeys.DbUser + NumToString(dbUser.Id)
	userData, _ := json.Marshal(dbUser)
	return redisInstances.User().Set(userKey, userData, time.Hour*24).Err()
}

//获取DBUser缓存
func GetUser(userId uint64) *dbModels.User {
	key := redisKeys.DbUser + NumToString(userId)
	val, err := redisInstances.User().Get(key).Result()
	if err != nil {
		return nil
	}

	var dbUser dbModels.User
	err = json.Unmarshal([]byte(val), &dbUser)
	return &dbUser
}
