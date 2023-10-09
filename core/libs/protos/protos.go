package protos

import (
	"encoding/binary"
	"reflect"

	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"go.uber.org/zap"

	"google.golang.org/protobuf/proto"
)

type ProtoMsg struct {
	ID   uint16
	Body proto.Message
}

var (
	NullProtoMsg = ProtoMsg{0, nil}
	MsgObjectMap = make(map[uint16]reflect.Type)
	MsgIDMap     = make(map[reflect.Type]uint16)
)

// 设置消息类型和消息ID的对应关系
func SetMsg(msgId uint16, data interface{}) {
	msgType := reflect.TypeOf(data)

	MsgObjectMap[msgId] = msgType
	MsgIDMap[reflect.TypeOf(reflect.New(msgType).Interface())] = msgId
}

// 根据消息ID获取消息实体
func GetMsgObject(msgId uint16) proto.Message {
	if msgType, exists := MsgObjectMap[msgId]; exists {
		return reflect.New(msgType).Interface().(proto.Message)
	} else {
		ERR("No MsgId", zap.Uint16("MsgId", msgId))
	}
	return nil
}

// 根据一条消息获取消息ID
func GetMsgId(msg interface{}) uint16 {
	msgType := reflect.TypeOf(msg)
	if msgId, exists := MsgIDMap[msgType]; exists {
		return msgId
	} else {
		ERR("No MsgType", zap.Any("MsgType", msgType))
	}
	return 0
}

// 序列化
func MarshalProtoMsg(args proto.Message) []byte {
	msgId := GetMsgId(args)
	msgBody, _ := proto.Marshal(args)

	result := make([]byte, 2+len(msgBody))
	binary.BigEndian.PutUint16(result[:2], msgId)
	copy(result[2:], msgBody)

	return result
}

func UnmarshalProtoId(msg []byte) uint16 {
	return binary.BigEndian.Uint16(msg[:2])
}

// 反序列化
func UnmarshalProtoMsg(msg []byte) ProtoMsg {
	if len(msg) < 2 {
		return NullProtoMsg
	}

	msgId := UnmarshalProtoId(msg)
	msgBody := GetMsgObject(msgId)
	if msgBody == nil {
		return NullProtoMsg
	}

	err := proto.Unmarshal(msg[2:], msgBody)
	if err != nil {
		return NullProtoMsg
	}

	return ProtoMsg{
		ID:   msgId,
		Body: msgBody,
	}
}

func String(param string) *string {
	return proto.String(param)
}

func Bool(param bool) *bool {
	return proto.Bool(param)
}

func Float64(param float64) *float64 {
	return proto.Float64(param)
}

func Float32(param float32) *float32 {
	return proto.Float32(param)
}

func Int64(param int64) *int64 {
	return proto.Int64(param)
}

func Uint64(param uint64) *uint64 {
	return proto.Uint64(param)
}

func Int32(param int32) *int32 {
	return proto.Int32(param)
}

func Uint32(param uint32) *uint32 {
	return proto.Uint32(param)
}
