package service

import (
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/yicaoyimuys/GoGameServer/core"
	"github.com/yicaoyimuys/GoGameServer/core/config"
	"github.com/yicaoyimuys/GoGameServer/core/consts/ServiceType"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/common"
	"github.com/yicaoyimuys/GoGameServer/core/libs/consul"
	"github.com/yicaoyimuys/GoGameServer/core/libs/grpc/ipc"
	"github.com/yicaoyimuys/GoGameServer/core/libs/logger"
	"github.com/yicaoyimuys/GoGameServer/core/libs/mongo"
	"github.com/yicaoyimuys/GoGameServer/core/libs/mysql"
	"github.com/yicaoyimuys/GoGameServer/core/libs/redis"
	"github.com/yicaoyimuys/GoGameServer/core/libs/rpc"
	"github.com/yicaoyimuys/GoGameServer/core/libs/sessions"
	"github.com/yicaoyimuys/GoGameServer/core/libs/socket"
	"github.com/yicaoyimuys/GoGameServer/core/libs/stack"
	"github.com/yicaoyimuys/GoGameServer/core/libs/system"
	"github.com/yicaoyimuys/GoGameServer/core/libs/timer"
	"github.com/yicaoyimuys/GoGameServer/core/libs/websocket"
	"github.com/yicaoyimuys/GoGameServer/core/messages"

	"github.com/astaxie/beego"
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

/*********************************====以下为公开函数====*******************************/

func (this *Service) StartRedis() {
	this.redisClients = make(map[string]*redis.Client)

	redisConfigs := config.GetRedisConfig()
	for aliasName, redisConfig := range redisConfigs {
		client, err := redis.NewClient(redisConfig)
		CheckError(err)

		if client != nil {
			this.redisClients[aliasName] = client
			INFO("redis_" + aliasName + "连接成功...")
		}
	}
}

func (this *Service) StartMysql() {
	this.mysqlClients = make(map[string]*mysql.Client)

	mysqlConfigs := config.GetMysqlConfig()
	index := 0
	for key, mysqlConfig := range mysqlConfigs {
		dbAliasName := key
		if index == 0 {
			dbAliasName = "default"
		}
		index++

		client, err := mysql.NewClient(dbAliasName, mysqlConfig)
		CheckError(err)

		if client != nil {
			this.mysqlClients[key] = client
			INFO("mysql_" + key + "连接成功...")
		}
	}
}

func (this *Service) StartMongo() {
	this.mongoClients = make(map[string]*mongo.Client)

	mongoConfigs := config.GetMongoConfig()
	for aliasName, mongoConfig := range mongoConfigs {
		client, err := mongo.NewClient(mongoConfig)
		CheckError(err)

		if client != nil {
			this.mongoClients[aliasName] = client
			INFO("mongo_" + aliasName + "连接成功")
		} else {
			ERR("mongo_" + aliasName + "连接失败")
		}
	}
}

func (this *Service) StartHttpServer() {
	//Api服务配置
	serviceConfig := config.GetService("api")
	serviceNodeConfig := serviceConfig.ServiceNodes[this.id]
	port := serviceNodeConfig.ClientPort
	useSSL := serviceNodeConfig.UseSSL

	//Http服务配置
	if useSSL {
		tslCrt := serviceConfig.TslCrt
		tslKey := serviceConfig.TslKey

		beego.BConfig.Listen.EnableHTTPS = true
		beego.BConfig.Listen.HTTPSCertFile = tslCrt
		beego.BConfig.Listen.HTTPSKeyFile = tslKey
		beego.BConfig.Listen.HTTPSPort = cast.ToInt(port)
	} else {
		beego.BConfig.Listen.HTTPPort = cast.ToInt(port)
	}
	beego.BConfig.RunMode = beego.PROD

	//启动http服务
	go beego.Run()

	//服务注册
	this.registerService(ServiceType.HTTP, port)
}

func (this *Service) RegisterHttpRouter(rootPath string, controller beego.ControllerInterface) {
	beego.Router(rootPath, controller)
}

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
	server.SetSessionReceiveMsgHandle(handle)
	server.Start()
	server.StartPing()

	//服务注册
	this.registerService(ServiceType.WEBSOCKET, port)

	//service中保存websocketServer
	this.websocketServer = server
}

func (this *Service) StartSocket(handle sessions.FrontSessionReceiveMsgHandle) {
	//Socket配置
	serviceConfig := config.GetService("connector")
	serviceNodeConfig := serviceConfig.ServiceNodes[this.id]
	port := serviceNodeConfig.ClientPort

	//创建Socket Server
	server := socket.NewServer(port, this.id)
	server.SetSessionReceiveMsgHandle(handle)
	server.Start()
	server.StartPing()

	//服务注册
	this.registerService(ServiceType.SOCKET, port)

	//service中保存socketServer
	this.socketServer = server
}

func (this *Service) SetSessionCreateHandle(handle sessions.FrontSessionCreateHandle) {
	if this.websocketServer != nil {
		this.websocketServer.SetSessionCreateHandle(handle)
	}
	if this.socketServer != nil {
		this.socketServer.SetSessionCreateHandle(handle)
	}
}

func (this *Service) StartIpcClient(serviceNames []string) {
	this.ipcClients = make(map[string]*ipc.Client)

	//初始化consul客户端
	consulClient, err := consul.NewClient()
	CheckError(err)

	//初始化Ipc客户端
	for _, serviceName := range serviceNames {
		serviceName = packageServiceName(ServiceType.IPC, serviceName)
		this.ipcClients[serviceName] = ipc.NewClient(consulClient, serviceName, messages.IpcClientReceive)
		INFO("ipc client start...", serviceName)
	}
}

func (this *Service) StartIpcServer() {
	//开启ipcServer
	ipcServer, port, err := ipc.InitServer(messages.IpcServerReceive)
	CheckError(err)
	INFO("ipc server start...", port)

	//service中记录ipcServer
	this.ipcServer = ipcServer

	//服务注册
	this.registerService(ServiceType.IPC, port)

	//Log
	timer.DoTimer(20*1000, func() {
		INFO("当前BackSession数量:", sessions.BackSessionLen())
	})
}

func (this *Service) StartRpcClient(serviceNames []string) {
	this.rpcClients = make(map[string]*rpc.Client)

	//初始化consul客户端
	consulClient, err := consul.NewClient()
	CheckError(err)

	//初始化Rpc客户端
	for _, serviceName := range serviceNames {
		serviceName = packageServiceName(ServiceType.RPC, serviceName)
		this.rpcClients[serviceName] = rpc.NewClient(consulClient, serviceName)
		INFO("rpc client start...", serviceName)
	}
}

func (this *Service) StartRpcServer() {
	//开启rpcServer
	port, err := rpc.InitServer()
	CheckError(err)
	INFO("rpc server start...." + port)

	//服务注册
	this.registerService(ServiceType.RPC, port)
}

func (this *Service) RegisterRpcModule(rpcName string, rpcModule interface{}) {
	//rpc模块注册
	err := rpc.RegisterModule(rpcName, rpcModule)
	CheckError(err)
}

func (this *Service) StartPProf(port int) {
	port = port + this.id
	go func() {
		defer stack.TryError()
		http.ListenAndServe(":"+cast.ToString(port), nil)
	}()
	INFO("debug start...", port)
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

func (this *Service) GetIpcClient(serviceName string) *ipc.Client {
	serviceName = packageServiceName(ServiceType.IPC, serviceName)
	client, _ := this.ipcClients[serviceName]
	return client
}

func (this *Service) GetRpcClient(serviceName string) *rpc.Client {
	serviceName = packageServiceName(ServiceType.RPC, serviceName)
	client, _ := this.rpcClients[serviceName]
	return client
}

func (this *Service) GetRedisClient(redisAliasName string) *redis.Client {
	client, _ := this.redisClients[redisAliasName]
	return client
}

func (this *Service) GetMysqlClient(dbAliasName string) *mysql.Client {
	client, _ := this.mysqlClients[dbAliasName]
	return client
}

func (this *Service) GetMongoClient(dbAliasName string) *mongo.Client {
	client, _ := this.mongoClients[dbAliasName]
	return client
}

func (this *Service) GetIpcServer() *ipc.Server {
	return this.ipcServer
}
