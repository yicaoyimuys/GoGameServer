package main

import (
	"GoGameServer/core/consts/Service"
	. "GoGameServer/core/libs"
	"GoGameServer/core/service"
)

func main() {
	//初始化Service
	newService := service.NewService(Service.Log)
	newService.StartRpcServer()

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initModule() {

}
