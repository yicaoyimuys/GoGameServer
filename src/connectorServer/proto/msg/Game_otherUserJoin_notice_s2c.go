package msg

import "connectorServer/proto"
import "bytes"



type Game_otherUserJoin_notice_s2c struct {
	MsgId	uint16
	OtherUser	*UserInfo

}


func NewGame_otherUserJoin_notice_s2c() *Game_otherUserJoin_notice_s2c {
	return &Game_otherUserJoin_notice_s2c{
		MsgId: 	ID_Game_otherUserJoin_notice_s2c,
	}
}


func (this *Game_otherUserJoin_notice_s2c) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetEntity(buf, this.OtherUser)

	return buf.Bytes()
}

func (this *Game_otherUserJoin_notice_s2c) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.OtherUser = proto.GetEntity(buf, "UserInfo").(*UserInfo)

}