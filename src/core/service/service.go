package service

import (
	"core"
	"core/argv"
	"core/config"
	. "core/libs"
	"core/libs/consul"
	"core/libs/dict"
	"core/libs/grpc/ipc"
	"core/libs/logger"
	"core/libs/redis"
	"core/libs/stack"
	"core/libs/websocket"
	"core/message"
	"runtime"
)

type Service struct {
	env  string
	name string
	id   int

	port string

	ipcClients map[string]*ipc.Client
}

func NewService(name string) *Service {
	service := &Service{
		name: name,
	}
	service.init()

	core.Service = service
	return service
}

func (this *Service) init() {
	//错误捕获
	recoverErr()

	//初始化: 使用CPU数设置
	initMaxProcs()

	//初始化: 命令行参数
	initArgv(this)

	//初始化: 配置文件
	initConfig()

	//初始化: log
	initLog(this)

	//系统环境输出
	printEnv(this)
}

func initMaxProcs() {
	//runtime.GOMAXPROCS(1)
}

func initArgv(service *Service) {
	err := argv.Init()
	CheckError(err)

	service.env = argv.Values.Env
	service.id = argv.Values.ServiceId
}

func initConfig() {
	config.Init()
}

func initLog(service *Service) {
	logConfig := config.GetLog()

	logOpenDebug := dict.GetBool(logConfig, "debug")
	logOutput := dict.GetString(logConfig, "output")
	logFileName := service.name + "-" + NumToString(service.id)

	logger.SetLogFile(logFileName, logOutput)
	logger.SetLogDebug(logOpenDebug)
}

func initRedis() {
	redis.InitRedis(config.GetRedisList())
}

func printEnv(service *Service) {
	INFO("使用CPU数量:" + NumToString(runtime.GOMAXPROCS(-1)))
	INFO("初始GoroutineNum:" + NumToString(runtime.NumGoroutine()))
	INFO("服务平台:" + service.env)
	INFO("服务名称:" + service.name)
	INFO("服务ID:" + NumToString(service.id))
}

func recoverErr() {
	stack.PrintPanicStackError()
}

func (this *Service) registerService(servicePort string) {
	err := consul.InitServer(this.name, this.id, servicePort)
	CheckError(err)

	INFO("join consul service...." + servicePort)

	this.port = servicePort
}

/*********************************====以下为公开函数====*******************************/

func (this *Service) StartRedis() {
	initRedis()
}

func (this *Service) StartWebSocket() {
	//WebSocket配置
	serviceConfig := config.GetConnectorService(this.id)
	port := dict.GetString(serviceConfig, "clientPort")
	useSSL := dict.GetBool(serviceConfig, "useSSL")

	//创建WebSocket Server
	server := websocket.NewServer(port, this.id)
	if useSSL {
		tslCrt := config.GetConnectorServiceTslCrt()
		tslKey := config.GetConnectorServiceTslKey()
		server.SetTLS(tslCrt, tslKey)
	}
	server.SetSessionMsgHandle(message.FontReceive)
	server.SetSessionCloseHandle(message.SendMsgToBack_UserOffline)
	server.Start()
	server.StartPing()

	//服务注册
	this.registerService(port)
}

func (this *Service) StartIpcClient(serviceNames []string) {
	this.ipcClients = make(map[string]*ipc.Client)

	//初始化consul客户端
	consulClient, err := consul.InitClient()
	CheckError(err)

	//初始化Ipc客户端
	for _, serviceName := range serviceNames {
		this.ipcClients[serviceName] = ipc.InitClient(consulClient, serviceName, message.BackReceive)
		INFO("ipc client start....", serviceName)
	}
}

func (this *Service) Env() string {
	return this.env
}

func (this *Service) ID() int {
	return this.id
}

func (this *Service) Name() string {
	return this.name
}

func (this *Service) Port() string {
	return this.port
}

func (this *Service) GetIpcClient(serviceName string) *ipc.Client {
	client, _ := this.ipcClients[serviceName]
	return client
}
