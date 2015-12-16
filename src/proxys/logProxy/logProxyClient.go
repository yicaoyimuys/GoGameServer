package logProxy

import (
	"github.com/funny/link"
	"global"
	"protos"
	"protos/systemProto"
	. "tools"
	"tools/dispatch"
)

var (
	logClient 				 *link.Session
	clientMsgDispatch 		 dispatch.DispatchInterface
)

func init()  {
	//创建消息分派
	clientMsgDispatch = dispatch.NewDispatch(
		dispatch.Handle{
			systemProto.ID_System_ConnectLogServerS2C:		connectLogServerCallBack,
		},
	)
}

//初始化
func InitClient(ip string, port string) error {
	//连接LogServer
	addr := ip + ":" + port
	client, err := global.Connect("LogServer", "tcp", addr, global.PackCodecType_Async, clientMsgDispatch)
	if err != nil {
		return err
	}

	logClient = client
	sendConnectLogServer()

	return nil
}

//发送系统消息到LogServer
func sendSystemMsgToServer(msg []byte) {
	if logClient == nil {
		return
	}
	protos.Send(logClient, msg)
}

//发送Log消息到LogServer
func sendLogMsgToServer(msg []byte) {
	if logClient == nil {
		serverMsgDispatch.Process(nil, msg)
	} else {
		protos.Send(logClient, msg)
	}
}

//发送连接LogServer
func sendConnectLogServer() {
	INFO(global.ServerName + " Connect LogServer ...")
	send_msg := protos.MarshalProtoMsg(&systemProto.System_ConnectLogServerC2S{
		ServerName: protos.String(global.ServerName),
	})
	sendSystemMsgToServer(send_msg)
}

//连接Transfer服务器返回
func connectLogServerCallBack(session *link.Session, protoMsg protos.ProtoMsg) {
	//	rev_msg := protoMsg.Body.(*systemProto.System_ConnectLogServerS2C)
	INFO(global.ServerName + " Connect LogServer Success")
}
