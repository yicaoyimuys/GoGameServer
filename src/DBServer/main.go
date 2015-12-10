package main

import (
	"os"
)

import (
	"global"
	"proxys/dbProxy"
	"proxys/redisProxy"
	. "tools"
	"tools/cfg"
)

//各个模块
import (
	_ "module/cache"
	_ "module/config"
	_ "module/user"
)

var (
	db_ip   string
	db_port string
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
	global.Startup(global.ServerName, "db_log", stopDBServer)

	//连接Redis
	redisProxyErr := redisProxy.InitClient(cfg.GetValue("redis_ip"), cfg.GetValue("redis_port"), cfg.GetValue("redis_pwd"))
	checkError(redisProxyErr)

	//开启DBServer监听
	err := dbProxy.InitServer(db_port)
	checkError(err)

	//保持进程
	global.Run()
}

func getPort() {
	db_ip = cfg.GetValue("db_ip")
	db_port = cfg.GetValue("db_port")
	global.ServerName = "DBServer[" + db_port + "]"
}

func stopDBServer() {
	INFO("Waiting SyncDB...")
	dbProxy.SyncDB()
	INFO("SyncDB Success")
}

func checkError(err error) {
	if err != nil {
		ERR("Fatal error: %v", err)
		os.Exit(-1)
	}
}
