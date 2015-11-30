package main

import (
	"os"
)

import (
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/packet"
	//	"module"
	"global"
	"proxys/transferProxy"
	. "tools"
	"tools/cfg"
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

	//开启TransferProxy，由GateServer充当中转服务器
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
	listener, err := link.Serve("tcp", "0.0.0.0:"+gateway_port, packet.New(
		binary.SplitByUint32BE, 1024, 1024, 1024,
	))
	checkError(err)

	listener.Serve(func(session *link.Session) {
		session.AddCloseCallback(session, func() {
			transferProxy.SetClientSessionOffline(session.Id())
		})
		transferProxy.SetClientSessionOnline(session)

		var msg packet.RAW
		for {
			if err := session.Receive(&msg); err != nil {
				break
			}
			transferProxy.SendToGameServer(msg, session)
		}
	})
}

func checkError(err error) {
	if err != nil {
		ERR("Fatal error: %v", err)
		os.Exit(-1)
	}
}
