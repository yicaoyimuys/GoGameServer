package msg

import "proto"
import "bytes"

type Game_createReadyRoom_s2c struct {
	MsgId  uint16
	RoomId string
}

func NewGame_createReadyRoom_s2c() *Game_createReadyRoom_s2c {
	return &Game_createReadyRoom_s2c{
		MsgId: ID_Game_createReadyRoom_s2c,
	}
}

func (this *Game_createReadyRoom_s2c) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetString(buf, this.RoomId)

	return buf.Bytes()
}

func (this *Game_createReadyRoom_s2c) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.RoomId = proto.GetString(buf)

}
