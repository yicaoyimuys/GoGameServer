package service

import (
	"core"
	"core/config"
	. "core/libs"
	"core/libs/argv"
	"core/libs/consul"
	"core/libs/dict"
	"core/libs/grpc/ipc"
	"core/libs/logger"
	"core/libs/mysql"
	"core/libs/redis"
	"core/libs/rpc"
	"core/libs/stack"
	"core/libs/websocket"
	"core/message"
	"net/http"
	"runtime"
)

const (
	WS  = "ws"
	RPC = "rpc"
	IPC = "ipc"
)

type Service struct {
	env  string
	name string
	id   int

	port string

	ipcClients   map[string]*ipc.Client
	rpcClients   map[string]*rpc.Client
	redisClients map[string]*redis.Client
	dbClients    map[string]*mysql.Client
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
	err := argv.Init()
	CheckError(err)

	service.env = argv.Values.Env
	service.id = argv.Values.ServiceId
}

func initConfig(service *Service) {
	config.Init(service.env)
}

func initLog(service *Service) {
	logConfig := config.GetLogConfig()

	logOpenDebug := dict.GetBool(logConfig, "debug")
	logOutput := dict.GetString(logConfig, "output")
	logFileName := service.name + "-" + NumToString(service.id)

	logger.SetLogFile(logFileName, logOutput)
	logger.SetLogDebug(logOpenDebug)
}

func printEnv(service *Service) {
	INFO("使用CPU数量:", runtime.GOMAXPROCS(-1))
	INFO("初始GoroutineNum:", runtime.NumGoroutine())
	INFO("服务平台:", service.env)
	INFO("服务名称:", service.name)
	INFO("服务ID:", service.id)
}

func recoverErr() {
	stack.PrintPanicStackError()
}

func packageServiceName(serviceType string, serviceName string) string {
	return "<" + serviceType + ">" + serviceName
}

func (this *Service) registerService(serviceType string, servicePort string) {
	serviceName := packageServiceName(serviceType, this.name)
	err := consul.NewServive(serviceName, this.id, servicePort)
	CheckError(err)

	INFO("join consul service...", serviceName, servicePort)

	this.port = servicePort
}

/*********************************====以下为公开函数====*******************************/

func (this *Service) StartRedis() {
	this.redisClients = make(map[string]*redis.Client)

	redisConfigs := config.GetRedisConfig()
	for aliasName, redisConfig := range redisConfigs {
		client, err := redis.NewClient(redisConfig.(map[string]interface{}))
		CheckError(err)

		this.redisClients[aliasName] = client
		INFO("redis_" + aliasName + "连接成功...")
	}
}

func (this *Service) StartMysql() {
	this.dbClients = make(map[string]*mysql.Client)

	mysqlConfigs := config.GetMysqlConfig()
	for aliasName, mysqlConfig := range mysqlConfigs {
		client, err := mysql.NewClient(aliasName, mysqlConfig.(map[string]interface{}))
		CheckError(err)

		this.dbClients[aliasName] = client
		INFO("mysql_" + aliasName + "连接成功...")
	}
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
	this.registerService(WS, port)
}

func (this *Service) StartIpcClient(serviceNames []string) {
	this.ipcClients = make(map[string]*ipc.Client)

	//初始化consul客户端
	consulClient, err := consul.NewClient()
	CheckError(err)

	//初始化Ipc客户端
	for _, serviceName := range serviceNames {
		serviceName = packageServiceName(IPC, serviceName)
		this.ipcClients[serviceName] = ipc.NewClient(consulClient, serviceName, message.IpcClientReceive)
		INFO("ipc client start...", serviceName)
	}
}

func (this *Service) StartIpcServer() {
	//开启ipcServer
	port, err := ipc.InitServer(message.IpcServerReceive)
	CheckError(err)
	INFO("ipc server start...", port)

	//服务注册
	this.registerService(IPC, port)
}

func (this *Service) StartRpcClient(serviceNames []string) {
	this.rpcClients = make(map[string]*rpc.Client)

	//初始化consul客户端
	consulClient, err := consul.NewClient()
	CheckError(err)

	//初始化Rpc客户端
	for _, serviceName := range serviceNames {
		serviceName = packageServiceName(RPC, serviceName)
		this.rpcClients[serviceName] = rpc.NewClient(consulClient, serviceName)
		INFO("rpc client start...", serviceName)
	}
}

func (this *Service) StartRpcServer(rcvr interface{}) {
	//rpc模块注册
	err := rpc.RegisterModule(this.name, rcvr)
	CheckError(err)

	//开启rpcServer
	port, err := rpc.InitServer()
	CheckError(err)
	INFO("rpc server start...." + port)

	//服务注册
	this.registerService(RPC, port)
}

func (this *Service) StartDebug() {
	port := 6060 + this.id
	go func() {
		defer stack.PrintPanicStackError()
		http.ListenAndServe(":"+NumToString(port), nil)
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

func (this *Service) Port() string {
	return this.port
}

func (this *Service) GetIpcClient(serviceName string) *ipc.Client {
	serviceName = packageServiceName(IPC, serviceName)
	client, _ := this.ipcClients[serviceName]
	return client
}

func (this *Service) GetRpcClient(serviceName string) *rpc.Client {
	serviceName = packageServiceName(RPC, serviceName)
	client, _ := this.rpcClients[serviceName]
	return client
}

func (this *Service) GetRedisClient(redisAliasName string) *redis.Client {
	client, _ := this.redisClients[redisAliasName]
	return client
}

func (this *Service) GetMysqlClient(dbAliasName string) *mysql.Client {
	client, _ := this.dbClients[dbAliasName]
	return client
}
