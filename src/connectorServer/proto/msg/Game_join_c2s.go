package msg

import "connectorServer/proto"
import "bytes"



type Game_join_c2s struct {
	MsgId	uint16
	RoomId	string

}


func NewGame_join_c2s() *Game_join_c2s {
	return &Game_join_c2s{
		MsgId: 	ID_Game_join_c2s,
	}
}


func (this *Game_join_c2s) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetString(buf, this.RoomId)

	return buf.Bytes()
}

func (this *Game_join_c2s) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.RoomId = proto.GetString(buf)

}