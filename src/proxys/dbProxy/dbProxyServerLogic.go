package dbProxy

import (
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"module_db"
	"protos"
	"protos/dbProto"
	"proxys/redisProxy"
)

//处理接收到的同步DB消息
func dealReceiveDBMsgC2S(session *link.Session, msg packet.RAW) {
	protoMsg := dbProto.UnmarshalProtoMsg(msg)
	if protoMsg == dbProto.NullProtoMsg {
		return
	}

	switch protoMsg.ID {
	case dbProto.ID_DB_User_LoginC2S:
		userLogin(session, protoMsg)
	}
}

//用户登录
func userLogin(session *link.Session, protoMsg dbProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*dbProto.DB_User_LoginC2S)

	sendProtoMsg := &dbProto.DB_User_LoginS2C{}

	//从数据库中获取
	dbUser, _ := module_db.GetUserByUserName(rev_msg.GetName())
	if dbUser != nil {
		//将数据缓存到Redis
		redisProxy.SetDBUser(dbUser)

		sendProtoMsg.ID = protos.Uint64(dbUser.ID)
		sendProtoMsg.Name = protos.String(dbUser.Name)
	}

	send_msg := dbProto.MarshalProtoMsg(protoMsg.Identification, sendProtoMsg)
	sendDBMsgToClient(session, send_msg)
}
