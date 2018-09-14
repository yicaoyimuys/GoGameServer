package msg

import "connectorServer/proto"
import "bytes"



type Game_join_s2c struct {
	MsgId	uint16
	MyUser	*UserInfo
	OtherUser	*UserInfo
	IsAi	uint16
	AiLevel	uint16

}


func NewGame_join_s2c() *Game_join_s2c {
	return &Game_join_s2c{
		MsgId: 	ID_Game_join_s2c,
	}
}


func (this *Game_join_s2c) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetEntity(buf, this.MyUser)
	proto.SetEntity(buf, this.OtherUser)
	proto.SetUint16(buf, this.IsAi)
	proto.SetUint16(buf, this.AiLevel)

	return buf.Bytes()
}

func (this *Game_join_s2c) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.MyUser = proto.GetEntity(buf, "UserInfo").(*UserInfo)
	this.OtherUser = proto.GetEntity(buf, "UserInfo").(*UserInfo)
	this.IsAi = proto.GetUint16(buf)
	this.AiLevel = proto.GetUint16(buf)

}