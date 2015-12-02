package module

import (
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"protos"
	. "protos/gameProto"
	//	. "tools"
	. "model"
)

//接收消息处理
func ReceiveMessage(session *link.Session, msg packet.RAW) {
	protoMsg := UnmarshalProtoMsg(msg)
	if protoMsg == NullProtoMsg {
		SendErrorMsg(Message_Error, session)
		return
	}

	//	DEBUG("收到消息ID: " + strconv.Itoa(int(msgID)))

	switch protoMsg.ID {
	case ID_UserLoginC2S:
		login(protoMsg, session)
	case ID_GetUserInfoC2S:
		getUserInfo(protoMsg, session)
	case ID_AgainConnectC2S:
		againConnect(protoMsg, session)
	}
}

//发送通用错误消息
func SendErrorMsg(errID int32, session *link.Session) {
	send_msg := MarshalProtoMsg(&ErrorMsgS2C{
		ErrorID: protos.Int32(errID),
	})
	protos.Send(send_msg, session)
}

//其他客户端登录
func SendOtherLogin(session *link.Session) {
	protos.Send(MarshalProtoMsg(&OtherLoginS2C{}), session)
}

//服务器连接成功
func SendConnectSuccess(session *link.Session) {
	protos.Send(MarshalProtoMsg(&ConnectSuccessS2C{}), session)
}

//登录
func login(msg ProtoMsg, session *link.Session) {
	rev_msg := msg.Body.(*UserLoginC2S)
	User.Login(rev_msg.GetUserName(), session)
}

//发送登录结果
func SendLoginResult(userID uint64, session *link.Session) {
	send_msg := MarshalProtoMsg(&UserLoginS2C{
		UserID: protos.Uint64(userID),
	})
	protos.Send(send_msg, session)
}

//重新连接
func againConnect(msg ProtoMsg, session *link.Session) {
	rev_msg := msg.Body.(*AgainConnectC2S)

	newSessionID := User.AgainConnect(rev_msg.GetSessionID(), session)
	if newSessionID != 0 {
		send_msg := MarshalProtoMsg(&AgainConnectS2C{
			SessionID: protos.Uint64(newSessionID),
		})
		protos.Send(send_msg, session)
	}
}

//获取用户详细信息
func getUserInfo(msg ProtoMsg, session *link.Session) {
	rev_msg := msg.Body.(*GetUserInfoC2S)
	User.GetUserInfo(rev_msg.GetUserID(), session)
}

func SendGetUserInfoResult(errorCode int32, u *UserModel, session *link.Session) {
	if errorCode != 0 {
		SendErrorMsg(errorCode, session)
	} else {
		send_msg := MarshalProtoMsg(&GetUserInfoS2C{
			UserInfo: &Person{
				ID:        protos.Uint64(u.DBUser.ID),
				Name:      protos.String(u.DBUser.Name),
				Money:     protos.Int32(u.DBUser.Money),
				SessionID: protos.Uint64(session.Id()),
			},
		})
		protos.Send(send_msg, session)
	}
}
