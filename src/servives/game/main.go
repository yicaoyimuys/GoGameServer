package main

import (
	"core/consts/service"
	. "core/libs"
	"core/service"
	_ "net/http/pprof"
	"servives/game/module"
)

func main() {
	////开启RpcServer
	//go startRpcServer()

	//初始化Service
	newService := service.NewService(Service.Game)
	newService.StartIpcServer()
	newService.StartRpcServer(&module.GameRpcServer{})
	newService.StartRpcClient([]string{Service.Platform, Service.Log})
	newService.StartDebug()

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initModule() {
}
