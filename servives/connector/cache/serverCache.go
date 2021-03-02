package cache

import (
	"GoGameServer/servives/public/redisInstances"
	"GoGameServer/servives/public/redisKeys"
	"encoding/json"
)

func SetServerInfo(domainName string, serverPort string, onlineUsersNum int) {
	oldServerInfo := GetServerInfo(domainName, serverPort)

	//读取最高在线
	var oldMaxOnlineUsersNum = 0
	if oldServerInfo != nil {
		if num, exists := oldServerInfo["maxOnlineUsersNum"]; exists {
			oldMaxOnlineUsersNum = num
		}
	}

	redisKey := redisKeys.ServerInfo
	serverKey := domainName + ":" + serverPort

	serverInfo := make(map[string]int)
	serverInfo["onlineUsersNum"] = onlineUsersNum
	if onlineUsersNum > oldMaxOnlineUsersNum {
		serverInfo["maxOnlineUsersNum"] = onlineUsersNum
	} else {
		serverInfo["maxOnlineUsersNum"] = oldMaxOnlineUsersNum
	}

	byteData, _ := json.Marshal(serverInfo)
	redisInstances.Global().HSet(redisKey, serverKey, string(byteData))
}

func GetServerInfo(domainName string, serverPort string) map[string]int {
	redisKey := redisKeys.ServerInfo
	serverKey := domainName + ":" + serverPort
	val, err := redisInstances.Global().HGet(redisKey, serverKey).Result()
	if err != nil {
		return nil
	}

	var serverInfo map[string]int
	_ = json.Unmarshal([]byte(val), &serverInfo)
	return serverInfo
}
