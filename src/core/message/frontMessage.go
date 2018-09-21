package message

import (
	. "core/libs"
	"core/libs/sessions"
	"encoding/binary"
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
	} else if isLoginMsg(msgId) {
		//登录服务器消息
		dealLoginMsg(session, msgBody)
	} else if isGameMsg(msgId) {
		//游戏服务器消息
		dealGameMsg(session, msgBody)
	} else {
		ERR("what???", msgId)
	}
}

//1-999: 系统消息
//1000-1999: connector消息
//2000-2999: login消息
//3000-3999: game消息

func isSystemMsg(msgId uint16) bool {
	return msgId >= 1 && msgId <= 999
}

func isConnectorMsg(msgId uint16) bool {
	return msgId >= 1000 && msgId <= 1999
}

func isLoginMsg(msgId uint16) bool {
	return msgId >= 2000 && msgId <= 2999
}

func isGameMsg(msgId uint16) bool {
	return msgId >= 3000 && msgId <= 3999
}
