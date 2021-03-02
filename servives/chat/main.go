package main

import (
	"GoGameServer/core/consts/Service"
	. "GoGameServer/core/libs"
	"GoGameServer/core/messages"
	"GoGameServer/core/protos/gameProto"
	"GoGameServer/core/service"
	"GoGameServer/servives/chat/module"
	"GoGameServer/servives/public/rpcModules"
)

func main() {
	//初始化Service
	newService := service.NewService(Service.Chat)
	newService.StartIpcServer()
	newService.StartRpcServer()
	newService.StartRpcClient([]string{Service.Log})
	newService.StartRedis()
	newService.RegisterRpcModule("Client", &rpcModules.Client{})

	//消息初始化
	initMessage()

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initMessage() {
	messages.RegisterIpcServerHandle(gameProto.ID_user_joinChat_c2s, module.JoinChat)
	messages.RegisterIpcServerHandle(gameProto.ID_user_chat_c2s, module.Chat)
}

func initModule() {

}
