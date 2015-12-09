package main

import (
	"os"
)

import (
	"global"
	"proxys/logProxy"
	. "tools"
	"tools/cfg"
)

var (
	log_ip string
	log_port string
)

func main()  {
	defer func() {
		if x := recover(); x != nil {
			ERR("caught panic in main()", x)
		}
	}()

	//获取监听端口
	getPort()

	//启动
	global.Startup(global.ServerName, "log_log", nil)

	//启动LogProxy
	err := logProxy.InitServer(log_port)
	checkError(err)

	//保持进程
	global.Run()
}

func getPort() {
	log_ip = cfg.GetValue("log_ip")
	log_port = cfg.GetValue("log_port")
	global.ServerName = "LogServer[" + log_port + "]"
}

func checkError(err error) {
	if err != nil {
		ERR("Fatal error: %v", err)
		os.Exit(-1)
	}
}