package service

import (
	"runtime"

	"github.com/yicaoyimuys/GoGameServer/core"
	"github.com/yicaoyimuys/GoGameServer/core/config"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/common"
	"github.com/yicaoyimuys/GoGameServer/core/libs/consul"
	"github.com/yicaoyimuys/GoGameServer/core/libs/grpc/ipc"
	"github.com/yicaoyimuys/GoGameServer/core/libs/logger"
	"github.com/yicaoyimuys/GoGameServer/core/libs/mongo"
	"github.com/yicaoyimuys/GoGameServer/core/libs/mysql"
	"github.com/yicaoyimuys/GoGameServer/core/libs/redis"
	"github.com/yicaoyimuys/GoGameServer/core/libs/rpc"
	"github.com/yicaoyimuys/GoGameServer/core/libs/socket"
	"github.com/yicaoyimuys/GoGameServer/core/libs/stack"
	"github.com/yicaoyimuys/GoGameServer/core/libs/system"
	"github.com/yicaoyimuys/GoGameServer/core/libs/timer"
	"github.com/yicaoyimuys/GoGameServer/core/libs/websocket"

	"github.com/spf13/cast"
)

type Service struct {
	env  string
	name string
	id   int

	ip    string
	ports map[string]string

	ipcServer *ipc.Server

	ipcClients   map[string]*ipc.Client
	rpcClients   map[string]*rpc.Client
	redisClients map[string]*redis.Client
	mysqlClients map[string]*mysql.Client
	mongoClients map[string]*mongo.Client

	websocketServer *websocket.Server
	socketServer    *socket.Server
}

func NewService(name string) *Service {
	service := &Service{
		name:  name,
		ip:    common.GetLocalIp(),
		ports: make(map[string]string),
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
	initConfig(this)

	//初始化: log
	initLog(this)

	//系统环境输出
	printEnv(this)
}

func initMaxProcs() {
	//runtime.GOMAXPROCS(1)
}

func initArgv(service *Service) {
	service.env = system.Args.Env
	service.id = system.Args.ServiceId
}

func initConfig(service *Service) {
	config.Init(service.env)
}

func initLog(service *Service) {
	logConfig := config.GetLogConfig()

	logOpenDebug := logConfig.Debug
	logOutput := logConfig.Output
	logFileName := service.name + "-" + cast.ToString(service.id)

	logger.SetLogFile(logFileName, logOutput)
	logger.SetLogDebug(logOpenDebug)
}

func printEnv(service *Service) {
	INFO("使用CPU数量:", runtime.GOMAXPROCS(-1))
	INFO("初始GoroutineNum:", runtime.NumGoroutine())
	INFO("服务平台:", service.env)
	INFO("服务名称:", service.name)
	INFO("服务ID:", service.id)
	INFO("启动参数:", system.Args)

	timer.DoTimer(20*1000, func() {
		INFO("当前GoroutineNum:", runtime.NumGoroutine())
	})
}

func recoverErr() {
	stack.TryError()
}

func packageServiceName(serviceType string, serviceName string) string {
	return "<" + serviceType + ">" + serviceName
}

func (this *Service) registerService(serviceType string, servicePort string) {
	if _, exists := this.ports[serviceType]; exists {
		ERR("该类型的Service已经在本进程内启用", serviceType)
		return
	}

	//注册到Consul
	serviceName := packageServiceName(serviceType, this.name)
	err := consul.NewServive(this.ip, serviceName, this.id, servicePort)
	CheckError(err)

	INFO("join consul service...", serviceName, servicePort)

	//记录该进程启用的端口号
	this.ports[serviceType] = servicePort
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

func (this *Service) Ip() string {
	return this.ip
}

func (this *Service) Port(serviceType string) string {
	return this.ports[serviceType]
}

func (this *Service) Identify() string {
	return this.ip + "_" + this.name + "_" + cast.ToString(this.id)
}
