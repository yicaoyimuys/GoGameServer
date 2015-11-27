package transferProxy

import (
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"protos/gameProto"
	_ "protos/gameProto"
	"protos/systemProto"
	_ "protos/systemProto"
	"strconv"
	"strings"
	. "tools"
	"tools/hashs"
)

var (
	servers          map[string][]*link.Session
	gameConsistent   *hashs.Consistent
	gameUserSessions map[uint64]int
)

//初始化
func InitServer(port string) error {
	servers = make(map[string][]*link.Session)
	gameConsistent = hashs.NewConsistent()
	gameUserSessions = make(map[uint64]int)

	listener, err := link.Serve("tcp", "0.0.0.0:"+port, packet.New(
		binary.SplitByUint32BE, 1024, 1024, 1024,
	))
	if err != nil {
		return err
	}

	listener.Serve(func(session *link.Session) {
		var msg packet.RAW
		for {
			if err := session.Receive(&msg); err != nil {
				break
			}
			dealReceiveMsgC2S(session, msg)
		}
	})

	return nil
}

//开始处理游戏逻辑消息
func startDealReceiveMsgC2S(session *link.Session) {
	revMsgChan := make(chan *packet.RAW, 2048)
	go func() {
		for {
			data, ok := <-revMsgChan
			if !ok {
				return
			}
			dealReceiveMsgC2S(session, *data)
		}
	}()

	for {
		var msg packet.RAW
		if err := session.Receive(&msg); err != nil {
			break
		}
		revMsgChan <- &msg
	}
}

//处理接收到的系统消息
func dealReceiveSystemMsgC2S(session *link.Session, msg packet.RAW) {
	protoMsg := systemProto.UnmarshalProtoMsg(msg)
	if protoMsg == systemProto.NullProtoMsg {
		return
	}

	switch protoMsg.ID {
	case systemProto.ID_System_ConnectTransferServerC2S:
		connectTransferServer(session, protoMsg)
	case systemProto.ID_System_ClientSessionOnlineC2S:
		sendSystemMsg("GameServer", msg)
		sendSystemMsg("LoginServer", msg)
	case systemProto.ID_System_ClientSessionOfflineC2S:
		sendSystemMsg("GameServer", msg)
		sendSystemMsg("LoginServer", msg)
		noAllotGameServer(protoMsg)
	case systemProto.ID_System_ClientLoginSuccessC2S:
		sendSystemMsg("GameServer", msg)
		allotGameServer(protoMsg)
	}
}

//处理接收到的消息
func dealReceiveMsgC2S(session *link.Session, msg packet.RAW) {
	if len(msg) < 2 {
		return
	}

	msgID := binary.GetUint16LE(msg[:2])
	if systemProto.IsValidID(msgID) {
		dealReceiveSystemMsgC2S(session, msg)
	} else if gameProto.IsValidID(msgID) {
		if msgID%2 == 1 {
			//C2S消息，发送到GameServer或者LoginServer
			if gameProto.IsValidLoginID(msgID) {
				sendGameMsg("LoginServer", msg)
			} else {
				sendGameMsg("GameServer", msg)
			}
		} else {
			//S2C消息，发送到GateServer
			sendGameMsg("GateServer", msg)
		}
	}
}

//不再分配游戏服务器
func noAllotGameServer(protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientSessionOfflineC2S)
	clientSessionID := rev_msg.GetSessionID()
	if _, exists := gameUserSessions[clientSessionID]; exists {
		delete(gameUserSessions, clientSessionID)
	}
}

//分配游戏服务器
func allotGameServer(protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientLoginSuccessC2S)

	clientSessionID := rev_msg.GetSessionID()
	if gameNode, existsNode := gameConsistent.Get(strconv.FormatUint(clientSessionID, 10)); existsNode {
		nodeIndex := gameNode.Id - 1
		gameUserSessions[clientSessionID] = nodeIndex
	}
}

//发送系统消息到不同服务器
func sendSystemMsg(serverName string, msg packet.RAW) {
	//系统消息发送到所有服务器
	if sessions, exists := servers[serverName]; exists {
		for _, session := range sessions {
			session.Send(msg)
		}
	}
}

//发送游戏消息到不同服务器
func sendGameMsg(serverName string, msg packet.RAW) {
	if sessions, exists := servers[serverName]; exists {
		if serverName == "GameServer" {
			//游戏消息发送到用户对应的GameServer
			clientSessionID := binary.GetUint64LE(msg[2:10])
			if nodeIndex, existsNodeID := gameUserSessions[clientSessionID]; existsNodeID {
				session := sessions[nodeIndex]
				session.Send(msg)
			}
		} else {
			for _, session := range sessions {
				session.Send(msg)
			}
		}
	}
}

//其他客户端连接TransferServer处理
func connectTransferServer(session *link.Session, protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ConnectTransferServerC2S)

	serverName := rev_msg.GetServerName()
	INFO(serverName + " Connect TransferServer")

	useServerName := strings.Split(serverName, "[")[0]
	sessions, exists := servers[useServerName]
	if !exists {
		sessions = make([]*link.Session, 0, 10)
	}
	sessions = append(sessions, session)
	servers[useServerName] = sessions

	//GameServer可以有多个
	if useServerName == "GameServer" {
		addr := strings.Split(session.Conn().RemoteAddr().String(), ":")
		addrIp := addr[0]
		addrPort, _ := strconv.Atoi(addr[1])
		gameConsistent.Add(hashs.NewNode(len(sessions), addrIp, addrPort, serverName, 1))
	}

	send_msg := systemProto.MarshalProtoMsg(&systemProto.System_ConnectTransferServerS2C{})
	systemProto.Send(send_msg, session)

	startDealReceiveMsgC2S(session)
}
