package main

import (
	"github.com/yicaoyimuys/GoGameServer/core/consts/Service"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/service"
	"github.com/yicaoyimuys/GoGameServer/servives/api/controllers"
)

func main() {
	//初始化Service
	newService := service.NewService(Service.Api)
	newService.StartRedis()
	newService.StartMongo()
	newService.StartHttpServer()
	newService.RegisterHttpRouter("/", &controllers.DefaultController{})
	newService.RegisterHttpRouter("/GetConnector", &controllers.ConnectorController{})

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initModule() {
}
