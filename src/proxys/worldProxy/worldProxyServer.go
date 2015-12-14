package worldProxy

import (
	"github.com/funny/link"
	"github.com/funny/binary"
	"global"
	"module"
	"protos"
	"protos/gameProto"
	"protos/systemProto"
	. "tools"
	"tools/codecType"
	"proxys"
)

var (
	servers 		map[uint32]*link.Session
)

//初始化
func InitServer(port string) error {
	servers = make(map[uint32]*link.Session)

	err := global.Listener("tcp", "0.0.0.0:"+port, global.PackCodecType, func(session *link.Session) {
		var msg []byte
		for {
			if err := session.Receive(&msg); err != nil {
				break
			}
			dealReceiveMsgC2S(session, msg)
		}
	})

	if err != nil {
		return err
	}

	return nil
}

//处理接收到的系统消息
func dealReceiveSystemMsgC2S(session *link.Session, msg []byte) {
	protoMsg := systemProto.UnmarshalProtoMsg(msg)
	if protoMsg == systemProto.NullProtoMsg {
		return
	}

	switch protoMsg.ID {
	case systemProto.ID_System_ConnectWorldServerC2S:
		connectWorldServer(session, protoMsg)
	case systemProto.ID_System_ClientSessionOfflineC2S:
		setSessionOffline(protoMsg)
	case systemProto.ID_System_ClientLoginSuccessC2S:
		setClientLoginSuccess(session, protoMsg)
	}
}

//处理接收到的消息
func dealReceiveMsgC2S(session *link.Session, msg []byte) {
	if len(msg) < 2 {
		return
	}

	msgID := binary.GetUint16LE(msg[:2])
	if systemProto.IsValidID(msgID) {
		//系统消息
		dealReceiveSystemMsgC2S(session, msg)
	} else if gameProto.IsValidID(msgID) {
		if msgID%2 == 1 {
			//C2S消息，由WorldServer处理消息
			dealGameMsg(msg)
		}
	}
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

//在World服务器设置用户登录成功
func setClientLoginSuccess(session *link.Session, protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientLoginSuccessC2S)
	userConn := proxys.NewDummyConn(rev_msg.GetSessionID(), rev_msg.GetNetwork(), rev_msg.GetAddr(), session)
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
}

//在World服务端删除虚拟用户
func setSessionOffline(protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientSessionOfflineC2S)
	userSession := global.GetSession(rev_msg.GetSessionID())
	if userSession != nil {
		userSession.Close()
	}
}

//其他客户端连接WorldServer处理
func connectWorldServer(session *link.Session, protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ConnectWorldServerC2S)

	serverName := rev_msg.GetServerName()
	serverID := rev_msg.GetServerID()
	servers[serverID] = session

	//GameServer断开连接处理
	session.AddCloseCallback(session, func(){
		delete(servers, serverID)
		ERR(serverName + " Disconnect At " + global.ServerName)
	})

	send_msg := systemProto.MarshalProtoMsg(&systemProto.System_ConnectWorldServerS2C{})
	protos.Send(send_msg, session)
}
