package messages

import (
	. "core/libs"
	"core/libs/grpc/ipc"
	"core/libs/sessions"
	"core/protos"
)

func IpcClientReceive(stream ipc.Ipc_TransferClient, msg *ipc.Res) {
	if msg.UserSessionId == 0 {
		//发送给所有人
		sessions.FetchFrontSession(func(clientSession *sessions.FrontSession) {
			clientSession.Send(msg.Data)
		})
	} else {
		//发送给单个人
		clientSession := sessions.GetFrontSession(msg.UserSessionId)
		if clientSession != nil {
			clientSession.Send(msg.Data)
		} else {
			msgId := protos.UnmarshalProtoId(msg.Data)
			WARN("frontSession no exists", msgId)
		}
	}
}
