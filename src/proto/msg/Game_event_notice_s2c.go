package msg

import "proto"
import "bytes"

type Game_event_notice_s2c struct {
	MsgId   uint16
	UserId  string
	Event   string
	TrunNum uint32
}

func NewGame_event_notice_s2c() *Game_event_notice_s2c {
	return &Game_event_notice_s2c{
		MsgId: ID_Game_event_notice_s2c,
	}
}

func (this *Game_event_notice_s2c) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetString(buf, this.UserId)
	proto.SetString(buf, this.Event)
	proto.SetUint32(buf, this.TrunNum)

	return buf.Bytes()
}

func (this *Game_event_notice_s2c) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.UserId = proto.GetString(buf)
	this.Event = proto.GetString(buf)
	this.TrunNum = proto.GetUint32(buf)

}
