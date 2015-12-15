package module

import (
	"github.com/funny/link"
	"tools/dispatch"
	"protos"
	. "protos/gameProto"
)

var (
	MsgDispatch dispatch.DispatchInterface
)

func init() {
	MsgDispatch = dispatch.NewDispatch(
		dispatch.Handle{
			ID_UserLoginC2S:			login,
			ID_GetUserInfoC2S:			getUserInfo,
			ID_AgainConnectC2S:			againConnect,
		},
	)
}

//接收消息处理
func ReceiveMessage(session *link.Session, msg []byte) {
	MsgDispatch.Process(session, msg)
}


//登录
func login(session *link.Session, msg protos.ProtoMsg) {
	rev_msg := msg.Body.(*UserLoginC2S)
	User.Login(rev_msg.GetUserName(), session)
}

//获取用户详细信息
func getUserInfo(session *link.Session, msg protos.ProtoMsg) {
	rev_msg := msg.Body.(*GetUserInfoC2S)
	User.GetUserInfo(rev_msg.GetUserID(), session)
}

//重新连接
func againConnect(session *link.Session, msg protos.ProtoMsg) {
	rev_msg := msg.Body.(*AgainConnectC2S)

	newSessionID := User.AgainConnect(rev_msg.GetSessionID(), session)
	if newSessionID != 0 {
		send_msg := MarshalProtoMsg(&AgainConnectS2C{
			SessionID: protos.Uint64(newSessionID),
		})
		protos.Send(session, send_msg)
	}
}
