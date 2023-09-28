package service

import (
	"github.com/yicaoyimuys/GoGameServer/core/config"
	"github.com/yicaoyimuys/GoGameServer/core/consts"
	"github.com/yicaoyimuys/GoGameServer/core/libs/sessions"
	"github.com/yicaoyimuys/GoGameServer/core/libs/websocket"
)

func (this *Service) StartWebSocket(handle sessions.FrontSessionReceiveMsgHandle) {
	//WebSocket配置
	serviceConfig := config.GetService("connector")
	serviceNodeConfig := serviceConfig.ServiceNodes[this.id]
	port := serviceNodeConfig.ClientPort
	useSSL := serviceNodeConfig.UseSSL

	//创建WebSocket Server
	server := websocket.NewServer(port, this.id)
	if useSSL {
		tslCrt := serviceConfig.TslCrt
		tslKey := serviceConfig.TslKey
		server.SetTLS(tslCrt, tslKey)
	}
	server.SetSessionCreateHandle(this.frontSessionCreateHandle)
	server.SetSessionReceiveMsgHandle(handle)
	server.Start()
	server.StartPing()

	//服务注册
	this.registerService(consts.ServiceType_WebSocket, port)

	//service中保存websocketServer
	this.websocketServer = server
}
