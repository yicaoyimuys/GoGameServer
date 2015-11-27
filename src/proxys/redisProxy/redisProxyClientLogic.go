package redisProxy

import (
	"encoding/json"
	. "model"
	"strconv"
)

const DB_User_Key string = "DB_User"
const DB_UserName_Key string = "DB_UserName"

//设置DBUser缓存
func SetDBUser(dbUser *DBUserModel) {
	userID := strconv.FormatUint(dbUser.ID, 10)

	userKey := DB_User_Key + userID
	data, _ := json.Marshal(dbUser)
	client.Set(userKey, data)

	userNameKey := DB_UserName_Key + dbUser.Name
	client.Set(userNameKey, []byte(userID))
}

//根据UserID获取用户DB数据
func GetDBUser(userID uint64) *DBUserModel {
	key := DB_User_Key + strconv.FormatUint(userID, 10)
	data, err := client.Get(key)
	if err != nil {
		return nil
	}
	var dbUser *DBUserModel = NewDBUserModel()
	json.Unmarshal(data, dbUser)
	return dbUser
}

//根据UserName获取用户DB数据
func GetDBUserByUserName(userName string) *DBUserModel {
	userNameKey := DB_UserName_Key + userName
	data, err := client.Get(userNameKey)
	if err != nil {
		return nil
	}
	userID, err := strconv.ParseUint(string(data), 10, 64)
	return GetDBUser(userID)
}

//更新用户最后登录时间
func UpdateUserLastLoginTime(userID uint64, time int64) {
	var dbUser *DBUserModel = GetDBUser(userID)
	if dbUser == nil {
		return
	}
	dbUser.LastLoginTime = time
	SetDBUser(dbUser)
}
