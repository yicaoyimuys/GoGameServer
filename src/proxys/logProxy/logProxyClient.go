package logProxy

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
	addr := ip + ":" + port
	client, err := link.Connect("tcp", addr, packet.New(binary.SplitByUint32BE, 1024, 1024, 1024))
	if err != nil {
		return err
	}

	session = client
	go dealReceiveMsg()
	ConnectLogServer()

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
	case systemProto.ID_System_ConnectLogServerS2C:
		connectLogServerCallBack(protoMsg)
	}
}

//处理接收到的消息
func dealReceiveMsgS2C(msg packet.RAW) {
	if len(msg) < 2 {
		return
	}

	msgID := binary.GetUint16LE(msg[:2])
	if systemProto.IsValidID(msgID) {
		dealReceiveSystemMsgS2C(msg)
	} else {
		ERR(global.ServerName, "收到未处理消息")
	}
}

//发送系统消息到LogServer
func SendSystemMsgToServer(msg []byte) {
	if session == nil {
		return
	}
	protos.Send(msg, session)
}

//发送Log消息到LogServer
func SendLogMsgToServer(msg []byte) {
	if session == nil {
		dealLogMsgC2S(msg)
		return
	}
	protos.Send(msg, session)
}

//发送连接LogServer
func ConnectLogServer() {
	INFO(global.ServerName + " Connect LogServer ...")
	send_msg := systemProto.MarshalProtoMsg(&systemProto.System_ConnectLogServerC2S{
		ServerName: protos.String(global.ServerName),
	})
	SendSystemMsgToServer(send_msg)
}

//连接Transfer服务器返回
func connectLogServerCallBack(protoMsg systemProto.ProtoMsg) {
	//	rev_msg := protoMsg.Body.(*systemProto.System_ConnectLogServerS2C)
	INFO(global.ServerName + " Connect LogServer Success")
}
