package main

import (
	"core"
	"core/consts/service"
	. "core/libs"
	"core/libs/sessions"
	"core/service"
	_ "net/http/pprof"
	"servives/connector/module"
	module2 "servives/login/module"
)

func main() {
	//初始化Service
	newService := service.NewService(Service.Connector)
	newService.StartRedis()
	newService.StartWebSocket()
	newService.SetSessionCreateHandle(sessionCreate)
	newService.StartIpcClient([]string{Service.Game, Service.Login})
	newService.StartRpcClient([]string{Service.Game, Service.Login})

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
	loginService := core.Service.GetRpcClient("login")

	method := "ClientOffline"
	args := &module2.RpcClientOfflineReq{
		ServiceName: GetLocalIp() + "_" + core.Service.Name() + "_" + NumToString(core.Service.ID()),
		SessionId:   session.ID(),
	}
	reply := new(module2.RpcClientOfflineRes)
	loginService.Call(method, args, reply, "")
}
