package message

import (
	"connectorServer/sessions"
	. "connectorServer/tools"
	"connectorServer/tools/grpc/ipc"
	"encoding/binary"
)

func BackReceive(stream ipc.Ipc_TransferClient, msg *ipc.Res) {
	frontSession := sessions.GetFrontSession(msg.SessionId)
	msgBody := msg.Data
	if frontSession != nil {
		frontSession.Send(msgBody)
	} else {
		msgId := binary.BigEndian.Uint16(msgBody[:2])
		WARN("frontSession no exists", msgId)
	}
}
