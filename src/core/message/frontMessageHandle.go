package message

import (
	"core"
	"core/consts/errCode"
	"core/consts/service"
	. "core/libs"
	"core/libs/grpc/ipc"
	"core/libs/sessions"
	"encoding/binary"
	"encoding/json"
	"errors"
	"proto"
	"proto/msg"
)

func dealConnectorMsg(session *sessions.FrontSession, msgBody []byte) {
	msgId := binary.BigEndian.Uint16(msgBody[:2])

	//Ping消息特殊处理
	if msgId == msg.ID_Client_ping_c2s {
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

func dealMatchingMsg(session *sessions.FrontSession, msgBody []byte) {
	err := sendMsgToBack(Service.Matching, session, msgBody)
	if err != nil {
		ERR("dealMatchingMsg", err)
		sendMsgToClient_Error(session)
	}
}

func getGameService(session *sessions.FrontSession, msgBody []byte, ipcClient *ipc.Client) string {
	msgId := binary.BigEndian.Uint16(msgBody[:2])
	if msgId == msg.ID_Platform_login_c2s {
		//平台登录，根据roomId分配
		msgData := proto.DecodeMsg(msgId, msgBody)
		if msgData != nil {
			data := msgData.(*msg.Platform_login_c2s)

			var platformDataRequest map[string]interface{}
			json.Unmarshal([]byte(data.PlatformData), &platformDataRequest)
			roomId := platformDataRequest["roomId"].(string)
			return ipcClient.GetServiceByFlag(roomId)
		} else {
			return ""
		}
	} else {
		//其他
		return session.IpcService()
	}
}

func getMatchingService(session *sessions.FrontSession, msgBody []byte, ipcClient *ipc.Client) string {
	msgId := binary.BigEndian.Uint16(msgBody[:2])
	if msgId == msg.ID_Game_matching_c2s {
		//匹配，根据gameId分配
		msgData := proto.DecodeMsg(msgId, msgBody)
		if msgData != nil {
			data := msgData.(*msg.Game_matching_c2s)
			return ipcClient.GetServiceByFlag(NumToString(data.GameId))
		} else {
			return ""
		}
	} else if msgId == msg.ID_Game_createReadyRoom_c2s {
		//创建准备房间，根据roomId分配
		msgData := proto.DecodeMsg(msgId, msgBody)
		if msgData != nil {
			data := msgData.(*msg.Game_createReadyRoom_c2s)
			roomId := createReadyRoomId(data.GameId, data.UserId)
			return ipcClient.GetServiceByFlag(roomId)
		} else {
			return ""
		}
	} else if msgId == msg.ID_Game_joinReadyRoom_c2s {
		//加入准备房间，根据roomId分配
		msgData := proto.DecodeMsg(msgId, msgBody)
		if msgData != nil {
			data := msgData.(*msg.Game_joinReadyRoom_c2s)
			return ipcClient.GetServiceByFlag(data.RoomId)
		} else {
			return ""
		}
	} else if msgId == msg.ID_Game_againJoinReadyRoom_c2s {
		//再来一局，根据roomId分配
		msgData := proto.DecodeMsg(msgId, msgBody)
		if msgData != nil {
			data := msgData.(*msg.Game_againJoinReadyRoom_c2s)
			return ipcClient.GetServiceByFlag(data.RoomId)
		} else {
			return ""
		}
	} else {
		//其他
		return session.IpcService()
	}
}

func createReadyRoomId(gameId uint16, createUserId string) string {
	return "readyRoom-" + NumToString(gameId) + "-" + createUserId
}

func sendMsgToClient_Error(session *sessions.FrontSession) {
	sendMsg := msg.NewError_notice_s2c()
	sendMsg.ErrorCode = ErrCode.ERR_SYSTEM
	session.Send(sendMsg.Encode())
}

func sendMsgToBack(serviceName string, session *sessions.FrontSession, msgBody []byte) error {
	ipcClient := core.Service.GetIpcClient(serviceName)
	if ipcClient == nil {
		ERR("ipcClient not exists", serviceName)
	}

	var service string
	if serviceName == Service.Game {
		service = getGameService(session, msgBody, ipcClient)
	} else if serviceName == Service.Matching {
		service = getMatchingService(session, msgBody, ipcClient)
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
	sendMsg := msg.NewSystem_userOffline_c2s()
	sendMsgToBack(session.IpcServiceName(), session, sendMsg.Encode())
}
