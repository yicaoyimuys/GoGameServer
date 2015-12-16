package gameProxy

import (
	"github.com/funny/link"
	"tools/dispatch"
	"module"
	"protos"
	"protos/gameProto"
)

var (
	MsgDispatch dispatch.DispatchInterface
)

func init() {
	MsgDispatch = dispatch.NewDispatch(
		dispatch.Handle{
			gameProto.ID_UserLoginC2S:				login,
			gameProto.ID_GetUserInfoC2S:			getUserInfo,
			gameProto.ID_AgainConnectC2S:			againConnect,
		},
	)
}

//登录
func login(session *link.Session, msg protos.ProtoMsg) {
	rev_msg := msg.Body.(*gameProto.UserLoginC2S)
	module.User.Login(rev_msg.GetUserName(), session)
}

//获取用户详细信息
func getUserInfo(session *link.Session, msg protos.ProtoMsg) {
	rev_msg := msg.Body.(*gameProto.GetUserInfoC2S)
	module.User.GetUserInfo(rev_msg.GetUserID(), session)
}

//重新连接
func againConnect(session *link.Session, msg protos.ProtoMsg) {
	rev_msg := msg.Body.(*gameProto.AgainConnectC2S)

	newSessionID := module.User.AgainConnect(rev_msg.GetSessionID(), session)
	if newSessionID != 0 {
		send_msg := protos.MarshalProtoMsg(&gameProto.AgainConnectS2C{
			SessionID: protos.Uint64(newSessionID),
		})
		protos.Send(session, send_msg)
	}
}
