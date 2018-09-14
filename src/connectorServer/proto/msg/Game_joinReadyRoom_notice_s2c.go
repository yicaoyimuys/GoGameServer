package msg

import "connectorServer/proto"
import "bytes"



type Game_joinReadyRoom_notice_s2c struct {
	MsgId	uint16
	User	*UserInfo

}


func NewGame_joinReadyRoom_notice_s2c() *Game_joinReadyRoom_notice_s2c {
	return &Game_joinReadyRoom_notice_s2c{
		MsgId: 	ID_Game_joinReadyRoom_notice_s2c,
	}
}


func (this *Game_joinReadyRoom_notice_s2c) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetEntity(buf, this.User)

	return buf.Bytes()
}

func (this *Game_joinReadyRoom_notice_s2c) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.User = proto.GetEntity(buf, "UserInfo").(*UserInfo)

}