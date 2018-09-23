package module

import (
	. "core/libs"
	"core/libs/sessions"
)

type LoginRpcServer struct {
}

type RpcClientOfflineReq struct {
	ServiceName string
	SessionId   uint64
}

type RpcClientOfflineRes struct {
}

func (this *LoginRpcServer) ClientOffline(args *RpcClientOfflineReq, reply *RpcClientOfflineRes) error {
	id := args.ServiceName + "_" + NumToString(args.SessionId)
	session := sessions.GetBackSession(id)
	if session != nil {
		session.Close()
	}
	return nil
}
