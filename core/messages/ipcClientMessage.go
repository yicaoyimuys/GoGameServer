package messages

import (
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/grpc/ipc"
	"github.com/yicaoyimuys/GoGameServer/core/libs/protos"
	"github.com/yicaoyimuys/GoGameServer/core/libs/sessions"
	"go.uber.org/zap"
)

func IpcClientReceive(stream ipc.Ipc_TransferClient, msg *ipc.Res) {
	if msg.UserSessionIds == nil {
		//发送给所有人
		sessions.FetchFrontSession(func(clientSession *sessions.FrontSession) {
			clientSession.Send(msg.Data)
		})
	} else {
		//发送给多个人
		for _, userSessionId := range msg.UserSessionIds {
			clientSession := sessions.GetFrontSession(userSessionId)
			if clientSession != nil {
				clientSession.Send(msg.Data)
			} else {
				msgId := protos.UnmarshalProtoId(msg.Data)
				WARN("FrontSession No Exists", zap.Uint16("MsgId", msgId))
			}
		}
	}
}
