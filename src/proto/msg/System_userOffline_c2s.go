package msg

import "proto"
import "bytes"

type System_userOffline_c2s struct {
	MsgId uint16
}

func NewSystem_userOffline_c2s() *System_userOffline_c2s {
	return &System_userOffline_c2s{
		MsgId: ID_System_userOffline_c2s,
	}
}

func (this *System_userOffline_c2s) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)

	return buf.Bytes()
}

func (this *System_userOffline_c2s) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)

}
