package msg

import "proto"
import "bytes"

type Platform_login_s2c struct {
	MsgId uint16
}

func NewPlatform_login_s2c() *Platform_login_s2c {
	return &Platform_login_s2c{
		MsgId: ID_Platform_login_s2c,
	}
}

func (this *Platform_login_s2c) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)

	return buf.Bytes()
}

func (this *Platform_login_s2c) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)

}
