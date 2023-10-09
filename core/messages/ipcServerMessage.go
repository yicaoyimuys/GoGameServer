package messages

import (
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/grpc/ipc"
	"github.com/yicaoyimuys/GoGameServer/core/libs/protos"
	"github.com/yicaoyimuys/GoGameServer/core/libs/sessions"
	"go.uber.org/zap"
)

func IpcServerReceive(stream *ipc.Stream, msg *ipc.Req) {
	msgBody := msg.Data

	//获取Session
	id := sessions.CreateBackSessionId(msg.ServiceIdentify, msg.UserSessionId)
	session := sessions.GetBackSession(id)
	if session == nil {
		session = sessions.NewBackSession(id, msg.UserSessionId, stream)
		session.SetMsgHandle(dealMessage)
		sessions.SetBackSession(session)
	}
	session.Receive(msgBody)
}

func dealMessage(session *sessions.BackSession, msgBody []byte) {
	//消息解析
	protoMsg := protos.UnmarshalProtoMsg(msgBody)
	if protoMsg == protos.NullProtoMsg {
		msgId := protos.UnmarshalProtoId(msgBody)
		ERR("收到错误消息ID", zap.Uint16("MsgId", msgId))
		session.Close()
		return
	}

	//消息处理
	msgId := protoMsg.ID
	msgData := protoMsg.Body
	handle := GetIpcServerHandle(msgId)
	if handle == nil {
		ERR("收到未处理的消息ID", zap.Uint16("MsgId", msgId))
		return
	}
	handle(session, msgData)
}
