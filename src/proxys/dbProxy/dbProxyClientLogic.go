package dbProxy

import (
	"github.com/funny/link"
	"module"
	"protos"
	"protos/dbProto"
	"tools/dispatch"
)

var (
	ClientDbMsgDispatchHandle   dispatch.HandleInterface
)

func init()  {
	ClientDbMsgDispatchHandle = dispatch.Handle{
		dbProto.ID_DB_User_LoginS2C:					userLoginCallBack,
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
func userLoginCallBack(userSession *link.Session, protoMsg protos.ProtoMsg) {
	rev_msg := protoMsg.Body.(*dbProto.DB_User_LoginS2C)
	module.User.UserLoginHandle(userSession, rev_msg.GetName(), rev_msg.GetID())
}
