package systemProto

import (
	"protos"
)

//初始化消息ID和消息类型的对应关系
func init() {
	protos.SetMsg(ID_System_ConnectDBServerC2S, System_ConnectDBServerC2S{})
	protos.SetMsg(ID_System_ConnectDBServerS2C, System_ConnectDBServerS2C{})
	protos.SetMsg(ID_System_ConnectTransferServerC2S, System_ConnectTransferServerC2S{})
	protos.SetMsg(ID_System_ConnectTransferServerS2C, System_ConnectTransferServerS2C{})
	protos.SetMsg(ID_System_ConnectWorldServerC2S, System_ConnectWorldServerC2S{})
	protos.SetMsg(ID_System_ConnectWorldServerS2C, System_ConnectWorldServerS2C{})
	protos.SetMsg(ID_System_ConnectLogServerC2S, System_ConnectLogServerC2S{})
	protos.SetMsg(ID_System_ConnectLogServerS2C, System_ConnectLogServerS2C{})
	protos.SetMsg(ID_System_ClientSessionOnlineC2S, System_ClientSessionOnlineC2S{})
	protos.SetMsg(ID_System_ClientSessionOfflineC2S, System_ClientSessionOfflineC2S{})
	protos.SetMsg(ID_System_ClientLoginSuccessC2S, System_ClientLoginSuccessC2S{})
}

//是否是有效的消息ID
func IsValidID(msgID uint16) bool {
	return msgID >= 10000 && msgID <= 10999
}
