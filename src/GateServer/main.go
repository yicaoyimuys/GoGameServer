package main

import (
	"os"
)

import (
	"github.com/funny/link"
	"global"
	"proxys/transferProxy"
	. "tools"
	"tools/cfg"
	"tools/dispatch"
)

var (
	gateway_ip    string
	gateway_port  string
	transfer_ip   string
	transfer_port string
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
	global.Startup(global.ServerName, "gateway_log", nil)

	//开启TransferServer，由GateServer充当中转服务器
	err := transferProxy.InitServer(transfer_port)
	checkError(err)
	INFO("Starting TransferServer")

	//开启GateServer监听
	startGateway()

	//保持进程
	global.Run()
}

func getPort() {
	//端口号
	gateway_ip = cfg.GetValue("gateway_ip")
	gateway_port = cfg.GetValue("gateway_port")
	global.ServerName = "GateServer[" + gateway_port + "]"

	transfer_ip = cfg.GetValue("transfer_ip")
	transfer_port = cfg.GetValue("transfer_port")
}

func startGateway() {
	msgDispatch := dispatch.NewDispatch(
		dispatch.HandleFunc{
			H: transferProxy.SendToGameServer,
		},
	)

	addr := "0.0.0.0:" + gateway_port
	err := global.Listener("tcp", addr, global.PackCodecType_UnSafe,
		func(session *link.Session) {
			//将此Session记录在缓存内，消息回传时使用
			global.AddSession(session)
			//通知LoginServer用户上线
			transferProxy.SetClientSessionOnline(session)
			//添加session关闭时回调
			session.AddCloseCallback(session, func() {
				//通知LoginServer、GameServer用户下线
				transferProxy.SetClientSessionOffline(session.Id())
			})
		},
		msgDispatch,
	)

	checkError(err)
}

func checkError(err error) {
	if err != nil {
		ERR("Fatal error: %v", err)
		os.Exit(-1)
	}
}
