package dbProxy

import (
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"global"
	"protos"
	"protos/systemProto"
	. "tools"
)

var (
	session *link.Session
)

//初始化
func InitClient(ip string, port string) error {
	//连接DB服务器
	addr := ip + ":" + port
	client, err := link.Connect("tcp", addr, packet.New(binary.SplitByUint32BE, 1024, 1024, 1024))
	if err != nil {
		return err
	}

	session = client
	go dealReceiveMsgS2C()
	ConnectDBServer()

	return nil
}

//发送DB消息到服务器
func sendDBMsgToServer(msg []byte) {
	if session == nil {
		dealReceiveDBMsgC2S(session, msg)
		dealReceiveAsyncDBMsgC2S(msg)
	} else {
		protos.Send(msg, session)
	}
}

//处理从DBServer发回的消息
func dealReceiveMsgS2C() {
	var msg packet.RAW
	for {
		if err := session.Receive(&msg); err != nil {
			break
		}
		dealReceiveSystemMsgS2C(msg)
		dealReceiveDBMsgS2C(msg)
	}
}

//处理接收到的系统消息
func dealReceiveSystemMsgS2C(msg packet.RAW) {
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
	protos.Send(send_msg, session)
}

//连接DB服务器返回
func connectDBServerCallBack(protoMsg systemProto.ProtoMsg) {
	//	rev_msg := protoMsg.Body.(*systemProto.System_ConnectDBServerS2C)
	INFO(global.ServerName + " Connect DBServer Success")
}
