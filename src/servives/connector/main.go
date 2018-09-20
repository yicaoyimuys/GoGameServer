package main

import (
	"core/consts/service"
	. "core/libs"
	"core/service"
	_ "net/http/pprof"
	"servives/connector/module"
)

func main() {
	//初始化Service
	newService := service.NewService(Service.Connector)
	newService.StartRedis()
	newService.StartWebSocket()
	newService.StartIpcClient([]string{Service.Game, Service.Matching})

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initModule() {
	module.StartServerTimer()
}
