package config

import (
	. "GoGameServer/core/libs"
	"GoGameServer/core/libs/system"
	"encoding/json"
	"io/ioutil"
	"sync"
)

var (
	env           string
	serviceConfig map[string]interface{}
	redisConfig   map[string]interface{}
	logConfig     map[string]interface{}
	mysqlConfig   map[string]interface{}
	mongoConfig   map[string]interface{}
	lock          sync.Mutex
)

func Init(_env string) {
	env = _env
	load()
}

func load() {
	var serviceConfigPath = getConfigPath("service.json")
	var redisConfigPath = getConfigPath("redis.json")
	var mysqlConfigPath = getConfigPath("mysql.json")
	var mongoConfigPath = getConfigPath("mongo.json")
	var logConfigPath = getConfigPath("log.json")

	lock.Lock()
	loadConfig(&serviceConfig, serviceConfigPath)
	loadConfig(&redisConfig, redisConfigPath)
	loadConfig(&mysqlConfig, mysqlConfigPath)
	loadConfig(&mongoConfig, mongoConfigPath)
	loadConfig(&logConfig, logConfigPath)
	lock.Unlock()
}

func getConfigPath(configFile string) string {
	return system.Root + "/config/" + env + "/" + configFile
}

func loadConfig(data *map[string]interface{}, configPath string) {
	fileData, _ := ioutil.ReadFile(configPath)
	json.Unmarshal(fileData, data)
}

func GetConnectorService(serviceId int) map[string]interface{} {
	serviceData := serviceConfig["connector"].(map[string]interface{})
	serverDatas := serviceData["services"].(map[string]interface{})
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

func GetApiService(serviceId int) map[string]interface{} {
	serviceData := serviceConfig["api"].(map[string]interface{})
	serverDatas := serviceData["services"].(map[string]interface{})
	return serverDatas[NumToString(serviceId)].(map[string]interface{})
}

func GetApiServiceTslCrt() string {
	serverData := serviceConfig["api"].(map[string]interface{})
	return serverData["tslCrt"].(string)
}

func GetApiServiceTslKey() string {
	serverData := serviceConfig["api"].(map[string]interface{})
	return serverData["tslKey"].(string)
}

func GetRedisConfig() map[string]interface{} {
	return redisConfig
}

func GetMysqlConfig() map[string]interface{} {
	return mysqlConfig
}

func GetLogConfig() map[string]interface{} {
	return logConfig
}

func GetMongoConfig() map[string]interface{} {
	return mongoConfig
}
