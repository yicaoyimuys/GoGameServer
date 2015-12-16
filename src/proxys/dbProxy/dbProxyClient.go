package dbProxy

import (
	"github.com/funny/link"
	"github.com/funny/binary"
	"global"
	"protos"
	"protos/systemProto"
	"protos/dbProto"
	. "tools"
	"tools/dispatch"
)

var (
	dbClient 				 *link.Session
	clientMsgDispatch   	 dispatch.DispatchInterface
)

func init()  {
	handle := dispatch.NewHandleConditions()
	//系统消息
	handle.Add(dispatch.HandleCondition{
		Condition: systemProto.IsValidID,
		H: dispatch.Handle{
			systemProto.ID_System_ConnectDBServerS2C:		connectDBServerCallBack,
		},
	})
	//DB消息
	handle.Add(dispatch.HandleFuncCondition{
		Condition: dbProto.IsValidID,
		H: func(session *link.Session, msg []byte) {
			identification := binary.GetUint64LE(msg[2:10])
			var userSession *link.Session = global.GetSession(identification)
			if userSession == nil {
				return
			}
			dbMsgChan := userSession.State.(chan []byte)
			dbMsgChan <- msg
		},
	})

	//创建消息分派
	clientMsgDispatch = dispatch.NewDispatch(handle)
}

//初始化
func InitClient(ip string, port string) error {
	//连接DB服务器
	addr := ip + ":" + port
	client, err := global.Connect("DBServer", "tcp", addr, global.PackCodecType_Safe, clientMsgDispatch)
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
	send_msg := protos.MarshalProtoMsg(&systemProto.System_ConnectDBServerC2S{
		ServerName: protos.String(global.ServerName),
	})
	protos.Send(dbClient, send_msg)
}

//连接DB服务器返回
func connectDBServerCallBack(session *link.Session, protoMsg protos.ProtoMsg) {
	//	rev_msg := protoMsg.Body.(*systemProto.System_ConnectDBServerS2C)
	INFO(global.ServerName + " Connect DBServer Success")
}
