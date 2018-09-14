package cache

import (
	"encoding/json"
	"tools/redis"
)

const ServerInfo_KEY = "BattleServer_ServerInfo"

func SetServerInfo(domainName string, serverPort string, onlineUsersNum int) {
	oldServerInfo := GetServerInfo(domainName, serverPort)

	//读取最高在线
	var oldMaxOnlineUsersNum = 0
	if oldServerInfo != nil {
		if num, exists := oldServerInfo["maxOnlineUsersNum"]; exists {
			oldMaxOnlineUsersNum = num
		}
	}

	serverKey := domainName + ":" + serverPort

	serverInfo := make(map[string]int)
	serverInfo["onlineUsersNum"] = onlineUsersNum
	if onlineUsersNum > oldMaxOnlineUsersNum {
		serverInfo["maxOnlineUsersNum"] = onlineUsersNum
	} else {
		serverInfo["maxOnlineUsersNum"] = oldMaxOnlineUsersNum
	}

	byteData, _ := json.Marshal(serverInfo)
	redis.GetLink("main").HSet(ServerInfo_KEY, serverKey, string(byteData))
}

func GetServerInfo(domainName string, serverPort string) map[string]int {
	serverKey := domainName + ":" + serverPort
	val, err := redis.GetLink("main").HGet(ServerInfo_KEY, serverKey).Result()
	if err != nil {
		return nil
	}

	var serverInfo map[string]int
	_ = json.Unmarshal([]byte(val), &serverInfo)
	return serverInfo
}

//func GetUserCache(userId int64) map[string]interface{} {
//	key := USER_KEY + NumToString(userId)
//	val, err := global.RedisClient.Get(key).Result()
//
//	if err != nil {
//		return nil
//	}
//
//	var dbUser interface{}
//	err = json.Unmarshal([]byte(val), &dbUser)
//	return dbUser.(map[string]interface{})
//}
