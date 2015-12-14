package dbProxy

import (
	"github.com/funny/link"
	"global"
	"protos"
	"protos/systemProto"
	. "tools"
)

var (
	logClient *link.Session
)

//初始化
func InitClient(ip string, port string) error {
	//连接DB服务器
	addr := ip + ":" + port
	client, err := link.Connect("tcp", addr, global.PackCodecType)
	if err != nil {
		return err
	}
	client.AddCloseCallback(client, func(){
		ERR("DBServer Disconnect At " + global.ServerName)
	})

	logClient = client
	go dealReceiveMsgS2C()
	ConnectDBServer()

	return nil
}

//发送DB消息到服务器
func sendDBMsgToServer(msg []byte) {
	if logClient == nil {
		dealReceiveDBMsgC2S(logClient, msg)
		dealReceiveAsyncDBMsgC2S(msg)
	} else {
		protos.Send(msg, logClient)
	}
}

//处理从DBServer发回的消息
func dealReceiveMsgS2C() {
	var msg []byte
	for {
		if err := logClient.Receive(&msg); err != nil {
			break
		}
		dealReceiveSystemMsgS2C(msg)
		dealReceiveDBMsgS2C(msg)
	}
}

//处理接收到的系统消息
func dealReceiveSystemMsgS2C(msg []byte) {
	protoMsg := systemProto.UnmarshalProtoMsg(msg)
	if protoMsg == systemProto.NullProtoMsg {
		return
	}

	switch protoMsg.ID {
	case systemProto.ID_System_ConnectDBServerS2C:
		connectDBServerCallBack(protoMsg)
	}
}

//发送连接DB服务器
func ConnectDBServer() {
	INFO(global.ServerName + " Connect DBServer ...")
	send_msg := systemProto.MarshalProtoMsg(&systemProto.System_ConnectDBServerC2S{
		ServerName: protos.String(global.ServerName),
	})
	protos.Send(send_msg, logClient)
}

//连接DB服务器返回
func connectDBServerCallBack(protoMsg systemProto.ProtoMsg) {
	//	rev_msg := protoMsg.Body.(*systemProto.System_ConnectDBServerS2C)
	INFO(global.ServerName + " Connect DBServer Success")
}
