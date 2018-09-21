package protos

import (
	. "core/libs"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"reflect"
)

type ProtoMsg struct {
	ID   uint16
	Body interface{}
}

var (
	NullProtoMsg = ProtoMsg{0, nil}
	MsgObjectMap = make(map[uint16]reflect.Type)
	MsgIDMap     = make(map[reflect.Type]uint16)
)

//设置消息类型和消息ID的对应关系
func SetMsg(msgID uint16, data interface{}) {
	msgType := reflect.TypeOf(data)

	MsgObjectMap[msgID] = msgType
	MsgIDMap[reflect.TypeOf(reflect.New(msgType).Interface())] = msgID
}

//根据消息ID获取消息实体
func GetMsgObject(msgID uint16) proto.Message {
	if msgType, exists := MsgObjectMap[msgID]; exists {
		return reflect.New(msgType).Interface().(proto.Message)
	} else {
		ERR("No MsgID:", msgID)
	}
	return nil
}

//根据一条消息获取消息ID
func GetMsgID(msg interface{}) uint16 {
	msgType := reflect.TypeOf(msg)
	if msgID, exists := MsgIDMap[msgType]; exists {
		return msgID
	} else {
		ERR("No MsgType:", msgType)
	}
	return 0
}

//序列化
func MarshalProtoMsg(args proto.Message) []byte {
	msgID := GetMsgID(args)
	msgBody, _ := proto.Marshal(args)

	result := make([]byte, 2+len(msgBody))
	binary.BigEndian.PutUint16(result[:2], msgID)
	copy(result[2:], msgBody)

	return result
}

//反序列化
func UnmarshalProtoMsg(msg []byte) ProtoMsg {
	if len(msg) < 2 {
		return NullProtoMsg
	}

	msgID := binary.BigEndian.Uint16(msg[:2])
	msgBody := GetMsgObject(msgID)
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

//封装消息String类型字段
func String(param string) *string {
	return proto.String(param)
}

//封装消息Uint64类型字段
func Uint64(param uint64) *uint64 {
	return proto.Uint64(param)
}

//封装消息Int64类型字段
func Int64(param int64) *int64 {
	return proto.Int64(param)
}

//封装消息Int32类型字段
func Int32(param int32) *int32 {
	return proto.Int32(param)
}

//封装消息Uint32类型字段
func Uint32(param uint32) *uint32 {
	return proto.Uint32(param)
}
