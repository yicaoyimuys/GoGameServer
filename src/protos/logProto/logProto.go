package logProto

import (
	"protos"
)

//初始化消息ID和消息类型的对应关系
func init() {
	protos.SetMsg(ID_Log_CommonLogC2S, Log_CommonLogC2S{})
}

//是否是有效的消息ID
func IsValidID(msgID uint16) bool {
	return msgID >= 15000 && msgID <= 15999
}
