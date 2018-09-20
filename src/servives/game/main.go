package main

import (
	"core/consts/service"
	. "core/libs"
	"core/service"
	_ "net/http/pprof"
)

func main() {
	////开启RpcServer
	//go startRpcServer()

	//初始化Service
	newService := service.NewService(Service.Game)
	newService.StartIpcServer()
	newService.StartRpcClient([]string{Service.Platform, Service.Ai})
	newService.StartDebug()

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initModule() {
}
