package msg

import "proto"
import "bytes"

type Game_event_c2s struct {
	MsgId    uint16
	RoomId   string
	SendType uint8
	Event    string
}

func NewGame_event_c2s() *Game_event_c2s {
	return &Game_event_c2s{
		MsgId: ID_Game_event_c2s,
	}
}

func (this *Game_event_c2s) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetString(buf, this.RoomId)
	proto.SetUint8(buf, this.SendType)
	proto.SetString(buf, this.Event)

	return buf.Bytes()
}

func (this *Game_event_c2s) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.RoomId = proto.GetString(buf)
	this.SendType = proto.GetUint8(buf)
	this.Event = proto.GetString(buf)

}
