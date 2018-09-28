package public

import (
	"core"
	"core/libs/grpc/ipc"
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

func SendMsgToClientList(userSessionIds []uint64, sendMsg proto.Message) {
	data := protos.MarshalProtoMsg(sendMsg)
	streams := core.Service.GetIpcServerStreams()
	for _, stream := range streams {
		msg := &ipc.Res{
			UserSessionIds: userSessionIds,
			Data:           data,
		}
		stream.Send(msg)
	}
}

func SendMsgToAllClient(sendMsg proto.Message) {
	data := protos.MarshalProtoMsg(sendMsg)
	streams := core.Service.GetIpcServerStreams()
	for _, stream := range streams {
		msg := &ipc.Res{
			Data: data,
		}
		stream.Send(msg)
	}
}
