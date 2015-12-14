package logProto

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/funny/binary"
)
import (
	"protos"
)

type ProtoMsg struct {
	ID   uint16
	Body interface{}
}

var (
	NullProtoMsg ProtoMsg = ProtoMsg{0, nil}
)

//初始化消息ID和消息类型的对应关系
func init() {
	protos.SetMsg(ID_Log_CommonLogC2S, Log_CommonLogC2S{})
}

//是否是有效的消息ID
func IsValidID(msgID uint16) bool {
	return msgID >= 15000 && msgID <= 15999
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
func UnmarshalProtoMsg(msg []byte) ProtoMsg {
	if len(msg) < 2 {
		return NullProtoMsg
	}

	msgID := binary.GetUint16LE(msg[:2])
	if !IsValidID(msgID) {
		return NullProtoMsg
	}

	msgBody := protos.GetMsgObject(msgID)
	if msgBody == nil {
		return NullProtoMsg
	}

	err := proto.Unmarshal(msg[2:], msgBody)
	if err != nil {
		return NullProtoMsg
	}

	return ProtoMsg{
		ID:   msgID,
		Body: msgBody,
	}
}
