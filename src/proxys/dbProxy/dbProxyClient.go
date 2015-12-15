package dbProxy

import (
	"github.com/funny/link"
	"global"
	"protos"
	"protos/systemProto"
	"protos/dbProto"
	. "tools"
	"tools/dispatch"
)

var (
	dbClient 				 *link.Session
	clientMsgReceiveChan 	 dispatch.ReceiveMsgChan
	clientMsgDispatchAsync   dispatch.DispatchInterface
)

func init()  {
	//创建异步接收消息的Chan
	clientMsgReceiveChan = make(chan dispatch.ReceiveMsg, 4096)

	//创建消息分派
	clientMsgDispatchAsync = dispatch.NewDispatchAsync([]dispatch.ReceiveMsgChan{clientMsgReceiveChan},
		dispatch.Handle{
			systemProto.ID_System_ConnectDBServerS2C:		connectDBServerCallBack,
			dbProto.ID_DB_User_LoginS2C:					userLoginCallBack,
		},
	)
}

//初始化
func InitClient(ip string, port string) error {
	//连接DB服务器
	addr := ip + ":" + port
	client, err := global.Connect("DBServer", "tcp", addr, global.PackCodecType_Safe, clientMsgDispatchAsync)
	if err != nil {
		return err
	}

	dbClient = client
	sendConnectDBServer()

	return nil
}

//发送DB消息到服务器
func sendDBMsgToServer(msg []byte) {
	if dbClient == nil {
		serverMsgDispatch.Process(nil, msg)
		serverMsgDispatchAsync.Process(nil, msg)
	} else {
		protos.Send(dbClient, msg)
	}
}

//发送连接DB服务器
func sendConnectDBServer() {
	INFO(global.ServerName + " Connect DBServer ...")
	send_msg := systemProto.MarshalProtoMsg(&systemProto.System_ConnectDBServerC2S{
		ServerName: protos.String(global.ServerName),
	})
	protos.Send(dbClient, send_msg)
}

//连接DB服务器返回
func connectDBServerCallBack(session *link.Session, protoMsg protos.ProtoMsg) {
	//	rev_msg := protoMsg.Body.(*systemProto.System_ConnectDBServerS2C)
	INFO(global.ServerName + " Connect DBServer Success")
}
