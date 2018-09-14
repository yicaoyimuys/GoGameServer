package config

import (
	. "core/libs"
	"core/libs/cfg"
	"encoding/json"
	"global"
	"io/ioutil"
	"sync"
)

var (
	serverConfig map[string]interface{}
	redisConfig  map[string]interface{}
	logConfig    map[string]interface{}
	lock         sync.Mutex
)

func init() {
	var serverConfigPath = cfg.ROOT + "/config/server.json"
	var redisConfigPath = cfg.ROOT + "/config/redis.json"
	var logConfigPath = cfg.ROOT + "/config/log.json"

	lock.Lock()
	loadConfig(&serverConfig, serverConfigPath)
	loadConfig(&redisConfig, redisConfigPath)
	loadConfig(&logConfig, logConfigPath)
	lock.Unlock()
}

func loadConfig(data *map[string]interface{}, configPath string) {
	fileData, _ := ioutil.ReadFile(configPath)
	json.Unmarshal(fileData, data)
}

func GetConnectorServer(serverId int) map[string]interface{} {
	serverData := serverConfig[global.Env].(map[string]interface{})
	gameData := serverData["connector"].(map[string]interface{})
	return gameData[NumToString(serverId)].(map[string]interface{})
}

func GetGameServerTslCrt() string {
	serverData := serverConfig[global.Env].(map[string]interface{})
	return serverData["tslCrt"].(string)
}

func GetGameServerTslKey() string {
	serverData := serverConfig[global.Env].(map[string]interface{})
	return serverData["tslKey"].(string)
}

func GetRedisList() map[string]interface{} {
	return redisConfig[global.Env].(map[string]interface{})
}

func GetLog() map[string]interface{} {
	return logConfig[global.Env].(map[string]interface{})
}
