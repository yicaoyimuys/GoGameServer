package msg

import "proto"
import "bytes"

type Game_loadProgress_notice_s2c struct {
	MsgId    uint16
	UserId   string
	Progress uint16
}

func NewGame_loadProgress_notice_s2c() *Game_loadProgress_notice_s2c {
	return &Game_loadProgress_notice_s2c{
		MsgId: ID_Game_loadProgress_notice_s2c,
	}
}

func (this *Game_loadProgress_notice_s2c) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetString(buf, this.UserId)
	proto.SetUint16(buf, this.Progress)

	return buf.Bytes()
}

func (this *Game_loadProgress_notice_s2c) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.UserId = proto.GetString(buf)
	this.Progress = proto.GetUint16(buf)

}
