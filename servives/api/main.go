package main

import (
	"GoGameServer/core/consts/Service"
	. "GoGameServer/core/libs"
	"GoGameServer/core/service"
	"GoGameServer/servives/api/controllers"
)

func main() {
	//初始化Service
	newService := service.NewService(Service.Api)
	newService.StartRedis()
	newService.StartMongo()
	newService.StartHttpServer()
	newService.RegisterHttpRouter("/", &controllers.DefaultController{})

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initModule() {
}
