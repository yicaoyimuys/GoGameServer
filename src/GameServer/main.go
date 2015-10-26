package main

import (
	"flag"
	"os"
	"strconv"
)

import (
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/gateway"
	"github.com/funny/link/packet"
	//	"module"
	"protos"
	. "tools"
	"tools/cfg"
	_ "tools/db"
)

//各个模块
import (
	_ "module/cache"
	_ "module/config"
	_ "module/user"
)

var (
	game_ip   string
	game_port string
)

func main() {
	defer func() {
		if x := recover(); x != nil {
			ERR("caught panic in main()", x)
		}
	}()

	// 获取监听端口
	getPort()

	// 开启Log记录
	SetLogFile("game_log")
	SetLogPrefix("GameServer[" + game_port + "]")

	// 信号量监听
	go SignalProc()

	// 开启系统环境信息监测
	go SysRoutine()

	// 开启游戏服务器
	INFO("Starting GameServer : " + game_port)
	startGameServer()
}

func getPort() {
	var s int
	flag.IntVar(&s, "s", 0, "tcp listen port")
	flag.Parse()
	if s == 0 {
		ERR("please set gameserver port")
		os.Exit(-1)
	}
	game_ip = cfg.GetValue("game_ip_" + strconv.Itoa(s))
	game_port = cfg.GetValue("game_port_" + strconv.Itoa(s))
}

func pingGateway() {
	//尝试链接网关服务器，通知网关服务器，该游戏服务器已上线，动态增加游戏服务器时使用
	gateway_ip := cfg.GetValue("gateway_ip")
	gateway_port_by_game := cfg.GetValue("gateway_port_by_game")
	gateway_addr := gateway_ip + ":" + gateway_port_by_game
	if client, err := link.Connect("tcp", gateway_addr, packet.New(binary.SplitByUint32BE, 1024, 1024, 1024)); err == nil {
		simple.SendGamePingGateway(game_ip, game_port, client)
	}
}

func startGameServer() {
	backend, err := link.Serve("tcp", "0.0.0.0:"+game_port, gateway.NewBackend(1024, 1024, 1024))
	checkError(err)

	pingGateway()

	backend.Serve(func(session *link.Session) {
		simple.SendConnectSuccess(session)

		var msg packet.RAW
		for {
			if err := session.Receive(&msg); err != nil {
				break
			}
			dealMsg(session, msg)
		}
	})
}

func dealMsg(session *link.Session, msg packet.RAW) {
	msgID := binary.GetUint16LE(msg[:2])
	msgBody := msg[2:]

	//	DEBUG("收到消息ID: " + strconv.Itoa(int(msgID)))

	switch msgID {
	case simple.ID_GetUserInfoC2S:
		simple.GetUserInfo(msgBody, session)
	case simple.ID_AgainConnectC2S:
		simple.AgainConnect(msgBody, session)
	}
}

func checkError(err error) {
	if err != nil {
		ERR("Fatal error: %v", err)
		os.Exit(-1)
	}
}
