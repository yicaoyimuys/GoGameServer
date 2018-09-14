package msg

import "connectorServer/proto"
import "bytes"



type Game_trunNumNotice_s2c struct {
	MsgId	uint16
	TrunNum	uint32

}


func NewGame_trunNumNotice_s2c() *Game_trunNumNotice_s2c {
	return &Game_trunNumNotice_s2c{
		MsgId: 	ID_Game_trunNumNotice_s2c,
	}
}


func (this *Game_trunNumNotice_s2c) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetUint32(buf, this.TrunNum)

	return buf.Bytes()
}

func (this *Game_trunNumNotice_s2c) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.TrunNum = proto.GetUint32(buf)

}