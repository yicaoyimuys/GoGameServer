package messages

import (
	"core/libs/sessions"
	"github.com/golang/protobuf/proto"
)

type ipcServerMsgHandle func(clientSession *sessions.BackSession, msgData proto.Message)

var (
	backHandles = make(map[uint16]ipcServerMsgHandle)
)

func init() {
	//backHandles[msg.ID_Game_join_c2s] = module.RoomJoin
	//backHandles[msg.ID_Game_loadProgress_c2s] = module.RoomLoadProgress
	//backHandles[msg.ID_Game_event_c2s] = module.RoomEvent
	//backHandles[msg.ID_Game_result_c2s] = module.RoomResult
	//backHandles[msg.ID_Game_exit_c2s] = module.RoomExit
	//backHandles[msg.ID_Client_network_c2s] = module.ClientNetwork
	//
	//backHandles[msg.ID_Platform_login_c2s] = module.PlatformLogin
	//
	//backHandles[msg.ID_System_userOffline_c2s] = module.ClientOffline
}

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
