package service

import (
	"runtime"

	beegoLogs "github.com/astaxie/beego/logs"
	beegoOrm "github.com/astaxie/beego/orm"
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
	"go.uber.org/zap"

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

	logger.Init(
		logger.WithDebug(logConfig.Debug),
		logger.WithBoth(logConfig.Both),
		logger.WithFile(logConfig.File),
		logger.WithName(service.name+"-"+cast.ToString(service.id)),
	)

	// 设置beego logs
	if logConfig.Debug {
		beegoLogs.SetLevel(beegoLogs.LevelDebug)
	} else {
		beegoLogs.SetLevel(beegoLogs.LevelInfo)
	}

	// 设置beego orm
	beegoOrm.Debug = logConfig.Debug
}

func printEnv(service *Service) {
	INFO("CPU数量", zap.Int("CpuNum", runtime.GOMAXPROCS(-1)))
	INFO("协程数量", zap.Int("GoroutineNum", runtime.NumGoroutine()))
	INFO("Go版本", zap.String("GoVersion", runtime.Version()))
	INFO("启动路径", zap.String("Root", system.Root))
	INFO("服务器环境", zap.String("ServiceEnv", service.env))
	INFO("服务器名称", zap.String("ServiceName", service.name))
	INFO("服务器ID", zap.Int("ServiceId", service.id))
	INFO("服务器IP", zap.String("ServiceIp", common.GetLocalIp()))

	timer.DoTimer(20*1000, func() {
		INFO("协程数量", zap.Int("GoroutineNum", runtime.NumGoroutine()))
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
		ERR("该类型的Service已经在本进程内启用", zap.String("ServiceType", serviceType))
		return
	}

	//注册到Consul
	serviceName := packageServiceName(serviceType, this.name)
	err := consul.NewServive(this.ip, serviceName, this.id, servicePort)
	CheckError(err)

	INFO("Join Consul Service", zap.String("ServiceName", serviceName), zap.String("ServicePort", servicePort))

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
