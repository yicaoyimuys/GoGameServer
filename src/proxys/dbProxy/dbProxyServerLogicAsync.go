package dbProxy

import (
	"dao"
	"protos/dbProto"
	"protos"
	"github.com/funny/link"
)

/**
此文件处理接收到的异步的DB消息
*/

//更新用户最后登录时间
func updateUserLastLoginTime(session *link.Session, protoMsg protos.ProtoMsg) {
	rev_msg := protoMsg.Body.(*dbProto.DB_User_UpdateLastLoginTimeC2S)
	dao.UpdateUserLastLoginTime(rev_msg.GetUserID(), rev_msg.GetTime())
}
