package msg

import "connectorServer/proto"
import "bytes"



type Game_matching_notice_s2c struct {
	MsgId	uint16
	GameId	uint16
	OtherUserId	string
	OtherUserName	string
	OtherUserPic	string
	OtherUserSex	uint8
	IsAi	uint8
	AiLevel	uint16
	RoomId	string

}


func NewGame_matching_notice_s2c() *Game_matching_notice_s2c {
	return &Game_matching_notice_s2c{
		MsgId: 	ID_Game_matching_notice_s2c,
	}
}


func (this *Game_matching_notice_s2c) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetUint16(buf, this.GameId)
	proto.SetString(buf, this.OtherUserId)
	proto.SetString(buf, this.OtherUserName)
	proto.SetString(buf, this.OtherUserPic)
	proto.SetUint8(buf, this.OtherUserSex)
	proto.SetUint8(buf, this.IsAi)
	proto.SetUint16(buf, this.AiLevel)
	proto.SetString(buf, this.RoomId)

	return buf.Bytes()
}

func (this *Game_matching_notice_s2c) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.GameId = proto.GetUint16(buf)
	this.OtherUserId = proto.GetString(buf)
	this.OtherUserName = proto.GetString(buf)
	this.OtherUserPic = proto.GetString(buf)
	this.OtherUserSex = proto.GetUint8(buf)
	this.IsAi = proto.GetUint8(buf)
	this.AiLevel = proto.GetUint16(buf)
	this.RoomId = proto.GetString(buf)

}