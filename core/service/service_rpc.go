package service

import (
	"github.com/yicaoyimuys/GoGameServer/core/consts"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/consul"
	"github.com/yicaoyimuys/GoGameServer/core/libs/rpc"
	"go.uber.org/zap"
)

func (this *Service) StartRpcClient(serviceNames []string) {
	this.rpcClients = make(map[string]*rpc.Client)

	//初始化consul客户端
	consulClient, err := consul.NewClient()
	CheckError(err)

	//初始化Rpc客户端
	for _, serviceName := range serviceNames {
		serviceName = packageServiceName(consts.ServiceType_Rpc, serviceName)
		this.rpcClients[serviceName] = rpc.NewClient(consulClient, serviceName)
		INFO("Rpc Client Start", zap.String("ServiceName", serviceName))
	}
}

func (this *Service) StartRpcServer() {
	//开启rpcServer
	port, err := rpc.InitServer()
	CheckError(err)
	INFO("Rpc Server Start", zap.String("Port", port))

	//注册客户端下线回调
	err = rpc.RegisterModule("ClientOffline", &ClientOffline{})
	CheckError(err)

	//服务注册
	this.registerService(consts.ServiceType_Rpc, port)
}

func (this *Service) RegisterRpcModule(rpcName string, rpcModule interface{}) {
	//rpc模块注册
	err := rpc.RegisterModule(rpcName, rpcModule)
	CheckError(err)
}

func (this *Service) GetRpcClient(serviceName string) *rpc.Client {
	serviceName = packageServiceName(consts.ServiceType_Rpc, serviceName)
	client, _ := this.rpcClients[serviceName]
	return client
}
