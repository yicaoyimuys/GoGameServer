package msg

import "proto"
import "bytes"

type Game_cancelMatching_s2c struct {
	MsgId uint16
}

func NewGame_cancelMatching_s2c() *Game_cancelMatching_s2c {
	return &Game_cancelMatching_s2c{
		MsgId: ID_Game_cancelMatching_s2c,
	}
}

func (this *Game_cancelMatching_s2c) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)

	return buf.Bytes()
}

func (this *Game_cancelMatching_s2c) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)

}
