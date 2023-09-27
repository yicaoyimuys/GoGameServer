package main

import (
	"github.com/yicaoyimuys/GoGameServer/core"
	"github.com/yicaoyimuys/GoGameServer/core/consts/Service"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/sessions"
	"github.com/yicaoyimuys/GoGameServer/core/service"
	"github.com/yicaoyimuys/GoGameServer/servives/connector/messages"
	"github.com/yicaoyimuys/GoGameServer/servives/connector/module"
	"github.com/yicaoyimuys/GoGameServer/servives/public/rpcModules"
)

func main() {
	//初始化Service
	newService := service.NewService(Service.Connector)
	newService.StartRedis()
	// newService.StartWebSocket(messages.FontReceive)
	newService.StartSocket(messages.FontReceive)
	newService.SetSessionCreateHandle(sessionCreate)
	newService.StartIpcClient([]string{Service.Game, Service.Login, Service.Chat})
	newService.StartRpcClient([]string{Service.Game, Service.Login, Service.Chat})
	newService.StartPProf(6000)

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initModule() {
	module.StartServerTimer()
}

func sessionCreate(session *sessions.FrontSession) {
	session.AddCloseCallback(nil, "FrontSessionOffline", func() {
		sessionOffline(session)
	})
}

func sessionOffline(session *sessions.FrontSession) {
	method := "Client.Offline"
	args := &rpcModules.ClientOfflineReq{
		ServiceIdentify: core.Service.Identify(),
		UserSessionId:   session.ID(),
	}
	reply := &rpcModules.ClientOfflineRes{}

	//通知登录服务器
	go func() {
		loginService := core.Service.GetRpcClient(Service.Login)
		loginService.CallAll(method, args, reply)
	}()

	//通知聊天服务器
	go func() {
		chatService := core.Service.GetRpcClient(Service.Chat)
		chatService.CallAll(method, args, reply)
	}()

	//通知游戏服务器
	go func() {
		gameService := core.Service.GetRpcClient(Service.Game)
		gameService.CallAll(method, args, reply)
	}()
}
