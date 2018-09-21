package messages

import (
	. "core/libs"
	"core/libs/grpc/ipc"
	"core/libs/sessions"
	"core/protos"
)

func IpcClientReceive(stream ipc.Ipc_TransferClient, msg *ipc.Res) {
	frontSession := sessions.GetFrontSession(msg.SessionId)
	msgBody := msg.Data
	if frontSession != nil {
		frontSession.Send(msgBody)
	} else {
		msgId := protos.UnmarshalProtoId(msgBody)
		WARN("frontSession no exists", msgId)
	}
}
