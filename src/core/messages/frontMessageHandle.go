package messages

import (
	"core"
	"core/consts/errCode"
	"core/consts/service"
	. "core/libs"
	"core/libs/grpc/ipc"
	"core/libs/sessions"
	"core/protos"
	"core/protos/gameProto"
	"errors"
)

func dealConnectorMsg(session *sessions.FrontSession, msgBody []byte) {
	protoMsg := protos.UnmarshalProtoMsg(msgBody)
	if protoMsg == protos.NullProtoMsg {
		return
	}

	//Ping消息特殊处理
	if protoMsg.ID == gameProto.ID_client_ping_c2s {
		session.UpdatePingTime()
		return
	}
}

func dealGameMsg(session *sessions.FrontSession, msgBody []byte) {
	err := sendMsgToBack(Service.Game, session, msgBody)
	if err != nil {
		ERR("dealGameMsg", err)
		sendMsgToClient_Error(session)
	}
}

func dealLoginMsg(session *sessions.FrontSession, msgBody []byte) {
	err := sendMsgToBack(Service.Login, session, msgBody)
	if err != nil {
		ERR("dealLoginMsg", err)
		sendMsgToClient_Error(session)
	}
}

func getGameService(session *sessions.FrontSession, msgBody []byte, ipcClient *ipc.Client) string {
	//1: 获取用户数据，根据Token分配
	msgId := protos.UnmarshalProtoId(msgBody)
	if msgId == gameProto.ID_user_getInfo_c2s {
		protoMsg := protos.UnmarshalProtoMsg(msgBody)
		if protoMsg == protos.NullProtoMsg {
			return ""
		}
		protoMsgData := protoMsg.Body.(*gameProto.UserGetInfoC2S)
		return ipcClient.GetServiceByFlag(protoMsgData.GetToken())
	} else {
		return session.IpcService()
	}
}

func getLoginService(session *sessions.FrontSession, msgBody []byte, ipcClient *ipc.Client) string {
	//1: 登录，根据Account分配
	msgId := protos.UnmarshalProtoId(msgBody)
	if msgId == gameProto.ID_user_login_c2s {
		protoMsg := protos.UnmarshalProtoMsg(msgBody)
		if protoMsg == protos.NullProtoMsg {
			return ""
		}
		protoMsgData := protoMsg.Body.(*gameProto.UserLoginC2S)
		return ipcClient.GetServiceByFlag(protoMsgData.GetAccount())
	} else {
		return session.IpcService()
	}
}

func sendMsgToClient_Error(session *sessions.FrontSession) {
	sendMsg := protos.MarshalProtoMsg(&gameProto.ErrorNoticeS2C{
		ErrorCode: protos.Int32(ErrCode.ERR_SYSTEM),
	})
	session.Send(sendMsg)
}

func sendMsgToBack(serviceName string, session *sessions.FrontSession, msgBody []byte) error {
	ipcClient := core.Service.GetIpcClient(serviceName)
	if ipcClient == nil {
		ERR("ipcClient not exists", serviceName)
	}

	var service string
	if serviceName == Service.Login {
		service = getLoginService(session, msgBody, ipcClient)
	} else if serviceName == Service.Game {
		service = getGameService(session, msgBody, ipcClient)
	}

	if ipcClient == nil {
		return errors.New("serverName not exists")
	}

	if service == "" {
		return errors.New("service not exists")
	}

	err := ipcClient.Send(core.Service.Name(), core.Service.ID(), session.ID(), msgBody, service)
	if err == nil {
		session.SetIpcService(serviceName, service)
	}
	return err
}

func SendMsgToBack_UserOffline(session *sessions.FrontSession) {
	//sendMsg := msg.NewSystem_userOffline_c2s()
	//sendMsgToBack(session.IpcServiceName(), session, sendMsg.Encode())

	//TODO
}
