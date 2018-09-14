package message

import (
	"connectorServer/sessions"
	. "connectorServer/tools"
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
