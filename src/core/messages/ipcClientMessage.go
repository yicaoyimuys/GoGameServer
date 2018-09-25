package messages

import (
	. "core/libs"
	"core/libs/grpc/ipc"
	"core/libs/sessions"
	"core/protos"
)

func IpcClientReceive(stream ipc.Ipc_TransferClient, msg *ipc.Res) {
	clientSession := sessions.GetFrontSession(msg.UserSessionId)
	if clientSession != nil {
		clientSession.Send(msg.Data)
	} else {
		msgId := protos.UnmarshalProtoId(msg.Data)
		WARN("frontSession no exists", msgId)
	}
}
