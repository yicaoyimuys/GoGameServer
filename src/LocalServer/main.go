package main

import (
	"os"
)

import (
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"global"
	"module"
	"proxys/redisProxy"
	. "tools"
	"tools/cfg"
	"tools/db"
)

//各个模块
import (
	_ "module/cache"
	_ "module/config"
	_ "module/user"
)

var (
	local_ip   string
	local_port string
)

func main() {
	defer func() {
		if x := recover(); x != nil {
			ERR("caught panic in main()", x)
		}
	}()

	// 获取端口号
	getPort()

	//启动
	global.Startup(global.ServerName, "local_log", nil)

	//开启LocalServer监听
	startLocalServer()

	//保持进程
	global.Run()
}

func getPort() {
	//端口号
	local_ip = cfg.GetValue("local_ip")
	local_port = cfg.GetValue("local_port")
	global.ServerName = "LocalServer[" + local_port + "]"
}

func startLocalServer() {
	//连接DB
	db.Init()

	//连接Redis
	redisProxyErr := redisProxy.InitClient(cfg.GetValue("redis_ip"), cfg.GetValue("redis_port"))
	checkError(redisProxyErr)

	listener, err := link.Serve("tcp", "0.0.0.0:"+local_port, packet.New(
		binary.SplitByUint32BE, 1024, 1024, 1024,
	))
	checkError(err)

	listener.Serve(func(session *link.Session) {
		session.AddCloseCallback(session, func() {
			session.Close()
		})
		global.AddSession(session)

		var msg packet.RAW
		for {
			if err := session.Receive(&msg); err != nil {
				break
			}
			module.ReceiveMessage(session, msg)
		}
	})
}

func checkError(err error) {
	if err != nil {
		ERR("Fatal error: %v", err)
		os.Exit(-1)
	}
}
