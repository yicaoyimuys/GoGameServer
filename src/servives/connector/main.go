package main

import (
	"core"
	"core/consts/service"
	. "core/libs"
	"core/libs/sessions"
	"core/service"
	moduleChat "servives/chat/module"
	"servives/connector/module"
	moduleLogin "servives/login/module"
)

func main() {
	//初始化Service
	newService := service.NewService(Service.Connector)
	newService.StartRedis()
	newService.StartWebSocket()
	newService.SetSessionCreateHandle(sessionCreate)
	newService.StartIpcClient([]string{Service.Game, Service.Login, Service.Chat})
	newService.StartRpcClient([]string{Service.Game, Service.Login, Service.Chat})

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
	{
		//通知登录服务器
		loginService := core.Service.GetRpcClient(Service.Login)

		method := "ClientOffline"
		args := &moduleLogin.RpcClientOfflineReq{
			ServiceIdentify: core.Service.Identify(),
			UserSessionId:   session.ID(),
		}
		reply := new(moduleLogin.RpcClientOfflineRes)
		loginService.Call(method, args, reply, "")
	}

	{
		//通知聊天服务器
		chatService := core.Service.GetRpcClient(Service.Chat)

		method := "ClientOffline"
		args := moduleChat.RpcClientOfflineReq{
			ServiceIdentify: core.Service.Identify(),
			UserSessionId:   session.ID(),
		}
		reply := new(moduleChat.RpcClientOfflineRes)
		chatService.Call(method, args, reply, "")
	}
}
