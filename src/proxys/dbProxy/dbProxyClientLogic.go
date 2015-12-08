package dbProxy

import (
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"global"
	"module"
	"protos"
	"protos/dbProto"
)

//处理接收到纯DB的消息
func dealReceiveDBMsgS2C(msg packet.RAW) {
	protoMsg := dbProto.UnmarshalProtoMsg(msg)
	if protoMsg == dbProto.NullProtoMsg {
		return
	}

	var session *link.Session = global.GetSession(protoMsg.Identification)
	if session == nil {
		return
	}

	switch protoMsg.ID {
	case dbProto.ID_DB_User_LoginS2C:
		userLoginCallBack(session, protoMsg)
	}
}

//用户登录使用
func UserLogin(identification uint64, userName string) {
	msg := dbProto.MarshalProtoMsg(identification, &dbProto.DB_User_LoginC2S{
		Name: protos.String(userName),
	})
	sendDBMsgToServer(msg)
}

//用户登录返回
func userLoginCallBack(session *link.Session, protoMsg dbProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*dbProto.DB_User_LoginS2C)
	module.User.UserLoginHandle(session, rev_msg.GetName(), rev_msg.GetID())
}
