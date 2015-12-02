package transferProxy

import (
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"global"
	"module"
	"protos"
	"protos/gameProto"
	"protos/systemProto"
	"proxys/worldProxy"
	. "tools"
)

var (
	session *link.Session
)

//初始化
func InitClient(ip string, port string) error {
	addr := ip + ":" + port
	client, err := link.Connect("tcp", addr, packet.New(binary.SplitByUint32BE, 1024, 1024, 1024))
	if err != nil {
		return err
	}

	session = client
	go dealReceiveMsg()
	ConnectTransferServer()

	return nil
}

//处理从TransferServer发回的消息
func dealReceiveMsg() {
	for {
		var msg packet.RAW
		if err := session.Receive(&msg); err != nil {
			break
		}
		dealReceiveMsgS2C(msg)
	}
}

//处理接收到的系统消息
func dealReceiveSystemMsgS2C(msg packet.RAW) {
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
	case systemProto.ID_System_ClientLoginSuccessS2C:
		setClientLoginSuccess(protoMsg)
	}
}

//处理接收到的消息
func dealReceiveMsgS2C(msg packet.RAW) {
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
	if session == nil {
		return
	}
	protos.Send(msg, session)
}

//发送连接TransferServer
func ConnectTransferServer() {
	send_msg := systemProto.MarshalProtoMsg(&systemProto.System_ConnectTransferServerC2S{
		ServerName: protos.String(global.ServerName),
		ServerID:   protos.Uint32(global.ServerID),
	})
	sendSystemMsgToServer(send_msg)
}

//连接DB服务器返回
func connectTransferServerCallBack(protoMsg systemProto.ProtoMsg) {
	//	rev_msg := protoMsg.Body.(*systemProto.System_ConnectTransferServerS2C)
	INFO(global.ServerName + " Connect TransferServer Success")
}

//通知游戏服务器用户登录成功
func SetClientLoginSuccess(userName string, userID uint64, sessionID uint64) {
	send_msg := systemProto.MarshalProtoMsg(&systemProto.System_ClientLoginSuccessC2S{
		UserID:    protos.Uint64(userID),
		UserName:  protos.String(userName),
		SessionID: protos.Uint64(sessionID),
	})
	sendSystemMsgToServer(send_msg)
}

//在游戏服务器设置用户登录成功
func setClientLoginSuccess(protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientLoginSuccessS2C)
	userSession := global.GetSession(rev_msg.GetSessionID())
	if userSession != nil {
		module.User.LoginSuccess(userSession, rev_msg.GetUserName(), rev_msg.GetUserID())
	}
}

//在游戏服务端创建虚拟用户
func setSessionOnline(protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientSessionOnlineC2S)
	userConn := NewTransferProxyConn(rev_msg.GetSessionID(), clientAddr{[]byte(rev_msg.GetNetwork()), []byte(rev_msg.GetAddr())}, session)
	userSession := link.NewSession(rev_msg.GetSessionID(), userConn)
	go func() {
		var msg packet.RAW
		for {
			if err := userSession.Receive(&msg); err != nil {
				break
			}
			module.ReceiveMessage(userSession, msg)
		}
	}()
	module.User.Online(userSession)
}

//在游戏服务端删除虚拟用户
func setSessionOffline(protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientSessionOfflineC2S)
	userSession := global.GetSession(rev_msg.GetSessionID())
	if userSession != nil {
		module.User.Offline(userSession)
	}
}

//处理游戏逻辑
func dealGameMsg(msg packet.RAW) {
	msgIdentification := binary.GetUint64LE(msg[2:10])

	userSession := global.GetSession(msgIdentification)
	if userSession == nil {
		return
	}

	conn := userSession.Conn().(*TransferProxyConn)
	conn.recvChan <- msg
}
