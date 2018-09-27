package module

import (
	"core/libs/sessions"
)

type ChatRpcServer struct {
}

type RpcClientOfflineReq struct {
	ServiceIdentify string
	UserSessionId   uint64
}

type RpcClientOfflineRes struct {
}

func (this *ChatRpcServer) ClientOffline(args *RpcClientOfflineReq, reply *RpcClientOfflineRes) error {
	id := sessions.CreateBackSessionId(args.ServiceIdentify, args.UserSessionId)
	session := sessions.GetBackSession(id)
	if session != nil {
		session.Close()
	}
	return nil
}
