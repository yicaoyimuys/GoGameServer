package dbProto

import (
	"code.google.com/p/goprotobuf/proto"
)
import (
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"protos"
	//	. "tools"
)

type ProtoMsg struct {
	ID             uint16
	Body           interface{}
	Identification uint64
}

var NullProtoMsg ProtoMsg = ProtoMsg{0, nil, 0}

//初始化消息ID和消息类型的对应关系
func init() {
	protos.SetMsg(ID_DB_User_LoginC2S, DB_User_LoginC2S{})
	protos.SetMsg(ID_DB_User_LoginS2C, DB_User_LoginS2C{})

	protos.SetMsg(ID_DB_User_UpdateLastLoginTimeC2S, DB_User_UpdateLastLoginTimeC2S{})
}

//是否是有效的消息ID
func IsValidID(msgID uint16) bool {
	return msgID >= 11000 && msgID <= 14999
}

//是否是有效的异步DB消息
func IsValidAsyncID(msgID uint16) bool {
	return msgID >= 12000 && msgID <= 14999
}

//发送消息
func Send(msgBody []byte, session *link.Session) {
	session.Send(packet.RAW(msgBody))
}

//序列化
func MarshalProtoMsg(identification uint64, args proto.Message) []byte {
	msgID := protos.GetMsgID(args)

	msgBody, _ := proto.Marshal(args)

	result := make([]byte, 2+8+len(msgBody))
	binary.PutUint16LE(result[:2], msgID)
	binary.PutUint64LE(result[2:10], identification)
	copy(result[10:], msgBody)

	return result
}

//反序列化消息
func UnmarshalProtoMsg(msg []byte) ProtoMsg {
	if len(msg) < 10 {
		return NullProtoMsg
	}

	msgID := binary.GetUint16LE(msg[:2])
	if !IsValidID(msgID) {
		return NullProtoMsg
	}

	identification := binary.GetUint64LE(msg[2:10])

	msgBody := protos.GetMsgObject(msgID)
	if msgBody == nil {
		return NullProtoMsg
	}

	err := proto.Unmarshal(msg[10:], msgBody)
	if err != nil {
		return NullProtoMsg
	}

	return ProtoMsg{
		ID:             msgID,
		Body:           msgBody,
		Identification: identification,
	}
}
