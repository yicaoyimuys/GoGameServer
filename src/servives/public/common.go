package public

import (
	"core/libs/sessions"
	"core/protos"
	"core/protos/gameProto"
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
