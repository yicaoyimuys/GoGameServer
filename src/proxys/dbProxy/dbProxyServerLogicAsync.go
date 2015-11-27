package dbProxy

import (
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"module_db"
	"protos/dbProto"
)

//处理接收到的异步的DB消息
func dealReceiveAsyncDBMsgC2S(msg packet.RAW) {
	protoMsg := dbProto.UnmarshalProtoMsg(msg)
	if protoMsg == dbProto.NullProtoMsg {
		return
	}

	switch protoMsg.ID {
	case dbProto.ID_DB_User_UpdateLastLoginTimeC2S:
		updateUserLastLoginTime(session, protoMsg)
	}
}

//更新用户最后登录时间
func updateUserLastLoginTime(session *link.Session, protoMsg dbProto.ProtoMsg) error {
	rev_msg := protoMsg.Body.(*dbProto.DB_User_UpdateLastLoginTimeC2S)
	return module_db.UpdateUserLastLoginTime(rev_msg.GetUserID(), rev_msg.GetTime())
}
