package cache

import (
	"core"
	"encoding/json"
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
	redisClient := core.Service.GetRedisClient("global")
	redisClient.HSet(ServerInfo_KEY, serverKey, string(byteData))
}

func GetServerInfo(domainName string, serverPort string) map[string]int {
	serverKey := domainName + ":" + serverPort
	redisClient := core.Service.GetRedisClient("global")
	val, err := redisClient.HGet(ServerInfo_KEY, serverKey).Result()
	if err != nil {
		return nil
	}

	var serverInfo map[string]int
	_ = json.Unmarshal([]byte(val), &serverInfo)
	return serverInfo
}
