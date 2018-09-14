package msg

import "connectorServer/proto"
import "bytes"



type Game_againJoinReadyRoom_c2s struct {
	MsgId	uint16
	RoomId	string
	UserId	string
	UserName	string
	UserPic	string
	UserSex	uint8

}


func NewGame_againJoinReadyRoom_c2s() *Game_againJoinReadyRoom_c2s {
	return &Game_againJoinReadyRoom_c2s{
		MsgId: 	ID_Game_againJoinReadyRoom_c2s,
	}
}


func (this *Game_againJoinReadyRoom_c2s) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetString(buf, this.RoomId)
	proto.SetString(buf, this.UserId)
	proto.SetString(buf, this.UserName)
	proto.SetString(buf, this.UserPic)
	proto.SetUint8(buf, this.UserSex)

	return buf.Bytes()
}

func (this *Game_againJoinReadyRoom_c2s) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.RoomId = proto.GetString(buf)
	this.UserId = proto.GetString(buf)
	this.UserName = proto.GetString(buf)
	this.UserPic = proto.GetString(buf)
	this.UserSex = proto.GetUint8(buf)

}