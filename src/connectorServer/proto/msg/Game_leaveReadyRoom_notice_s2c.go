package msg

import "connectorServer/proto"
import "bytes"



type Game_leaveReadyRoom_notice_s2c struct {
	MsgId	uint16
	UserId	string

}


func NewGame_leaveReadyRoom_notice_s2c() *Game_leaveReadyRoom_notice_s2c {
	return &Game_leaveReadyRoom_notice_s2c{
		MsgId: 	ID_Game_leaveReadyRoom_notice_s2c,
	}
}


func (this *Game_leaveReadyRoom_notice_s2c) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetString(buf, this.UserId)

	return buf.Bytes()
}

func (this *Game_leaveReadyRoom_notice_s2c) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.UserId = proto.GetString(buf)

}