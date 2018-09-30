package main

import (
	"core/consts/service"
	. "core/libs"
	"core/service"
	"servives/api/controllers"
)

func main() {
	//初始化Service
	newService := service.NewService(Service.Api)
	newService.StartRedis()
	newService.StartMysql()
	newService.StartHttpServer()
	newService.RegisterHttpRouter("/", &controllers.DefaultController{})

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initModule() {
}
