package module

import (
	"github.com/funny/link"
	"protos"
	. "protos/gameProto"
	. "model"
)


//发送通用错误消息
func SendErrorMsg(session *link.Session, errID int32) {
	send_msg := MarshalProtoMsg(&ErrorMsgS2C{
		ErrorID: protos.Int32(errID),
	})
	protos.Send(session, send_msg)
}

//发送其他客户端登录
func SendOtherLogin(session *link.Session) {
	protos.Send(session, MarshalProtoMsg(&OtherLoginS2C{}))
}

//发送服务器连接成功
func SendConnectSuccess(session *link.Session) {
	protos.Send(session, MarshalProtoMsg(&ConnectSuccessS2C{}))
}

//发送登录结果
func SendLoginResult(session *link.Session, userID uint64) {
	send_msg := MarshalProtoMsg(&UserLoginS2C{
		UserID: protos.Uint64(userID),
	})
	protos.Send(session, send_msg)
}

//发送获取用户数据结果
func SendGetUserInfoResult(session *link.Session, errorCode int32, u *UserModel) {
	if errorCode != 0 {
		SendErrorMsg(session, errorCode)
	} else {
		send_msg := MarshalProtoMsg(&GetUserInfoS2C{
			UserInfo: &Person{
				ID:        protos.Uint64(u.DBUser.ID),
				Name:      protos.String(u.DBUser.Name),
				Money:     protos.Int32(u.DBUser.Money),
				SessionID: protos.Uint64(session.Id()),
			},
		})
		protos.Send(session, send_msg)
	}
}
