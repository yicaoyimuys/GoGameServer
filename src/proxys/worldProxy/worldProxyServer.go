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
	"tools/dispatch"
	"proxys"
	"proxys/gameProxy"
)

var (
	servers 			map[uint32]*link.Session
	serverMsgDispatch 	dispatch.DispatchInterface
)


func init()  {
	handle := dispatch.NewHandleConditions()
	//系统消息处理
	handle.Add(dispatch.HandleCondition{
		Condition: systemProto.IsValidID,
		H: dispatch.Handle{
			systemProto.ID_System_ConnectWorldServerC2S:		connectWorldServer,
			systemProto.ID_System_ClientSessionOfflineC2S:		setSessionOffline,
			systemProto.ID_System_ClientLoginSuccessC2S:		setClientLoginSuccess,
		},
	})
	//游戏消息处理
	handle.Add(dispatch.HandleFuncCondition{
		Condition: gameProto.IsValidID,
		H: func(session *link.Session, msg []byte) {
			dealGameMsg(msg)
		},
	})

	//创建消息分派
	serverMsgDispatch = dispatch.NewDispatch(handle)
}

//初始化
func InitServer(port string) error {
	servers = make(map[uint32]*link.Session)

	//监听tcp连接
	addr := "0.0.0.0:" + port
	err := global.Listener("tcp", addr, global.PackCodecType_Safe,
		func(session *link.Session) { },
		serverMsgDispatch,
	)

	return err
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
func setClientLoginSuccess(session *link.Session, protoMsg protos.ProtoMsg) {
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
			gameProxy.MsgDispatch.Process(userSession, msg)
		}
	}()
	module.User.LoginSuccess(userSession, rev_msg.GetUserName(), rev_msg.GetUserID(), rev_msg.GetGameServerID())
}

//在World服务端删除虚拟用户
func setSessionOffline(session *link.Session, protoMsg protos.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientSessionOfflineC2S)
	userSession := global.GetSession(rev_msg.GetSessionID())
	if userSession != nil {
		userSession.Close()
	}
}

//其他客户端连接WorldServer处理
func connectWorldServer(session *link.Session, protoMsg protos.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ConnectWorldServerC2S)

	serverName := rev_msg.GetServerName()
	serverID := rev_msg.GetServerID()
	servers[serverID] = session

	//GameServer断开连接处理
	session.AddCloseCallback(session, func(){
		delete(servers, serverID)
		ERR(serverName + " Disconnect At " + global.ServerName)
	})

	send_msg := protos.MarshalProtoMsg(&systemProto.System_ConnectWorldServerS2C{})
	protos.Send(session, send_msg)
}
