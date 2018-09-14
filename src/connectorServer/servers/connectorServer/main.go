package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"runtime"
)

import (
	"connectorServer/config"
	"connectorServer/game/module"
	"connectorServer/global"
	"connectorServer/message"
	"connectorServer/sessions"
	. "connectorServer/tools"
	"connectorServer/tools/consul"
	"connectorServer/tools/grpc/ipc"
	"connectorServer/tools/redis"
	"connectorServer/tools/stack"
	"connectorServer/tools/timer"
)

import (
	_ "net/http/pprof"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	defer func() {
		if x := recover(); x != nil {
			ERR("caught panic in main()", x)
		}
	}()

	//CPU数设置
	//runtime.GOMAXPROCS(1)

	//初始化命令行参数
	flag.StringVar(&global.Env, "e", "development", "env")
	flag.IntVar(&global.ServerId, "g", 1, "serverId")
	flag.IntVar(&global.GameServerComputer, "s", 0, "serverNumber")
	flag.Parse()

	//Server配置
	serverConfig := config.GetConnectorServer(global.ServerId)
	global.ServerName = serverConfig["id"].(string)

	//Server启动端口设置
	if global.GameServerComputer > 0 {
		global.ServerPort = NumToString(10000 + 100*global.GameServerComputer + global.ServerId)
	} else {
		global.ServerPort = NumToString(serverConfig["clientPort"])
	}

	//Log配置
	logConfig := config.GetLog()
	SetLogDebug(logConfig["debug"].(bool))
	SetLogFile(global.ServerName, logConfig["output"].(string))

	//Guid初始化
	global.InitGuid(uint16(global.ServerId))

	//Redis配置
	redis.InitRedis(config.GetRedisList())

	//系统环境
	INFO("使用CPU数量:" + NumToString(runtime.GOMAXPROCS(-1)))
	INFO("初始GoroutineNum:" + NumToString(runtime.NumGoroutine()))
	INFO("服务器平台:" + global.Env)

	//Kv初始化
	err := consul.InitKV(true)
	checkError(err)

	//开启WebSocket
	go startWs(serverConfig["useSSL"].(bool))
	//开启Ipc
	go startIpcClient()
	//服务注册
	go registService()

	//开启Ping检测
	overTime := If(global.Env == "facebook", 60, 15).(int)
	sessions.FrontSessionOpenPing(int64(overTime))
	INFO("Session超时时间设置", overTime)

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initModule() {
	module.StartServerTimer()
}

func addFontSession(session *sessions.FrontSession) {
	sessions.AddFrontSession(session)
	session.SetMsgHandle(message.FontReceive)
	session.AddCloseCallback(nil, "main.FrontSessionOffline", func() {
		message.SendMsgToBack_UserOffline(session)
		//DEBUG("session count: ", sessions.FrontSessionLen())
	})
	//DEBUG("session count: ", sessions.FrontSessionLen())

	defer session.Close()
	for {
		msg, err := session.Receive()
		if err != nil || msg == nil {
			break
		}
	}
}

func startWs(useSSL bool) {
	gameServerPort := global.ServerPort

	INFO("front start websocket...." + gameServerPort)

	http.HandleFunc("/", wsHandler)
	var err error
	if useSSL {
		tslCrt := config.GetGameServerTslCrt()
		tslKey := config.GetGameServerTslKey()
		err = http.ListenAndServeTLS("0.0.0.0:"+gameServerPort, tslCrt, tslKey, nil)
	} else {
		err = http.ListenAndServe("0.0.0.0:"+gameServerPort, nil)
	}
	checkError(err)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ERR("wsHandler: ", err)
		return
	}

	defer stack.PrintPanicStackError()

	//Session创建
	session := sessions.NewFontSession(sessions.NewFrontCodec(conn))
	addFontSession(session)
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
	serverName := global.ServerName
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
