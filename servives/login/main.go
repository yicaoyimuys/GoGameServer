package main

import (
	"GoGameServer/core/consts/Service"
	. "GoGameServer/core/libs"
	"GoGameServer/core/messages"
	"GoGameServer/core/protos/gameProto"
	"GoGameServer/core/service"
	"GoGameServer/servives/login/module"
	"GoGameServer/servives/public/rpcModules"
)

func main() {
	//初始化Service
	newService := service.NewService(Service.Login)
	newService.StartIpcServer()
	newService.StartRpcServer()
	newService.StartRpcClient([]string{Service.Log})
	newService.StartRedis()
	newService.StartMysql()
	newService.RegisterRpcModule("Client", &rpcModules.Client{})

	//消息初始化
	initMessage()

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initMessage() {
	messages.RegisterIpcServerHandle(gameProto.ID_user_login_c2s, module.Login)
}

func initModule() {

}
