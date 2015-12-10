package main

import (
	"os"
)

import (
	"global"
	"proxys/redisProxy"
	"proxys/worldProxy"
	. "tools"
	"tools/cfg"
)

//各个模块
import (
	_ "module/cache"
	_ "module/config"
	_ "module/user"
	"proxys/logProxy"
)

var (
	world_ip   string
	world_port string
)

func main() {
	defer func() {
		if x := recover(); x != nil {
			ERR("caught panic in main()", x)
		}
	}()

	// 获取监听端口
	getPort()

	//启动
	global.Startup(global.ServerName, "world_log", nil)

	//连接Redis
	redisProxyErr := redisProxy.InitClient(cfg.GetValue("redis_ip"), cfg.GetValue("redis_port"), cfg.GetValue("redis_pwd"))
	checkError(redisProxyErr)

	//连接LogServer
	logProxyErr := logProxy.InitClient(cfg.GetValue("log_ip"), cfg.GetValue("log_port"))
	checkError(logProxyErr)

	//启动WorldServer
	worldProxyErr := worldProxy.InitServer(world_port)
	checkError(worldProxyErr)

	//保持进程
	global.Run()
}

func getPort() {
	world_ip = cfg.GetValue("world_ip")
	world_port = cfg.GetValue("world_port")
	global.ServerName = "WorldServer[" + world_port + "]"
}

func checkError(err error) {
	if err != nil {
		ERR("Fatal error: %v", err)
		os.Exit(-1)
	}
}
