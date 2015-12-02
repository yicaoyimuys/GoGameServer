package worldProxy

import (
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"global"
	"module"
	"protos/gameProto"
	"protos/systemProto"
	//	. "tools"
)

var (
	servers map[uint32]*link.Session
)

//初始化
func InitServer(port string) error {
	servers = make(map[uint32]*link.Session)

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
	case systemProto.ID_System_ClientSessionOnlineC2S:
		setSessionOnline(session, protoMsg)
	case systemProto.ID_System_ClientSessionOfflineC2S:
		setSessionOffline(protoMsg)
	case systemProto.ID_System_ClientLoginSuccessC2S:
		setClientLoginSuccess(protoMsg)
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
		if msgID%2 == 1 {
			//C2S消息，由WorldServer处理消息
			dealGameMsg(msg)
		}
	}
}

//处理游戏逻辑
func dealGameMsg(msg packet.RAW) {
	msgIdentification := binary.GetUint64LE(msg[2:10])

	userSession := global.GetSession(msgIdentification)
	if userSession == nil {
		return
	}

	conn := userSession.Conn().(*WorldProxyConn)
	conn.recvChan <- msg
}

//在World服务器设置用户登录成功
func setClientLoginSuccess(protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientLoginSuccessC2S)
	userSession := global.GetSession(rev_msg.GetSessionID())
	if userSession != nil {
		module.User.LoginSuccess(userSession, rev_msg.GetUserName(), rev_msg.GetUserID())
	}
}

//在World服务端创建虚拟用户
func setSessionOnline(session *link.Session, protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientSessionOnlineC2S)

	userConn := NewWorldProxyConn(rev_msg.GetSessionID(), clientAddr{[]byte(rev_msg.GetNetwork()), []byte(rev_msg.GetAddr())}, session)
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

//在World服务端删除虚拟用户
func setSessionOffline(protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientSessionOfflineC2S)
	userSession := global.GetSession(rev_msg.GetSessionID())
	if userSession != nil {
		module.User.Offline(userSession)
	}
}
