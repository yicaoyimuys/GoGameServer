package proto

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"reflect"
)

var defaultOrder = binary.BigEndian
var msgObjectMap_ByName map[string]reflect.Type = make(map[string]reflect.Type)
var msgObjectMap_ById map[uint16]reflect.Type = make(map[uint16]reflect.Type)

func SetMsgByName(msgName string, data interface{}) {
	msgType := reflect.TypeOf(data)
	msgObjectMap_ByName[msgName] = msgType
}

func SetMsgById(msgId uint16, data interface{}) {
	msgType := reflect.TypeOf(data)
	msgObjectMap_ById[msgId] = msgType
}

func DecodeMsg(msgId uint16, msgBody []byte) Msg {
	class := msgObjectMap_ById[msgId]
	if class == nil {
		return nil
	}
	entity := reflect.New(class).Interface().(Msg)
	entity.Decode(msgBody)
	return entity
}

func SetUint8(data *bytes.Buffer, num uint8) {
	binary.Write(data, defaultOrder, num)
}

func SetInt8(data *bytes.Buffer, num int8) {
	binary.Write(data, defaultOrder, num)
}

func SetUint16(data *bytes.Buffer, num uint16) {
	binary.Write(data, defaultOrder, num)
}

func SetInt16(data *bytes.Buffer, num int16) {
	binary.Write(data, defaultOrder, num)
}

func SetUint32(data *bytes.Buffer, num uint32) {
	binary.Write(data, defaultOrder, num)
}

func SetInt32(data *bytes.Buffer, num int32) {
	binary.Write(data, defaultOrder, num)
}

func SetUint64(data *bytes.Buffer, num uint64) {
	binary.Write(data, defaultOrder, num)
}

func SetInt64(data *bytes.Buffer, num int64) {
	//binary.Write(data, defaultOrder, num);

	var tmp float64 = float64(num)
	binary.Write(data, defaultOrder, tmp)
}

func SetFloat(data *bytes.Buffer, num float32) {
	binary.Write(data, defaultOrder, num)
}

func SetString(data *bytes.Buffer, str string) {
	strBytes := []byte(str)
	strLen := len(strBytes)

	SetUint8(data, 1)
	SetUint16(data, uint16(strLen))
	binary.Write(data, defaultOrder, strBytes)
}

func SetBuffer(data *bytes.Buffer, buff []byte) {
	buffLen := len(buff)

	SetUint8(data, 1)
	SetUint16(data, uint16(buffLen))
	binary.Write(data, defaultOrder, buff)
}

func SetEntity(data *bytes.Buffer, entity Msg) {
	SetBuffer(data, entity.Encode())
}

func SetArray(data *bytes.Buffer, list *list.List, arrType string) {
	listLen := uint16(list.Len())
	//DEBUG(listLen)

	SetUint16(data, listLen)
	for e := list.Front(); e != nil; e = e.Next() {
		if arrType == "uint8" {
			SetUint8(data, e.Value.(uint8))
		} else if arrType == "int8" {
			SetInt8(data, e.Value.(int8))
		} else if arrType == "uint16" {
			SetUint16(data, e.Value.(uint16))
		} else if arrType == "int16" {
			SetInt16(data, e.Value.(int16))
		} else if arrType == "uint32" {
			SetUint32(data, e.Value.(uint32))
		} else if arrType == "int32" {
			SetInt32(data, e.Value.(int32))
		} else if arrType == "uint64" {
			SetUint64(data, e.Value.(uint64))
		} else if arrType == "int64" {
			SetInt64(data, e.Value.(int64))
		} else if arrType == "string" {
			SetString(data, e.Value.(string))
		} else if arrType == "float" {
			SetFloat(data, e.Value.(float32))
		} else if arrType == "buffer" {
			SetBuffer(data, e.Value.([]byte))
		} else {
			SetEntity(data, e.Value.(Msg))
		}
	}
}

func GetEntity(data *bytes.Buffer, entityType string) Msg {
	buffer := GetBuffer(data)
	entity := reflect.New(msgObjectMap_ByName[entityType]).Interface().(Msg)
	entity.Decode(buffer)
	return entity
}

func GetArray(data *bytes.Buffer, arrType string) *list.List {
	l := list.New()
	len := GetUint16(data)

	for i := 0; i < int(len); i++ {
		if arrType == "uint8" {
			l.PushBack(GetUint8(data))
		} else if arrType == "int8" {
			l.PushBack(GetInt8(data))
		} else if arrType == "uint16" {
			l.PushBack(GetUint16(data))
		} else if arrType == "int16" {
			l.PushBack(GetInt16(data))
		} else if arrType == "uint32" {
			l.PushBack(GetUint32(data))
		} else if arrType == "int32" {
			l.PushBack(GetInt32(data))
		} else if arrType == "uint64" {
			l.PushBack(GetUint64(data))
		} else if arrType == "int64" {
			l.PushBack(GetInt64(data))
		} else if arrType == "string" {
			l.PushBack(GetString(data))
		} else if arrType == "float" {
			l.PushBack(GetFloat(data))
		} else if arrType == "buffer" {
			l.PushBack(GetBuffer(data))
		} else {
			l.PushBack(GetEntity(data, arrType))
		}
	}

	return l
}

func GetUint8(data *bytes.Buffer) uint8 {
	var num uint8
	binary.Read(data, defaultOrder, &num)
	return num
}

func GetInt8(data *bytes.Buffer) int8 {
	var num int8
	binary.Read(data, defaultOrder, &num)
	return num
}

func GetUint16(data *bytes.Buffer) uint16 {
	var num uint16
	binary.Read(data, defaultOrder, &num)
	return num
}

func GetInt16(data *bytes.Buffer) int16 {
	var num int16
	binary.Read(data, defaultOrder, &num)
	return num
}

func GetUint32(data *bytes.Buffer) uint32 {
	var num uint32
	binary.Read(data, defaultOrder, &num)
	return num
}

func GetInt32(data *bytes.Buffer) int32 {
	var num int32
	binary.Read(data, defaultOrder, &num)
	return num
}

func GetUint64(data *bytes.Buffer) uint64 {
	var num uint64
	binary.Read(data, defaultOrder, &num)
	return num
}

func GetInt64(data *bytes.Buffer) int64 {
	//var num int64;
	//binary.Read(data, defaultOrder, &num);
	//return num;

	var num float64
	binary.Read(data, defaultOrder, &num)
	return int64(num)
}

func GetFloat(data *bytes.Buffer) float32 {
	var num float32
	binary.Read(data, defaultOrder, &num)

	return num
}

func GetString(data *bytes.Buffer) string {
	GetUint8(data)
	len := GetUint16(data)

	strBytes := make([]byte, len)
	binary.Read(data, defaultOrder, &strBytes)
	return string(strBytes)
}

func GetBuffer(data *bytes.Buffer) []byte {
	GetUint8(data)
	len := GetUint16(data)

	buffer := make([]byte, len)
	binary.Read(data, defaultOrder, &buffer)
	return buffer
}
