package simple

import "code.google.com/p/goprotobuf/proto"
import (
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"module"
)

//封装发送消息
func send(msgID uint16, send_msg proto.Message, session *link.Session) {
	msgBody, _ := proto.Marshal(send_msg)

	msg := make([]byte, 2+len(msgBody))
	binary.PutUint16LE(msg, msgID)
	copy(msg[2:], msgBody)

	session.Send(packet.RAW(msg))
}

//发送通用错误消息
func SendErrorMsg(errID int32, session *link.Session) {
	send_msg := &ErrorMsg{
		ErrorID: proto.Int32(errID),
	}
	send(ID_ErrorMsg, send_msg, session)
}

//发送GameServer主动ping GateServer
func SendGamePingGateway(game_ip string, game_port string, session *link.Session) {
	send_msg := &GamePingGateway{
		IP:   proto.String(game_ip),
		Port: proto.String(game_port),
	}
	send(ID_GamePingGateway, send_msg, session)
}

//接收GamePingGateway消息
func ReceiveGamePingGateway(msgBody []byte) (string, string) {
	rev_msg := &GamePingGateway{}
	proto.Unmarshal(msgBody, rev_msg)
	return rev_msg.GetIP(), rev_msg.GetPort()
}

//其他客户端登录
func SendOtherLogin(session *link.Session) {
	send(ID_OtherLoginS2C, &OtherLoginS2C{}, session)
}

//服务器连接成功
func SendConnectSuccess(session *link.Session) {
	send(ID_ConnectSuccessS2C, &ConnectSuccessS2C{}, session)
}

//登录
func Login(msgBody []byte, session *link.Session) {
	rev_msg := &UserLoginC2S{}
	if err := proto.Unmarshal(msgBody, rev_msg); err != nil {
		return
	}

	userID := module.User.Login(rev_msg.GetUserName(), session)
	send_msg := &UserLoginS2C{
		UserID: proto.Int32(userID),
	}
	send(ID_UserLoginS2C, send_msg, session)
}

//重新连接
func AgainConnect(msgBody []byte, session *link.Session) {
	rev_msg := &AgainConnectC2S{}
	if err := proto.Unmarshal(msgBody, rev_msg); err != nil {
		return
	}

	newSessionID := module.User.AgainConnect(rev_msg.GetSessionID(), session)
	if newSessionID != 0 {
		send_msg := &AgainConnectS2C{
			SessionID: proto.Uint64(newSessionID),
		}
		send(ID_AgainConnectS2C, send_msg, session)
	}
}

//获取用户详细信息
func GetUserInfo(msgBody []byte, session *link.Session) {
	rev_msg := &GetUserInfoC2S{}
	if err := proto.Unmarshal(msgBody, rev_msg); err != nil {
		return
	}

	u := module.User.GetUserInfo(rev_msg.GetUserID(), session)
	if u != nil {
		send_msg := &GetUserInfoS2C{
			UserInfo: &Person{
				ID:        proto.Int32(u.DBUser.ID),
				Name:      proto.String(u.DBUser.Name),
				Money:     proto.Int32(u.DBUser.Money),
				SessionID: proto.Uint64(u.Session.Id()),
			},
		}

		send(ID_GetUserInfoS2C, send_msg, session)
	} else {
		SendErrorMsg(USER_NOT_EXISTS, session)
	}
}
