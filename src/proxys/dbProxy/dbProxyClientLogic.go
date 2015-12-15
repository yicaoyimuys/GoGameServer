package dbProxy

import (
	"github.com/funny/link"
	"global"
	"module"
	"protos"
	"protos/dbProto"
)

//用户登录使用
func UserLogin(identification uint64, userName string) {
	msg := dbProto.MarshalProtoMsg(identification, &dbProto.DB_User_LoginC2S{
		Name: protos.String(userName),
	})
	sendDBMsgToServer(msg)
}

//用户登录返回
func userLoginCallBack(session *link.Session, protoMsg protos.ProtoMsg) {
	var userSession *link.Session = global.GetSession(protoMsg.Identification)
	if userSession == nil {
		return
	}

	rev_msg := protoMsg.Body.(*dbProto.DB_User_LoginS2C)
	module.User.UserLoginHandle(userSession, rev_msg.GetName(), rev_msg.GetID())
}
