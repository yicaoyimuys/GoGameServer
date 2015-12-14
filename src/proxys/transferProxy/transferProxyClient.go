package transferProxy

import (
	"github.com/funny/link"
	"github.com/funny/binary"
	"global"
	"module"
	"protos"
	"protos/gameProto"
	"protos/systemProto"
	"proxys/worldProxy"
	. "tools"
	"tools/codecType"
	"proxys"
)

var (
	transferClient *link.Session
)

//初始化
func InitClient(ip string, port string) error {
	addr := ip + ":" + port
	client, err := link.Connect("tcp", addr, global.PackCodecType)
	if err != nil {
		return err
	}
	client.AddCloseCallback(client, func(){
		ERR("TransferServer Disconnect At " + global.ServerName)
	})

	transferClient = client
	go dealReceiveMsg()
	ConnectTransferServer()

	return nil
}

//处理从TransferServer发回的消息
func dealReceiveMsg() {
	for {
		var msg []byte
		if err := transferClient.Receive(&msg); err != nil {
			break
		}
		dealReceiveMsgS2C(msg)
	}
}

//处理接收到的系统消息
func dealReceiveSystemMsgS2C(msg []byte) {
	protoMsg := systemProto.UnmarshalProtoMsg(msg)
	if protoMsg == systemProto.NullProtoMsg {
		return
	}

	switch protoMsg.ID {
	case systemProto.ID_System_ConnectTransferServerS2C:
		connectTransferServerCallBack(protoMsg)
	case systemProto.ID_System_ClientSessionOnlineC2S:
		setSessionOnline(protoMsg)
	case systemProto.ID_System_ClientSessionOfflineC2S:
		setSessionOffline(protoMsg)
	case systemProto.ID_System_ClientLoginSuccessC2S:
		setClientLoginSuccess(protoMsg)
	}
}

//处理接收到的消息
func dealReceiveMsgS2C(msg []byte) {
	if len(msg) < 2 {
		return
	}

	msgID := binary.GetUint16LE(msg[:2])
	//	DEBUG(global.ServerName, msgID)
	if systemProto.IsValidID(msgID) {
		dealReceiveSystemMsgS2C(msg)
	} else if gameProto.IsValidWorldID(msgID) {
		if msgID%2 == 1 {
			//C2S消息，由WorldServer处理
			worldProxy.SendGameMsgToServer(msg)
		}
	} else if gameProto.IsValidID(msgID) {
		if msgID%2 == 1 {
			//C2S消息，由GameServer处理
			dealGameMsg(msg)
		}
	} else {
		ERR(global.ServerName, "收到未处理消息")
	}
}

//发送系统消息到TransferServer
func sendSystemMsgToServer(msg []byte) {
	if transferClient == nil {
		return
	}
	protos.Send(msg, transferClient)
}

//发送连接TransferServer
func ConnectTransferServer() {
	INFO(global.ServerName + " Connect TransferServer ...")
	send_msg := systemProto.MarshalProtoMsg(&systemProto.System_ConnectTransferServerC2S{
		ServerName: protos.String(global.ServerName),
		ServerID:   protos.Uint32(global.ServerID),
	})
	sendSystemMsgToServer(send_msg)
}

//连接TransferServer返回
func connectTransferServerCallBack(protoMsg systemProto.ProtoMsg) {
	//	rev_msg := protoMsg.Body.(*systemProto.System_ConnectTransferServerS2C)
	INFO(global.ServerName + " Connect TransferServer Success")
}

//通知GameServer用户登录成功
func SetClientLoginSuccess(userName string, userID uint64, session *link.Session) {
	send_msg := systemProto.MarshalProtoMsg(&systemProto.System_ClientLoginSuccessC2S{
		UserID:    		protos.Uint64(userID),
		UserName:  		protos.String(userName),
		SessionID: 		protos.Uint64(session.Id()),
		GameServerID:	protos.Uint32(0),
		Network:      	protos.String(session.Conn().RemoteAddr().Network()),
		Addr:         	protos.String(session.Conn().RemoteAddr().String()),
	})
	sendSystemMsgToServer(send_msg)
}

//在GameServer设置用户登录成功
func setClientLoginSuccess(protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientLoginSuccessC2S)
	userConn := proxys.NewDummyConn(rev_msg.GetSessionID(), rev_msg.GetNetwork(), rev_msg.GetAddr(), transferClient)
	userSession := link.NewSessionByID(userConn, codecType.DummyCodecType{}, rev_msg.GetSessionID())
	global.AddSession(userSession)
	go func() {
		var msg []byte
		for {
			if err := userSession.Receive(&msg); err != nil {
				break
			}
			module.ReceiveMessage(userSession, msg)
		}
	}()
	module.User.LoginSuccess(userSession, rev_msg.GetUserName(), rev_msg.GetUserID(), rev_msg.GetGameServerID())

	//通知WorldServer用户登录成功
	worldProxy.SendSystemMsgToServer(systemProto.MarshalProtoMsg(rev_msg))
}

//在LoginServer创建虚拟用户
func setSessionOnline(protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientSessionOnlineC2S)
	userConn := proxys.NewDummyConn(rev_msg.GetSessionID(), rev_msg.GetNetwork(), rev_msg.GetAddr(), transferClient)
	userSession := link.NewSessionByID(userConn, codecType.DummyCodecType{}, rev_msg.GetSessionID())
	global.AddSession(userSession)
	go func() {
		var msg []byte
		for {
			if err := userSession.Receive(&msg); err != nil {
				break
			}
			module.ReceiveMessage(userSession, msg)
		}
	}()
}

//在游戏服务端删除虚拟用户
func setSessionOffline(protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientSessionOfflineC2S)
	userSession := global.GetSession(rev_msg.GetSessionID())
	if userSession != nil {
		userSession.Close()
	}

	//通知WorldServer用户下线
	worldProxy.SendSystemMsgToServer(systemProto.MarshalProtoMsg(rev_msg))
}

//处理游戏逻辑
func dealGameMsg(msg []byte) {
	msgIdentification := binary.GetUint64LE(msg[2:10])

	userSession := global.GetSession(msgIdentification)
	if userSession == nil {
		return
	}

	conn := userSession.Conn().(*proxys.DummyConn)
	conn.PutMsg(msg)
}
