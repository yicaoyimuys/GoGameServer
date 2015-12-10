package dbProxy

import (
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"dao"
	"protos"
	"protos/dbProto"
	"proxys/redisProxy"
//	."tools"
	"time"
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
	userName := rev_msg.GetName()

	//先从缓存中读取
	dbUser := redisProxy.GetDBUserByUserName(userName)
	if dbUser == nil {
		//从数据库中获取
		dbUser, _ = dao.GetUserByUserName(userName)
		//将数据缓存到Redis
		redisProxy.SetDBUser(dbUser)
	}

	//返回消息
	sendProtoMsg := &dbProto.DB_User_LoginS2C{}
	if dbUser != nil {
		sendProtoMsg.ID = protos.Uint64(dbUser.ID)
		sendProtoMsg.Name = protos.String(dbUser.Name)
	}
	send_msg := dbProto.MarshalProtoMsg(protoMsg.Identification, sendProtoMsg)
	sendDBMsgToClient(session, send_msg)

	//更新最后登录时间
	if dbUser != nil{
		dbUser.LastLoginTime = time.Now().Unix()
		redisProxy.UpdateUserLastLoginTime(dbUser)
	}
}
