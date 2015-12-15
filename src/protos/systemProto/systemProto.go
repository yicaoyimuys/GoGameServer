package systemProto

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/funny/binary"
	//	"strconv"
)
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

//序列化
func MarshalProtoMsg(args proto.Message) []byte {
	msgID := protos.GetMsgID(args)
	msgBody, _ := proto.Marshal(args)

	result := make([]byte, 2+len(msgBody))
	binary.PutUint16LE(result[:2], msgID)
	copy(result[2:], msgBody)

	return result
}

//反序列化
func UnmarshalProtoMsg(msg []byte) protos.ProtoMsg {
	if len(msg) < 2 {
		return protos.NullProtoMsg
	}

	msgID := binary.GetUint16LE(msg[:2])
	if !IsValidID(msgID) {
		return protos.NullProtoMsg
	}

	msgBody := protos.GetMsgObject(msgID)
	if msgBody == nil {
		return protos.NullProtoMsg
	}

	err := proto.Unmarshal(msg[2:], msgBody)
	if err != nil {
		return protos.NullProtoMsg
	}

	return protos.ProtoMsg{
		ID:   msgID,
		Body: msgBody,
	}
}
