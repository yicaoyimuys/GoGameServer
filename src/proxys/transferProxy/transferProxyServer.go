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
	//	. "tools"
	"tools/hashs"
)

var (
	servers          map[string][]Server
	gameConsistent   *hashs.Consistent
	gameUserSessions map[uint64]int
)

type Server struct {
	session     *link.Session
	serverID    uint32
	serverIndex int
}

//初始化
func InitServer(port string) error {
	servers = make(map[string][]Server)
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
		clientLoginSuccess(protoMsg)
	}
}

//处理接收到的消息
func dealReceiveMsgC2S(session *link.Session, msg packet.RAW) {
	if len(msg) < 2 {
		return
	}

	msgID := binary.GetUint16LE(msg[:2])
	if systemProto.IsValidID(msgID) {
		//系统消息
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
		gameUserSessions[clientSessionID] = gameNode.Id
	}
}

//获取给用户分配的GameServerID
func getUserGameServerID(sessionID uint64) uint32 {
	nodeIndex, existsNodeID := gameUserSessions[sessionID]
	if !existsNodeID {
		return 0
	}
	gameServerID := servers["GameServer"][nodeIndex].serverID
	return gameServerID
}

//发送系统消息到不同服务器
func sendSystemMsg(serverName string, msg packet.RAW) {
	//系统消息发送到所有服务器
	if sessions, exists := servers[serverName]; exists {
		for _, s := range sessions {
			s.session.Send(msg)
		}
	}
}

//发送系统消息到指定服务器
func sendSystemMsg2(serverName string, gameServerID uint32, msg packet.RAW) {
	if serverList, exists := servers[serverName]; exists {
		for _, s := range serverList {
			if s.serverID == 0 || s.serverID == gameServerID {
				s.session.Send(msg)
			}
		}
	}
}

//发送游戏消息到不同服务器
func sendGameMsg(serverName string, msg packet.RAW) {
	if serverList, exists := servers[serverName]; exists {
		if serverName == "GameServer" {
			//游戏消息发送到用户对应的GameServer
			clientSessionID := binary.GetUint64LE(msg[2:10])
			if nodeIndex, existsNodeID := gameUserSessions[clientSessionID]; existsNodeID {
				s := serverList[nodeIndex]
				s.session.Send(msg)
			}
		} else {
			for _, s := range serverList {
				s.session.Send(msg)
			}
		}
	}
}

//其他客户端连接TransferServer处理
func connectTransferServer(session *link.Session, protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ConnectTransferServerC2S)

	serverName := rev_msg.GetServerName()
	serverID := rev_msg.GetServerID()

	useServerName := strings.Split(serverName, "[")[0]
	serverList, exists := servers[useServerName]
	if !exists {
		serverList = make([]Server, 0, 10)
	}
	server := Server{
		session:     session,
		serverID:    serverID,
		serverIndex: len(serverList),
	}
	serverList = append(serverList, server)
	servers[useServerName] = serverList

	//GameServer可以有多个
	if useServerName == "GameServer" {
		addr := strings.Split(session.Conn().RemoteAddr().String(), ":")
		addrIp := addr[0]
		addrPort, _ := strconv.Atoi(addr[1])
		gameConsistent.Add(hashs.NewNode(server.serverIndex, addrIp, addrPort, serverName, 1))
	}

	send_msg := systemProto.MarshalProtoMsg(&systemProto.System_ConnectTransferServerS2C{})
	protos.Send(send_msg, session)
}

//处理玩家登录成功，分配服务器
func clientLoginSuccess(protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientLoginSuccessC2S)

	//用户SessionID
	userSessionID := rev_msg.GetSessionID()

	//给该用户分配GameServer
	allotGameServer(userSessionID)

	//给该用户所分配的GameServerID
	gameServerID := getUserGameServerID(userSessionID)
	if gameServerID == 0 {
		return
	}

	//给该消息填充上所分配的GameServerID
	rev_msg.GameServerID = protos.Uint32(gameServerID)

	//通知GameServer用户登录成功
	send_msg := systemProto.MarshalProtoMsg(rev_msg)
	sendSystemMsg2("GameServer", gameServerID, packet.RAW(send_msg))
}

//LoginServer用户上线
func SetClientSessionOnline(userSession *link.Session) {
	//给该用户所分配的GameServerID
	gameServerID := getUserGameServerID(userSession.Id())

	//发送用户上线消息到serverName
	protoMsg := &systemProto.System_ClientSessionOnlineC2S{
		SessionID:    protos.Uint64(userSession.Id()),
		Network:      protos.String(userSession.Conn().RemoteAddr().Network()),
		Addr:         protos.String(userSession.Conn().RemoteAddr().String()),
		GameServerID: protos.Uint32(gameServerID),
	}
	send_msg := systemProto.MarshalProtoMsg(protoMsg)

	sendSystemMsg2("LoginServer", gameServerID, packet.RAW(send_msg))
}

//通知GameServer、LoginServer用户下线, 网关调用
func SetClientSessionOffline(sessionID uint64) {
	//给该用户所分配的GameServerID
	gameServerID := getUserGameServerID(sessionID)

	//发送消息到GameServer和LoginServer
	protoMsg := &systemProto.System_ClientSessionOfflineC2S{
		SessionID: protos.Uint64(sessionID),
	}
	send_msg := systemProto.MarshalProtoMsg(protoMsg)

	sendSystemMsg2("GameServer", gameServerID, packet.RAW(send_msg))
	sendSystemMsg2("LoginServer", 0, packet.RAW(send_msg))

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
