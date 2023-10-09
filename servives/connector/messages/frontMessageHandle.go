package messages

import (
	"errors"

	"github.com/yicaoyimuys/GoGameServer/core"
	"github.com/yicaoyimuys/GoGameServer/core/consts"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/grpc/ipc"
	"github.com/yicaoyimuys/GoGameServer/core/libs/protos"
	"github.com/yicaoyimuys/GoGameServer/core/libs/sessions"
	"github.com/yicaoyimuys/GoGameServer/servives/public/gameProto"
	"go.uber.org/zap"
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
	err := sendMsgToIpcService(consts.Service_Game, session, msgBody)
	if err != nil {
		ERR("DealGameMsg", zap.Error(err))
		sendErrorMsgToClient(session)
	}
}

func dealChatMsg(session *sessions.FrontSession, msgBody []byte) {
	err := sendMsgToIpcService(consts.Service_Chat, session, msgBody)
	if err != nil {
		ERR("DealChatMsg", zap.Error(err))
		sendErrorMsgToClient(session)
	}
}

func dealLoginMsg(session *sessions.FrontSession, msgBody []byte) {
	err := sendMsgToIpcService(consts.Service_Login, session, msgBody)
	if err != nil {
		ERR("DealLoginMsg", zap.Error(err))
		sendErrorMsgToClient(session)
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
		return session.GetIpcService(consts.Service_Game)
	}
}

func getChatService(session *sessions.FrontSession, msgBody []byte, ipcClient *ipc.Client) string {
	//1: 加入聊天，根据Token分配
	msgId := protos.UnmarshalProtoId(msgBody)
	if msgId == gameProto.ID_user_joinChat_c2s {
		protoMsg := protos.UnmarshalProtoMsg(msgBody)
		if protoMsg == protos.NullProtoMsg {
			return ""
		}
		protoMsgData := protoMsg.Body.(*gameProto.UserJoinChatC2S)
		return ipcClient.GetServiceByFlag(protoMsgData.GetToken())
	} else {
		return session.GetIpcService(consts.Service_Chat)
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
		return session.GetIpcService(consts.Service_Login)
	}
}

func sendErrorMsgToClient(session *sessions.FrontSession) {
	sendMsg := protos.MarshalProtoMsg(&gameProto.ErrorNoticeS2C{
		ErrorCode: protos.Int32(consts.ErrCode_SystemError),
	})
	session.Send(sendMsg)
}

func sendMsgToIpcService(serviceName string, clientSession *sessions.FrontSession, msgBody []byte) error {
	ipcClient := core.Service.GetIpcClient(serviceName)
	if ipcClient == nil {
		return errors.New(serviceName + ": ipcClient not exists")
	}

	var service string
	if serviceName == consts.Service_Login {
		service = getLoginService(clientSession, msgBody, ipcClient)
	} else if serviceName == consts.Service_Game {
		service = getGameService(clientSession, msgBody, ipcClient)
	} else if serviceName == consts.Service_Chat {
		service = getChatService(clientSession, msgBody, ipcClient)
	}

	if service == "" {
		return errors.New(serviceName + ": service not exists")
	}

	err := ipcClient.Send(core.Service.Identify(), clientSession.ID(), msgBody, service)
	if err == nil {
		clientSession.SetIpcService(serviceName, service)
	}
	return err
}
