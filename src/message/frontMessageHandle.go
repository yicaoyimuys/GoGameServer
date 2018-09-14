package message

import (
	. "core/libs"
	"encoding/binary"
	"global"
	"proto/msg"
	"sessions"
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
	err := sendMsgToBack(global.Services.Game, session, msgBody)
	if err != nil {
		ERR("dealGameMsg", err)
		sendMsgToClient_Error(session)
	}
}

func dealMatchingMsg(session *sessions.FrontSession, msgBody []byte) {
	err := sendMsgToBack(global.Services.Matching, session, msgBody)
	if err != nil {
		ERR("dealMatchingMsg", err)
		sendMsgToClient_Error(session)
	}
}
