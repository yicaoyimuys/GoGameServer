package message

import (
	. "core/libs"
	"core/libs/grpc/ipc"
	"core/libs/sessions"
	"encoding/binary"
	"proto"
)

func IpcServerReceive(stream *ipc.Stream, msg *ipc.Req) {
	msgBody := msg.Data

	//获取Session
	id := msg.ServiceName + "_" + NumToString(msg.SessionId)
	session := sessions.GetBackSession(id)
	if session == nil {
		session = sessions.NewBackSession(msg.ServiceName, msg.SessionId, stream)
		session.SetMsgHandle(dealMessage)
		sessions.SetBackSession(session)
	}
	session.Receive(msgBody)
}

func dealMessage(session *sessions.BackSession, msgBody []byte) {
	//DEBUG(msgBody)
	//消息ID
	msgId := binary.BigEndian.Uint16(msgBody[:2])
	//DEBUG("BackMessage收到消息ID：", msgId)

	//消息解析
	msgData := proto.DecodeMsg(msgId, msgBody)
	if msgData == nil {
		ERR("收到错误消息ID: " + NumToString(msgId))
		session.Close()
		return
	}

	//消息处理
	handle := GetIpcServerHandle(msgId)
	if handle == nil {
		ERR("收到未处理的消息ID: " + NumToString(msgId))
		return
	}
	handle(session, msgData)
}
