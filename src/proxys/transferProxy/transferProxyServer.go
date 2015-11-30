package transferProxy

import (
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"global"
	"protos"
	"protos/gameProto"
	"protos/systemProto"
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

	go func() {
		listener.Serve(func(session *link.Session) {
			var msg packet.RAW
			for {
				if err := session.Receive(&msg); err != nil {
					break
				}
				dealReceiveMsgC2S(session, msg)
			}
		})
	}()

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
	case systemProto.ID_System_ClientLoginSuccessC2S:
		sendSystemMsg("GameServer", msg)
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
		if msgID%2 == 0 {
			//S2C消息，发送到GateServer
			SendToGateServer(msg)
		}
	}
}

//不再分配游戏服务器
func noAllotGameServer(clientSessionID uint64) {
	if _, exists := gameUserSessions[clientSessionID]; exists {
		delete(gameUserSessions, clientSessionID)
	}
}

//分配游戏服务器
func allotGameServer(clientSessionID uint64) {
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

//通知游戏服务器用户上线, 网关调用
func SetClientSessionOnline(userSession *link.Session) {
	global.AddSession(userSession)

	protoMsg := &systemProto.System_ClientSessionOnlineC2S{
		SessionID: protos.Uint64(userSession.Id()),
		Network:   protos.String(userSession.Conn().RemoteAddr().Network()),
		Addr:      protos.String(userSession.Conn().RemoteAddr().String()),
	}
	send_msg := systemProto.MarshalProtoMsg(protoMsg)

	sendSystemMsg("GameServer", packet.RAW(send_msg))
	sendSystemMsg("LoginServer", packet.RAW(send_msg))

	//给该用户分配游戏服务器
	allotGameServer(userSession.Id())
}

//通知游戏服务器用户下线, 网关调用
func SetClientSessionOffline(sessionID uint64) {
	protoMsg := &systemProto.System_ClientSessionOfflineC2S{
		SessionID: protos.Uint64(sessionID),
	}
	send_msg := systemProto.MarshalProtoMsg(protoMsg)

	sendSystemMsg("GameServer", packet.RAW(send_msg))
	sendSystemMsg("LoginServer", packet.RAW(send_msg))

	//给该用户不再分配游戏服务器
	noAllotGameServer(sessionID)
}

//发送消息到TransferServer, 网关调用
func SendToGameServer(msg packet.RAW, userSession *link.Session) {
	send_msg := make([]byte, 8+len(msg))
	copy(send_msg[:2], msg[:2])
	binary.PutUint64LE(send_msg[2:10], userSession.Id())
	copy(send_msg[10:], msg[2:])

	//C2S消息，发送到GameServer或者LoginServer
	msgID := binary.GetUint16LE(send_msg[:2])
	if gameProto.IsValidLoginID(msgID) {
		sendGameMsg("LoginServer", send_msg)
	} else {
		sendGameMsg("GameServer", send_msg)
	}
}

//发送消息到用户客户端
func SendToGateServer(msg packet.RAW) {
	if len(msg) < 10 {
		return
	}

	msgID := binary.GetUint16LE(msg[:2])
	msgIdentification := binary.GetUint64LE(msg[2:10])
	msgBody := msg[10:]

	userSession := global.GetSession(msgIdentification)
	if userSession == nil {
		return
	}

	result := make([]byte, len(msg)-8)
	binary.PutUint16LE(result[:2], msgID)
	copy(result[2:], msgBody)
	userSession.Send(packet.RAW(result))
}
