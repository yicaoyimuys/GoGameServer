package service

import (
	"github.com/yicaoyimuys/GoGameServer/core"
	"github.com/yicaoyimuys/GoGameServer/core/libs/sessions"
	"github.com/yicaoyimuys/GoGameServer/core/libs/stack"
)

type ClientOffline struct {
}

type ClientOfflineReq struct {
	ServiceIdentify string
	UserSessionId   uint64
}

type ClientOfflineRes struct {
}

func (this *ClientOffline) Do(args *ClientOfflineReq, reply *ClientOfflineRes) error {
	id := sessions.CreateBackSessionId(args.ServiceIdentify, args.UserSessionId)
	session := sessions.GetBackSession(id)
	if session != nil {
		session.Close()
	}
	return nil
}

func (this *Service) frontSessionCreateHandle(session *sessions.FrontSession) {
	session.AddCloseCallback(nil, "FrontSessionOffline", func() {
		this.frontSessionOfflineHandle(session)
	})
}

func (this *Service) frontSessionOfflineHandle(session *sessions.FrontSession) {
	method := "ClientOffline.Do"
	args := &ClientOfflineReq{
		ServiceIdentify: core.Service.Identify(),
		UserSessionId:   session.ID(),
	}
	reply := &ClientOfflineRes{}

	//通知后端服务器
	for _, v := range this.rpcClients {
		var client = v
		go func() {
			defer stack.TryError()

			client.CallAll(method, args, reply)
		}()
	}
}
