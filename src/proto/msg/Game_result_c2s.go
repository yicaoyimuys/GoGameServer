package msg

import "proto"
import "bytes"

type Game_result_c2s struct {
	MsgId  uint16
	RoomId string
	Result uint8
}

func NewGame_result_c2s() *Game_result_c2s {
	return &Game_result_c2s{
		MsgId: ID_Game_result_c2s,
	}
}

func (this *Game_result_c2s) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetString(buf, this.RoomId)
	proto.SetUint8(buf, this.Result)

	return buf.Bytes()
}

func (this *Game_result_c2s) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.RoomId = proto.GetString(buf)
	this.Result = proto.GetUint8(buf)

}
