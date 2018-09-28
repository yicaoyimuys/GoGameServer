package main

import (
	"core/consts/service"
	. "core/libs"
	"core/messages"
	"core/protos/gameProto"
	"core/service"
	"servives/chat/module"
	"servives/public/rpcModules"
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
