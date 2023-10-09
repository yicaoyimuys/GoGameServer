package service

import (
	"github.com/yicaoyimuys/GoGameServer/core/consts"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/consul"
	"github.com/yicaoyimuys/GoGameServer/core/libs/grpc/ipc"
	"github.com/yicaoyimuys/GoGameServer/core/libs/sessions"
	"github.com/yicaoyimuys/GoGameServer/core/libs/timer"
	"github.com/yicaoyimuys/GoGameServer/core/messages"
	"go.uber.org/zap"
)

func (this *Service) StartIpcClient(serviceNames []string) {
	this.ipcClients = make(map[string]*ipc.Client)

	//初始化consul客户端
	consulClient, err := consul.NewClient()
	CheckError(err)

	//初始化Ipc客户端
	for _, serviceName := range serviceNames {
		serviceName = packageServiceName(consts.ServiceType_Ipc, serviceName)
		this.ipcClients[serviceName] = ipc.NewClient(consulClient, serviceName, messages.IpcClientReceive)
		INFO("Ipc Client Start", zap.String("ServiceName", serviceName))
	}
}

func (this *Service) StartIpcServer() {
	//开启ipcServer
	ipcServer, port, err := ipc.InitServer(messages.IpcServerReceive)
	CheckError(err)
	INFO("Ipc Server Start", zap.String("Port", port))

	//service中记录ipcServer
	this.ipcServer = ipcServer

	//服务注册
	this.registerService(consts.ServiceType_Ipc, port)

	//Log
	timer.DoTimer(20*1000, func() {
		INFO("当前BackSession数量", zap.Int("BackSessionLen", sessions.BackSessionLen()))
	})
}

func (this *Service) GetIpcClient(serviceName string) *ipc.Client {
	serviceName = packageServiceName(consts.ServiceType_Ipc, serviceName)
	client, _ := this.ipcClients[serviceName]
	return client
}

func (this *Service) GetIpcServer() *ipc.Server {
	return this.ipcServer
}
