package dbProxy

import (
	"github.com/funny/link/packet"
	"dao"
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
		updateUserLastLoginTime(protoMsg)
	}
}

//更新用户最后登录时间
func updateUserLastLoginTime(protoMsg dbProto.ProtoMsg) error {
	rev_msg := protoMsg.Body.(*dbProto.DB_User_UpdateLastLoginTimeC2S)
	return dao.UpdateUserLastLoginTime(rev_msg.GetUserID(), rev_msg.GetTime())
}
