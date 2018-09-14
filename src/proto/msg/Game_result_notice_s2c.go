package msg

import "proto"
import "bytes"

type Game_result_notice_s2c struct {
	MsgId        uint16
	WinUserId    string
	UserExitFlag uint8
}

func NewGame_result_notice_s2c() *Game_result_notice_s2c {
	return &Game_result_notice_s2c{
		MsgId: ID_Game_result_notice_s2c,
	}
}

func (this *Game_result_notice_s2c) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetString(buf, this.WinUserId)
	proto.SetUint8(buf, this.UserExitFlag)

	return buf.Bytes()
}

func (this *Game_result_notice_s2c) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.WinUserId = proto.GetString(buf)
	this.UserExitFlag = proto.GetUint8(buf)

}
