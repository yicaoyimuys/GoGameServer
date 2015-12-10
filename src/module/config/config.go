package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

import (
	"module"
	. "tools"
)

type ConfigModule struct {
	mutex sync.RWMutex

	dropData map[string]DropStruct
}

type DropStruct struct {
	Id   int    `json:"id"`
	Desc string `json:"desc"`
}

// 在初始化的时候将模块注册到module包
func init() {
	module.Config = ConfigModule{}
	module.Config.Load()
}

// 读取文件并且序列化
func readFile(fileName string, v interface{}) {
	path := os.Getenv("GOGAMESERVER_PATH") + "/data/jsons/" + fileName
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		ERR("读取配置错误", fileName)
		return
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		ERR("配置内容转换为Json错误", fileName)
		return
	}
}

// 加载所有配置文件
func (this ConfigModule) Load() {
	this.mutex.Lock()
	this.mutex.Unlock()

	this.loadDrop()
}

func (this ConfigModule) loadDrop() {
	readFile("drop.json", &this.dropData)
	//	DEBUG(this.dropData)
}
