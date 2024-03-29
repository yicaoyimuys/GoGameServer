package messages

import (
	"github.com/yicaoyimuys/GoGameServer/core/libs/sessions"

	"google.golang.org/protobuf/proto"
)

type ipcServerMsgHandle func(clientSession *sessions.BackSession, msgData proto.Message)

var (
	backHandles = make(map[uint16]ipcServerMsgHandle)
)

func RegisterIpcServerHandle(msgId uint16, handle ipcServerMsgHandle) {
	backHandles[msgId] = handle
}

func GetIpcServerHandle(msgId uint16) ipcServerMsgHandle {
	handle, ok := backHandles[msgId]
	if ok {
		return handle
	} else {
		return nil
	}
}
