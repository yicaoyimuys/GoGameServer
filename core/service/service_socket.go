package service

import (
	"github.com/yicaoyimuys/GoGameServer/core/config"
	"github.com/yicaoyimuys/GoGameServer/core/consts"
	"github.com/yicaoyimuys/GoGameServer/core/libs/sessions"
	"github.com/yicaoyimuys/GoGameServer/core/libs/socket"
)

func (this *Service) StartSocket(handle sessions.FrontSessionReceiveMsgHandle) {
	//Socket配置
	serviceConfig := config.GetService("connector")
	serviceNodeConfig := serviceConfig.ServiceNodes[this.id]
	port := serviceNodeConfig.ClientPort

	//创建Socket Server
	server := socket.NewServer(port, this.id)
	server.SetSessionCreateHandle(this.frontSessionCreateHandle)
	server.SetSessionReceiveMsgHandle(handle)
	server.Start()
	server.StartPing()

	//服务注册
	this.registerService(consts.ServiceType_Socket, port)

	//service中保存socketServer
	this.socketServer = server
}
