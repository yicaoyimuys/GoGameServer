package redisProxy

import (
	"encoding/json"
	. "model"
	"strconv"
	"protos/dbProto"
	"protos"
//	. "tools"
)

const (
	DB_User_Key     = "DB_User_"
	DB_UserName_Key = "DB_UserName_"
	DB_UserLastLoginTime_Key = "DB_UserLastLoginTime_"
)

//设置DBUser缓存
func SetDBUser(dbUser *DBUserModel) {
	if dbUser == nil{
		return
	}
	userID := strconv.FormatUint(dbUser.ID, 10)

	userKey := DB_User_Key + userID
	userData, _ := json.Marshal(dbUser)

	userNameKey := DB_UserName_Key + dbUser.Name
	userNameData := []byte(userID)

	mapping := make(map[string][]byte)
	mapping[userKey] = userData
	mapping[userNameKey] = userNameData

	client.Mset(mapping)
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

//删除用户数据
func RemoveDBUser(userID uint64) {
	userIDStr := strconv.FormatUint(userID, 10)
	//删除用户数据
	client.Del(DB_User_Key + userIDStr)
	//删除最后登录数据
	client.Hdel(DB_UserLastLoginTime_Key, userIDStr)
}

//更新用户最后登录时间
func UpdateUserLastLoginTime(dbUser *DBUserModel) {
	userIDStr := strconv.FormatUint(dbUser.ID, 10)
	userLastLoginTime := strconv.FormatInt(dbUser.LastLoginTime, 10)

	//更新内存
	userKey := DB_User_Key + userIDStr
	data, _ := json.Marshal(dbUser)
	client.Set(userKey, data)

	//单独设置用户的最后登录时间，清理用户缓存数据使用
	client.Hset(DB_UserLastLoginTime_Key, userIDStr, []byte(userLastLoginTime))

	//更新DB
	msg := dbProto.MarshalProtoMsg(0, &dbProto.DB_User_UpdateLastLoginTimeC2S{
		UserID: protos.Uint64(dbUser.ID),
		Time:   protos.Int64(dbUser.LastLoginTime),
	})
	PushDBWriteMsg(msg)
}

//获取当前缓冲中的所有用户最后登录时间
func GetAllUserLastLoginTime() map[string]int64 {
	mapping := make(map[string]int64)
	client.Hgetall(DB_UserLastLoginTime_Key, mapping)
	return mapping
}