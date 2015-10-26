package main

import (
	"errors"
	"os"
	"strconv"
)

import (
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/gateway"
	"github.com/funny/link/packet"
	"github.com/funny/link/stream"
	//	"module"
	"protos"
	. "tools"
	"tools/cfg"
	"tools/hashs"
)

//各个模块
import (
	_ "module/cache"
	_ "module/user"
)

var (
	consistent           *hashs.Consistent
	frontend             *gateway.Frontend
	gateway_port         string
	gateway_port_by_game string
)

func main() {
	defer func() {
		if x := recover(); x != nil {
			ERR("caught panic in main()", x)
		}
	}()

	// 开启Log记录
	SetLogFile("gateway_log")
	SetLogPrefix("GateServer")

	// 信号量监听
	go SignalProc()

	// 开启系统环境信息监测
	go SysRoutine()

	// 获取端口号
	getPort()

	// 开启网关服务器
	INFO("Starting Gateway: " + gateway_port)
	startGateway()

	// 开启对游戏服务器的侦听
	startGatewayByGameServer()

	//	// 保持进程
	//	temp := make(chan int32, 10)
	//	for {
	//		select {
	//		case <-temp:
	//		}
	//	}
}

func getPort() {
	//端口号
	gateway_port = cfg.GetValue("gateway_port")
	gateway_port_by_game = cfg.GetValue("gateway_port_by_game")
}

func startGateway() {
	//获取游戏服务器个数
	serverNum, err := strconv.Atoi(cfg.GetValue("init_game_num"))
	checkError(err)

	//开启客户端链接
	listener, err := link.Listen("tcp", "0.0.0.0:"+gateway_port, packet.New(
		binary.SplitByUint32BE, 1024, 1024, 1024,
	))
	checkError(err)

	//一致性哈希
	consistent = hashs.NewConsistent()

	//开启网关前端
	frontend = gateway.NewFrontend(listener.(*packet.Listener),
		func(session *link.Session) (uint64, error) {
			//分配逻辑服务器
			key := "session_" + strconv.FormatUint(session.Id(), 10)
			node, exists := consistent.Get(key)
			if exists {
				//			DEBUG(node.Id)
				return uint64(node.Id), nil
			} else {
				errStr := "No Has GameServer"
				ERR(errStr)
				return 0, errors.New(errStr)
			}
		},
	)

	//开启网关后端
	for i := 1; i <= serverNum; i++ {
		ia := strconv.Itoa(i)
		game_ip := cfg.GetValue("game_ip_" + ia)
		game_port := cfg.GetValue("game_port_" + ia)
		createBackend(game_ip, game_port)
	}
}

func startGatewayByGameServer() {
	listener, err := link.Serve("tcp", "0.0.0.0:"+gateway_port_by_game, packet.New(
		binary.SplitByUint32BE, 1024, 1024, 1024,
	))
	checkError(err)

	listener.Serve(func(session *link.Session) {
		var msg packet.RAW
		for {
			if err := session.Receive(&msg); err != nil {
				break
			}
			dealMsg(session, msg)
		}
	})
}

func dealMsg(client *link.Session, msg packet.RAW) bool {
	msgID := binary.GetUint16LE(msg[:2])
	msgBody := msg[2:]

	//	DEBUG("收到消息ID: " + strconv.Itoa(int(msgID)))

	//登录、注册消息在网关处理
	switch msgID {
	case simple.ID_UserLoginC2S:
		simple.Login(msgBody, client)
		return false
	case simple.ID_GamePingGateway:
		game_ip, game_port := simple.ReceiveGamePingGateway(msgBody)
		createBackend(game_ip, game_port)
		client.Close()
		return false
	}
	return true
}

func createBackend(ip string, port string) {
	addr := ip + ":" + port
	frontendLink, err := frontend.AddBackend("tcp",
		addr,
		stream.New(1024, 1024, 1024),
	)

	if err != nil {
		ERR("Connect GameServer Error", err)
		return
	}

	linkId := frontendLink.LinkId()
	nodeId := int(linkId)

	//Link接收到客户端发送的消息处理
	frontendLink.SetReceiveClientMsgHandle(dealMsg)

	//Link关闭时处理
	frontendLink.SetCloseHandle(func() {
		frontend.RemoveBackend(linkId)

		if node, exists := consistent.GetNodeByID(nodeId); exists {
			INFO("GameServer", node.Ip, node.Port, "退出集群")
			consistent.Remove(node)
		}
	})

	INFO("GameServer:", addr, "加入集群")

	//添加到节点
	weight := 1
	portNum, _ := strconv.Atoi(port)
	consistent.Add(hashs.NewNode(nodeId, ip, portNum, addr, weight))
}

func checkError(err error) {
	if err != nil {
		ERR("Fatal error: %v", err)
		os.Exit(-1)
	}
}
