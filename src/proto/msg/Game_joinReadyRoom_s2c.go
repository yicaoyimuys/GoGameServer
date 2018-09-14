package msg

import "proto"
import "bytes"
import "container/list"

type Game_joinReadyRoom_s2c struct {
	MsgId uint16
	Users *list.List
}

func NewGame_joinReadyRoom_s2c() *Game_joinReadyRoom_s2c {
	return &Game_joinReadyRoom_s2c{
		MsgId: ID_Game_joinReadyRoom_s2c,
	}
}

func (this *Game_joinReadyRoom_s2c) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetArray(buf, this.Users, "UserInfo")

	return buf.Bytes()
}

func (this *Game_joinReadyRoom_s2c) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.Users = proto.GetArray(buf, "UserInfo")

}
