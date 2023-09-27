package public

import (
	"github.com/yicaoyimuys/GoGameServer/core"
	"github.com/yicaoyimuys/GoGameServer/core/libs/protos"
	"github.com/yicaoyimuys/GoGameServer/core/libs/sessions"
	"github.com/yicaoyimuys/GoGameServer/servives/public/gameProto"

	"github.com/golang/protobuf/proto"
)

func SendErrorMsgToClient(session *sessions.BackSession, errorCode int32) {
	sendMsg := &gameProto.ErrorNoticeS2C{
		ErrorCode: protos.Int32(errorCode),
	}
	SendMsgToClient(session, sendMsg)
}

func SendMsgToClient(session *sessions.BackSession, sendMsg proto.Message) {
	if session == nil || sendMsg == nil {
		return
	}
	session.Send(protos.MarshalProtoMsg(sendMsg))
}

func SendMsgToClientList(userSessionIds []uint64, sendMsg proto.Message) {
	data := protos.MarshalProtoMsg(sendMsg)
	core.Service.GetIpcServer().SendToAllClient(userSessionIds, data)
}

func SendMsgToAllClient(sendMsg proto.Message) {
	data := protos.MarshalProtoMsg(sendMsg)
	core.Service.GetIpcServer().SendToAllClient(nil, data)
}
