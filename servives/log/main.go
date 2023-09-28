package main

import (
	"github.com/yicaoyimuys/GoGameServer/core/consts"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/service"
)

func main() {
	//初始化Service
	newService := service.NewService(consts.Service_Log)
	newService.StartRpcServer()

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initModule() {

}
