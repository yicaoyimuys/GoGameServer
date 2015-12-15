package transferProxy

import (
	"github.com/funny/link"
	"github.com/funny/binary"
	"global"
	"protos"
	"protos/gameProto"
	"protos/systemProto"
	"strconv"
	"strings"
	. "tools"
	"tools/hashs"
	"tools/dispatch"
)

var (
	servers          	map[string][]Server
	gameConsistent   	*hashs.Consistent
	gameUserSessions 	map[uint64]int
	serverMsgDispatch 	dispatch.DispatchInterface
)

type Server struct {
	session     *link.Session
	serverID    uint32
	serverIndex int
}

func init()  {
	handle := dispatch.NewHandleConditions()
	//系统消息处理
	handle.Add(dispatch.HandleCondition{
		Condition: systemProto.IsValidID,
		H: dispatch.Handle{
			systemProto.ID_System_ConnectTransferServerC2S:		connectTransferServer,
			systemProto.ID_System_ClientLoginSuccessC2S:		clientLoginSuccess,
		},
	})
	//游戏消息处理
	handle.Add(dispatch.HandleFuncCondition{
		Condition: gameProto.IsValidID,
		H: func(session *link.Session, msg []byte) {
			sendToGateServer(msg)
		},
	})

	//创建消息分派
	serverMsgDispatch = dispatch.NewDispatch(handle)
}

//初始化
func InitServer(port string) error {
	servers = make(map[string][]Server)
	gameConsistent = hashs.NewConsistent()
	gameUserSessions = make(map[uint64]int)

	//监听tcp
	addr := "0.0.0.0:" + port
	err := global.Listener("tcp", addr, global.PackCodecType_Safe,
		func(session *link.Session) {},
		serverMsgDispatch,
	)

	return err
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
func sendSystemMsg(serverName string, msg []byte) {
	//系统消息发送到所有服务器
	if sessions, exists := servers[serverName]; exists {
		for _, s := range sessions {
			s.session.Send(msg)
		}
	}
}

//发送系统消息到指定服务器
func sendSystemMsg2(serverName string, gameServerID uint32, msg []byte) {
	if serverList, exists := servers[serverName]; exists {
		for _, s := range serverList {
			if s.serverID == 0 || s.serverID == gameServerID {
				s.session.Send(msg)
			}
		}
	}
}

//发送游戏消息到不同服务器
func sendGameMsg(serverName string, msg []byte) {
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
func connectTransferServer(session *link.Session, protoMsg protos.ProtoMsg) {
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

	//服务器断开连接处理
	session.AddCloseCallback(session, func(){
		serverList = append(serverList[:server.serverIndex], serverList[server.serverIndex+1:]...)
		servers[useServerName] = serverList
		ERR(serverName + " Disconnect At " + global.ServerName)
	})

	//GameServer可以有多个
	if useServerName == "GameServer" {
		addr := strings.Split(session.Conn().RemoteAddr().String(), ":")
		addrIp := addr[0]
		addrPort, _ := strconv.Atoi(addr[1])
		gameNode := hashs.NewNode(server.serverIndex, addrIp, addrPort, serverName, 1)
		gameConsistent.Add(gameNode)

		//GameServer断开连接处理
		session.AddCloseCallback(session, func(){
			//移除此Node
			gameConsistent.Remove(gameNode)
			//将此Node的所有用户断开连接
			for clientSessionID, gameNodeIndex := range gameUserSessions {
				if server.serverIndex == gameNodeIndex {
					clientSession := global.GetSession(clientSessionID)
					if clientSession != nil {
						clientSession.Close()
					}
				}
			}
		})
	}

	//发送连接成功消息
	send_msg := systemProto.MarshalProtoMsg(&systemProto.System_ConnectTransferServerS2C{})
	protos.Send(session, send_msg)
}

//处理玩家登录成功，分配服务器
func clientLoginSuccess(session *link.Session, protoMsg protos.ProtoMsg) {
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
	sendSystemMsg2("GameServer", gameServerID, send_msg)
}

//LoginServer用户上线
func SetClientSessionOnline(userSession *link.Session) {
	//发送用户上线消息到serverName
	protoMsg := &systemProto.System_ClientSessionOnlineC2S{
		SessionID:    protos.Uint64(userSession.Id()),
		Network:      protos.String(userSession.Conn().RemoteAddr().Network()),
		Addr:         protos.String(userSession.Conn().RemoteAddr().String()),
	}
	send_msg := systemProto.MarshalProtoMsg(protoMsg)
	sendSystemMsg2("LoginServer", 0, send_msg)
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

	sendSystemMsg2("GameServer", gameServerID, send_msg)
	sendSystemMsg2("LoginServer", 0, send_msg)

	//给该用户不再分配游戏服务器
	noAllotGameServer(sessionID)
}

//发送消息到TransferServer, 网关调用
func SendToGameServer(userSession *link.Session, msg []byte) {
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
func sendToGateServer(msg []byte) {
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
	userSession.Send(result)
}
