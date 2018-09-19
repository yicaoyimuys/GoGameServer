package main

import (
	"os"
)

import (
	"core/argv"
	"core/config"
	. "core/libs"
	"core/libs/consul"
	"core/libs/grpc/ipc"
	"core/libs/redis"
	"core/libs/timer"
	"core/libs/websocket"
	"game/module"
	"global"
	"message"
)

import (
	"core"
	"core/libs/dict"
	_ "net/http/pprof"
)

func main() {
	//初始化
	serviceName := "connector"
	core.NewService(serviceName)

	//Service配置
	global.ServiceName = argv.Values.ServiceName + "-" + NumToString(argv.Values.ServiceId)

	//Server启动端口设置
	serviceConfig := config.GetConnectorService(argv.Values.ServiceId)
	global.ServerPort = NumToString(serviceConfig["clientPort"])

	//Guid初始化
	global.InitGuid(uint16(argv.Values.ServiceId))

	//Redis配置
	redis.InitRedis(config.GetRedisList())

	//Kv初始化
	err := consul.InitKV(true)
	checkError(err)

	//开启WebSocket
	startWs()
	//开启Ipc
	go startIpcClient()
	//服务注册
	go registService()

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initModule() {
	module.StartServerTimer()
}

func startWs() {
	//WebSocket配置
	serviceConfig := config.GetConnectorService(argv.Values.ServiceId)
	port := dict.GetString(serviceConfig, "clientPort")
	useSSL := dict.GetBool(serviceConfig, "useSSL")

	//创建WebSocket Server
	server := websocket.NewServer(port)
	if useSSL {
		tslCrt := config.GetConnectorServiceTslCrt()
		tslKey := config.GetConnectorServiceTslKey()
		server.SetTLS(tslCrt, tslKey)
	}
	server.SetSessionMsgHandle(message.FontReceive)
	server.SetSessionCloseHandle(message.SendMsgToBack_UserOffline)
	server.Start()
	server.StartPing()
}

func startIpcClient() {
	//初始化consul客户端
	consulClient, err := consul.InitClient()
	checkError(err)

	//初始化Ipc客户端(Game)
	serviceName := global.Services.Game
	global.IpcClients.Game = ipc.InitClient(consulClient, serviceName, message.BackReceive)
	INFO("ipc client start....", serviceName)

	//初始化Ipc客户端(Matching)
	serviceName = global.Services.Matching
	global.IpcClients.Matching = ipc.InitClient(consulClient, serviceName, message.BackReceive)
	INFO("ipc client start....", serviceName)
}

func registService() {
	serverName := global.ServiceName
	serverPort := global.ServerPort

	INFO("join consul service...." + serverPort)

	err := consul.InitServer(serverName, serverPort)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		ERR("Fatal error: %v", err)
		timer.SetTimeOut(1000, func() {
			os.Exit(-1)
		})
	}
}
