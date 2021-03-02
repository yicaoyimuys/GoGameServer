package main

import (
	"GoGameServer/core"
	"GoGameServer/core/consts/Service"
	. "GoGameServer/core/libs"
	"GoGameServer/core/libs/sessions"
	"GoGameServer/core/service"
	"GoGameServer/servives/connector/module"
	"GoGameServer/servives/public/rpcModules"
)

func main() {
	//初始化Service
	newService := service.NewService(Service.Connector)
	newService.StartRedis()
	newService.StartWebSocket()
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
	reply := &rpcModules.ClientOfflineReq{}

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
