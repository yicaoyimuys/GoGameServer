package main

import (
	"os"
)

import (
	"global"
	"proxys/transferProxy"
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
	transfer_ip   string
	transfer_port string
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
	global.Startup(global.ServerName, "transfer_log", nil)

	//开启TransferServer
	err := transferProxy.InitServer(transfer_port)
	checkError(err)

	//保持进程
	global.Run()
}

func getPort() {
	transfer_ip = cfg.GetValue("transfer_ip")
	transfer_port = cfg.GetValue("transfer_port")
	global.ServerName = "TransferServer[" + transfer_port + "]"
}

func checkError(err error) {
	if err != nil {
		ERR("Fatal error: %v", err)
		os.Exit(-1)
	}
}
