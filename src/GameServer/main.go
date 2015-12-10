package main

import (
	"flag"
	"os"
	"strconv"
)

import (
	"global"
	"proxys/redisProxy"
	"proxys/transferProxy"
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
	game_ip   string
	game_port string
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
	global.Startup(global.ServerName, "game_log", nil)

	//连接TransferServer
	err := transferProxy.InitClient(cfg.GetValue("transfer_ip"), cfg.GetValue("transfer_port"))
	checkError(err)

	//连接WorldServer
	worldProxyErr := worldProxy.InitClient(cfg.GetValue("world_ip"), cfg.GetValue("world_port"))
	checkError(worldProxyErr)

	//连接Redis
	redisProxyErr := redisProxy.InitClient(cfg.GetValue("redis_ip"), cfg.GetValue("redis_port"), cfg.GetValue("redis_pwd"))
	checkError(redisProxyErr)

	//连接LogServer
	logProxyErr := logProxy.InitClient(cfg.GetValue("log_ip"), cfg.GetValue("log_port"))
	checkError(logProxyErr)

	//保持进程
	global.Run()
}

func getPort() {
	var s int
	flag.IntVar(&s, "s", 0, "tcp listen port")
	flag.Parse()
	if s == 0 {
		ERR("please set gameserver port")
		os.Exit(-1)
	}
	game_ip = cfg.GetValue("game_ip_" + strconv.Itoa(s))
	game_port = cfg.GetValue("game_port_" + strconv.Itoa(s))
	global.ServerName = "GameServer[" + game_port + "]"
	global.ServerID = uint32(s)
}

func checkError(err error) {
	if err != nil {
		ERR("Fatal error: %v", err)
		os.Exit(-1)
	}
}
