package main

import (
	"core/consts"
	. "core/libs"
	"core/service"
	"game/module"
	_ "net/http/pprof"
)

func main() {
	//初始化
	newService := service.NewService(consts.Service_Connector)
	//开启Redis
	newService.StartRedis()
	//开启WebSocket
	newService.StartWebSocket()
	//开启Ipc
	newService.StartIpcClient([]string{consts.Service_Game, consts.Service_Matching})

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initModule() {
	module.StartServerTimer()
}
