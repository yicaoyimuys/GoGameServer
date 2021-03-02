package rpcModules

import (
	"GoGameServer/core/libs/sessions"
)

type Client struct {
}

type ClientOfflineReq struct {
	ServiceIdentify string
	UserSessionId   uint64
}

type ClientOfflineRes struct {
}

func (this *Client) Offline(args *ClientOfflineReq, reply *ClientOfflineRes) error {
	id := sessions.CreateBackSessionId(args.ServiceIdentify, args.UserSessionId)
	session := sessions.GetBackSession(id)
	if session != nil {
		session.Close()
	}
	return nil
}
