package msg

import "connectorServer/proto"
import "bytes"



type Game_loadProgress_c2s struct {
	MsgId	uint16
	RoomId	string
	Progress	uint16

}


func NewGame_loadProgress_c2s() *Game_loadProgress_c2s {
	return &Game_loadProgress_c2s{
		MsgId: 	ID_Game_loadProgress_c2s,
	}
}


func (this *Game_loadProgress_c2s) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetString(buf, this.RoomId)
	proto.SetUint16(buf, this.Progress)

	return buf.Bytes()
}

func (this *Game_loadProgress_c2s) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.RoomId = proto.GetString(buf)
	this.Progress = proto.GetUint16(buf)

}