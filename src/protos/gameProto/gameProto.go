package gameProto

import "code.google.com/p/goprotobuf/proto"
import (
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"protos"
	//	. "tools"
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

//发送消息
func Send(msg []byte, session *link.Session) {
	session.Send(packet.RAW(msg))
}