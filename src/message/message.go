package message

import (
	. "core/libs"
	"core/libs/grpc/ipc"
	"encoding/binary"
	"encoding/json"
	"errors"
	"game/errCode"
	"global"
	"proto"
	"proto/msg"
	"sessions"
)

//func isConnectorMsg(msgId uint16) bool {
//	return msgId >= 1000 && msgId <= 1999
//}
//
//func isGameMsg(msgId uint16) bool {
//	return msgId >= 3000 && msgId <= 3999
//}
//
//func isMatchingMsg(msgId uint16) bool {
//	return msgId >= 4000 && msgId <= 4999
//}

func isSystemMsg(msgId uint16) bool {
	return msgId >= 1 && msgId <= 999
}

func isConnectorMsg(msgId uint16) bool {
	return msgId == msg.ID_Client_ping_c2s
}

func isGameMsg(msgId uint16) bool {
	return true
}

func isMatchingMsg(msgId uint16) bool {
	ids := []uint16{
		msg.ID_Game_matching_c2s,
		msg.ID_Game_cancelMatching_c2s,

		msg.ID_Game_createReadyRoom_c2s,
		msg.ID_Game_joinReadyRoom_c2s,
		msg.ID_Game_leaveReadyRoom_c2s,
		msg.ID_Game_dissolveReadyRoom_c2s,
		msg.ID_Game_startByReadyRoom_c2s,
		msg.ID_Game_refuseReadyRoom_c2s,
		msg.ID_Game_againJoinReadyRoom_c2s,
	}
	return InArray(ids, msgId)
}

func getGameService(session *sessions.FrontSession, msgBody []byte) string {
	msgId := binary.BigEndian.Uint16(msgBody[:2])
	if msgId == msg.ID_Platform_login_c2s {
		//平台登录，根据roomId分配
		msgData := proto.DecodeMsg(msgId, msgBody)
		if msgData != nil {
			data := msgData.(*msg.Platform_login_c2s)

			var platformDataRequest map[string]interface{}
			json.Unmarshal([]byte(data.PlatformData), &platformDataRequest)
			roomId := platformDataRequest["roomId"].(string)
			return global.IpcClients.Game.GetServiceByFlag(roomId)
		} else {
			return ""
		}
	} else {
		//其他
		return session.IpcService()
	}
}

func getMatchingService(session *sessions.FrontSession, msgBody []byte) string {
	msgId := binary.BigEndian.Uint16(msgBody[:2])
	if msgId == msg.ID_Game_matching_c2s {
		//匹配，根据gameId分配
		msgData := proto.DecodeMsg(msgId, msgBody)
		if msgData != nil {
			data := msgData.(*msg.Game_matching_c2s)
			return global.IpcClients.Matching.GetServiceByFlag(NumToString(data.GameId))
		} else {
			return ""
		}
	} else if msgId == msg.ID_Game_createReadyRoom_c2s {
		//创建准备房间，根据roomId分配
		msgData := proto.DecodeMsg(msgId, msgBody)
		if msgData != nil {
			data := msgData.(*msg.Game_createReadyRoom_c2s)
			roomId := createReadyRoomId(data.GameId, data.UserId)
			return global.IpcClients.Matching.GetServiceByFlag(roomId)
		} else {
			return ""
		}
	} else if msgId == msg.ID_Game_joinReadyRoom_c2s {
		//加入准备房间，根据roomId分配
		msgData := proto.DecodeMsg(msgId, msgBody)
		if msgData != nil {
			data := msgData.(*msg.Game_joinReadyRoom_c2s)
			return global.IpcClients.Matching.GetServiceByFlag(data.RoomId)
		} else {
			return ""
		}
	} else if msgId == msg.ID_Game_againJoinReadyRoom_c2s {
		//再来一局，根据roomId分配
		msgData := proto.DecodeMsg(msgId, msgBody)
		if msgData != nil {
			data := msgData.(*msg.Game_againJoinReadyRoom_c2s)
			return global.IpcClients.Matching.GetServiceByFlag(data.RoomId)
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
	var ipcClient *ipc.Client
	var service string
	if serviceName == global.Services.Game {
		ipcClient = global.IpcClients.Game
		service = getGameService(session, msgBody)
	} else if serviceName == global.Services.Matching {
		ipcClient = global.IpcClients.Matching
		service = getMatchingService(session, msgBody)
	}

	if ipcClient == nil {
		return errors.New("serverName not exists")
	}

	if service == "" {
		return errors.New("service not exists")
	}

	err := ipcClient.Send(global.ServerName, session.ID(), msgBody, service)
	if err == nil {
		session.SetIpcService(serviceName, service)
	}
	return err
}

func SendMsgToBack_UserOffline(session *sessions.FrontSession) {
	sendMsg := msg.NewSystem_userOffline_c2s()
	sendMsgToBack(session.IpcServiceName(), session, sendMsg.Encode())
}
