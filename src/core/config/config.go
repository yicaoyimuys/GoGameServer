package config

import (
	"core/argv"
	. "core/libs"
	"core/libs/cfg"
	"encoding/json"
	"io/ioutil"
	"sync"
)

var (
	serviceConfig map[string]interface{}
	redisConfig   map[string]interface{}
	logConfig     map[string]interface{}
	lock          sync.Mutex
)

func Init() {
	var serviceConfigPath = getConfigPath("service.json")
	var redisConfigPath = getConfigPath("redis.json")
	var logConfigPath = getConfigPath("log.json")

	lock.Lock()
	loadConfig(&serviceConfig, serviceConfigPath)
	loadConfig(&redisConfig, redisConfigPath)
	loadConfig(&logConfig, logConfigPath)
	lock.Unlock()
}

func getConfigPath(configFile string) string {
	return cfg.ROOT + "/config/" + argv.Values.Env + "/" + configFile
}

func loadConfig(data *map[string]interface{}, configPath string) {
	fileData, _ := ioutil.ReadFile(configPath)
	json.Unmarshal(fileData, data)
}

func GetConnectorService(serviceId int) map[string]interface{} {
	serviceData := serviceConfig["connector"].(map[string]interface{})
	serverDatas := serviceData["servers"].(map[string]interface{})
	return serverDatas[NumToString(serviceId)].(map[string]interface{})
}

func GetConnectorServiceTslCrt() string {
	serverData := serviceConfig["connector"].(map[string]interface{})
	return serverData["tslCrt"].(string)
}

func GetConnectorServiceTslKey() string {
	serverData := serviceConfig["connector"].(map[string]interface{})
	return serverData["tslKey"].(string)
}

func GetRedisList() map[string]interface{} {
	return redisConfig
}

func GetLog() map[string]interface{} {
	return logConfig
}
