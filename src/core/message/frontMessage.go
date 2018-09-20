package message

import (
	. "core/libs"
	"core/libs/array"
	"core/libs/sessions"
	"encoding/binary"
	"proto/msg"
)

func FontReceive(session *sessions.FrontSession, msgBody []byte) {
	//DEBUG(msgBody)
	//消息ID
	msgId := binary.BigEndian.Uint16(msgBody[:2])
	//DEBUG("FrontMessage收到消息ID：", msgId)

	//消息处理
	if isSystemMsg(msgId) {
		//系统消息
		ERR("ERR???", msgId)
	} else if isConnectorMsg(msgId) {
		//连接服务器消息
		dealConnectorMsg(session, msgBody)
	} else if isMatchingMsg(msgId) {
		//匹配服务器消息
		dealMatchingMsg(session, msgBody)
	} else if isGameMsg(msgId) {
		//游戏服务器消息
		dealGameMsg(session, msgBody)
	} else {
		ERR("what???", msgId)
	}
}

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
	return array.InArray(ids, msgId)
}
