package messages

import (
	. "GoGameServer/core/libs"
	"GoGameServer/core/libs/sessions"
	"GoGameServer/core/protos"
)

func FontReceive(session *sessions.FrontSession, msgBody []byte) {
	//消息ID
	msgId := protos.UnmarshalProtoId(msgBody)
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
	} else if isChatMsg(msgId) {
		//聊天服务器消息
		dealChatMsg(session, msgBody)
	} else {
		ERR("WHAT???", msgId)
	}
}

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

func isChatMsg(msgId uint16) bool {
	return msgId >= 4000 && msgId <= 4999
}
