package worldProxy

import (
	"github.com/funny/link"
	"github.com/funny/binary"
	"global"
	"protos"
	"protos/gameProto"
	"protos/systemProto"
	. "tools"
	"tools/dispatch"
)

var (
	worldClient 			 *link.Session
	clientMsgDispatch 	     dispatch.DispatchInterface
)

func init() {
	handle := dispatch.NewHandleConditions()
	//系统消息处理
	handle.Add(dispatch.HandleCondition{
		Condition: systemProto.IsValidID,
		H: dispatch.Handle{
			systemProto.ID_System_ConnectWorldServerS2C:		connectWorldServerCallBack,
		},
	})
	//游戏消息处理
	handle.Add(dispatch.HandleFuncCondition{
		Condition: gameProto.IsValidID,
		H: func(session *link.Session, msg []byte) {
			//发送到用户客户端
			msgIdentification := binary.GetUint64LE(msg[2:10])
			userSession := global.GetSession(msgIdentification)
			if userSession == nil {
				return
			}
			protos.Send(userSession, msg)
		},
	})

	//创建消息分派
	clientMsgDispatch = dispatch.NewDispatch(handle)
}

//初始化
func InitClient(ip string, port string) error {
	//连接WorldServer
	addr := ip + ":" + port
	client, err := global.Connect("WorldServer", "tcp", addr, global.PackCodecType_Safe, clientMsgDispatch)
	if err != nil {
		return err
	}

	worldClient = client
	sendConnectWorldServer()

	return nil
}

//发送系统消息到WorldServer
func SendSystemMsgToServer(msg []byte) {
	if worldClient == nil {
		return
	}
	protos.Send(worldClient, msg)
}

//发送游戏消息到WorldServer
func SendGameMsgToServer(msg []byte) {
	if worldClient == nil {
		return
	}
	protos.Send(worldClient, msg)
}

//发送连接WorldServer
func sendConnectWorldServer() {
	INFO(global.ServerName + " Connect WorldServer ...")
	send_msg := protos.MarshalProtoMsg(&systemProto.System_ConnectWorldServerC2S{
		ServerName: protos.String(global.ServerName),
		ServerID:   protos.Uint32(global.ServerID),
	})
	SendSystemMsgToServer(send_msg)
}

//连接Transfer服务器返回
func connectWorldServerCallBack(session *link.Session, protoMsg protos.ProtoMsg) {
	//	rev_msg := protoMsg.Body.(*systemProto.System_ConnectWorldServerS2C)
	INFO(global.ServerName + " Connect WorldServer Success")
}
