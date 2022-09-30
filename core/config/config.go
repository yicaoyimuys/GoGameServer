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
	serviceConfig map[string]ServiceConfig
	redisConfig   map[string]RedisConfig
	logConfig     LogConfig
	mysqlConfig   map[string]MysqlConfig
	mongoConfig   map[string]MongoConfig
	lock          sync.Mutex
)

func Init(_env string) {
	env = _env
	load()
}

func load() {
	lock.Lock()
	loadConfig(&serviceConfig, "service.json")
	loadConfig(&redisConfig, "redis.json")
	loadConfig(&mysqlConfig, "mysql.json")
	loadConfig(&mongoConfig, "mongo.json")
	loadConfig(&logConfig, "log.json")
	lock.Unlock()
}

func getConfigPath(configFile string) string {
	return system.Root + "/config/" + env + "/" + configFile
}

func loadConfig(data interface{}, configName string) {
	configPath := getConfigPath(configName)
	fileData, err := ioutil.ReadFile(configPath)
	if err != nil {
		ERR("Config读取失败", configPath, err)
		return
	}
	json.Unmarshal(fileData, data)
}

func GetService(serviceName string) ServiceConfig {
	return serviceConfig[serviceName]
}

func GetRedisConfig() map[string]RedisConfig {
	return redisConfig
}

func GetMysqlConfig() map[string]MysqlConfig {
	return mysqlConfig
}

func GetLogConfig() LogConfig {
	return logConfig
}

func GetMongoConfig() map[string]MongoConfig {
	return mongoConfig
}
