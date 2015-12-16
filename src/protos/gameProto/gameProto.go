package gameProto

import (
	"protos"
)

//初始化消息ID和消息类型的对应关系
func init() {
	protos.SetMsg(ID_ConnectSuccessS2C, ConnectSuccessS2C{})
	protos.SetMsg(ID_AgainConnectC2S, AgainConnectC2S{})
	protos.SetMsg(ID_AgainConnectS2C, AgainConnectS2C{})

	protos.SetMsg(ID_OtherLoginS2C, OtherLoginS2C{})
	protos.SetMsg(ID_ErrorMsgS2C, ErrorMsgS2C{})
	protos.SetMsg(ID_UserLoginC2S, UserLoginC2S{})
	protos.SetMsg(ID_UserLoginS2C, UserLoginS2C{})
	protos.SetMsg(ID_GetUserInfoC2S, GetUserInfoC2S{})
	protos.SetMsg(ID_GetUserInfoS2C, GetUserInfoS2C{})
}

//是否是有效的消息ID
func IsValidID(msgID uint16) bool {
	return msgID >= 1000 && msgID <= 9999
}

//是否是有效的登录消息
func IsValidLoginID(msgID uint16) bool {
	return msgID >= 1000 && msgID <= 1999
}

//是否是有效的WorldServer消息
func IsValidWorldID(msgID uint16) bool {
	return msgID >= 2000 && msgID <= 5999
}

//是否是有效的GameServer消息
func IsValidGameID(msgID uint16) bool {
	return msgID >= 6000 && msgID <= 9999
}
