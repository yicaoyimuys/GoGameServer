package main

import (
	"github.com/yicaoyimuys/GoGameServer/core/consts"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/service"
	"github.com/yicaoyimuys/GoGameServer/servives/connector/messages"
	"github.com/yicaoyimuys/GoGameServer/servives/connector/module"
)

func main() {
	//初始化Service
	newService := service.NewService(consts.Service_Connector)
	newService.StartRedis()
	// newService.StartWebSocket(messages.FontReceive)
	newService.StartSocket(messages.FontReceive)
	newService.StartIpcClient([]string{consts.Service_Game, consts.Service_Login, consts.Service_Chat})
	newService.StartRpcClient([]string{consts.Service_Game, consts.Service_Login, consts.Service_Chat})
	newService.StartPProf(6000)

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initModule() {
	module.StartServerTimer()
}
