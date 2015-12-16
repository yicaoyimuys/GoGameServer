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
	"tools/dispatch"
	"proxys"
	"proxys/dbProxy"
	"proxys/gameProxy"
)

var (
	transferClient 			 *link.Session
	clientMsgDispatch 		 dispatch.DispatchInterface
)

func init() {
	handle := dispatch.NewHandleConditions()
	//系统消息处理
	handle.Add(dispatch.HandleCondition{
		Condition: systemProto.IsValidID,
		H: dispatch.Handle{
			systemProto.ID_System_ConnectTransferServerS2C:		connectTransferServerCallBack,
			systemProto.ID_System_ClientSessionOnlineC2S:		setSessionOnline,
			systemProto.ID_System_ClientSessionOfflineC2S:		setSessionOffline,
			systemProto.ID_System_ClientLoginSuccessC2S:		setClientLoginSuccess,
		},
	})
	//LoginServer消息
	handle.Add(dispatch.HandleFuncCondition{
		Condition: gameProto.IsValidLoginID,
		H: func(session *link.Session, msg []byte) {
			dealGameMsg(msg)
		},
	})
	//GameServer消息
	handle.Add(dispatch.HandleFuncCondition{
		Condition: gameProto.IsValidGameID,
		H: func(session *link.Session, msg []byte) {
			dealGameMsg(msg)
		},
	})
	//WorldServer消息
	handle.Add(dispatch.HandleFuncCondition{
		Condition: gameProto.IsValidWorldID,
		H: func(session *link.Session, msg []byte) {
			worldProxy.SendGameMsgToServer(msg)
		},
	})

	//创建消息分派
	clientMsgDispatch = dispatch.NewDispatch(handle)
}

//初始化
func InitClient(ip string, port string) error {
	//连接TransferServer
	addr := ip + ":" + port
	client, err := global.Connect("TransferServer", "tcp", addr, global.PackCodecType_Safe, clientMsgDispatch)
	if err != nil {
		return err
	}

	transferClient = client
	sendConnectTransferServer()

	return nil
}

//发送系统消息到TransferServer
func sendSystemMsgToServer(msg []byte) {
	if transferClient == nil {
		return
	}
	protos.Send(transferClient, msg)
}

//发送连接TransferServer
func sendConnectTransferServer() {
	INFO(global.ServerName + " Connect TransferServer ...")
	send_msg := protos.MarshalProtoMsg(&systemProto.System_ConnectTransferServerC2S{
		ServerName: protos.String(global.ServerName),
		ServerID:   protos.Uint32(global.ServerID),
	})
	sendSystemMsgToServer(send_msg)
}

//连接TransferServer返回
func connectTransferServerCallBack(session *link.Session, protoMsg protos.ProtoMsg) {
	//	rev_msg := protoMsg.Body.(*systemProto.System_ConnectTransferServerS2C)
	INFO(global.ServerName + " Connect TransferServer Success")
}

//通知GameServer用户登录成功
func SetClientLoginSuccess(userName string, userID uint64, session *link.Session) {
	send_msg := protos.MarshalProtoMsg(&systemProto.System_ClientLoginSuccessC2S{
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
func setClientLoginSuccess(session *link.Session, protoMsg protos.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientLoginSuccessC2S)
	userConn := proxys.NewDummyConn(rev_msg.GetSessionID(), rev_msg.GetNetwork(), rev_msg.GetAddr(), transferClient)
	userSession := link.NewSessionByID(userConn, codecType.DummyCodecType{}, rev_msg.GetSessionID())
	global.AddSession(userSession)
	//接收游戏消息
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

	//通知WorldServer用户登录成功
	worldProxy.SendSystemMsgToServer(protos.MarshalProtoMsg(rev_msg))
}

//在LoginServer创建虚拟用户
func setSessionOnline(session *link.Session, protoMsg protos.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientSessionOnlineC2S)
	userConn := proxys.NewDummyConn(rev_msg.GetSessionID(), rev_msg.GetNetwork(), rev_msg.GetAddr(), transferClient)
	userSession := link.NewSessionByID(userConn, codecType.DummyCodecType{}, rev_msg.GetSessionID())
	userSession.State = make(chan []byte, 100)
	userSession.AddCloseCallback(userSession, func(){ close(userSession.State.(chan []byte)) })
	global.AddSession(userSession)
	//接收游戏消息
	go func() {
		var msg []byte
		for {
			if err := userSession.Receive(&msg); err != nil {
				break
			}
			gameProxy.MsgDispatch.Process(userSession, msg)
		}
	}()
	//接收DB消息
	go func() {
		dbMsgChan := userSession.State.(chan []byte)
		for {
			data, ok := <-dbMsgChan
			if !ok {
				return
			}
			dbProxy.ClientDbMsgDispatchHandle.DealMsg(userSession, data)
		}
	}()
}

//在游戏服务端删除虚拟用户
func setSessionOffline(session *link.Session, protoMsg protos.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ClientSessionOfflineC2S)
	userSession := global.GetSession(rev_msg.GetSessionID())
	if userSession != nil {
		userSession.Close()
	}

	//通知WorldServer用户下线
	worldProxy.SendSystemMsgToServer(protos.MarshalProtoMsg(rev_msg))
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
