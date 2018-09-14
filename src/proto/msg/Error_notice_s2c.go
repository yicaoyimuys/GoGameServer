package msg

import "proto"
import "bytes"

type Error_notice_s2c struct {
	MsgId     uint16
	ErrorCode uint32
}

func NewError_notice_s2c() *Error_notice_s2c {
	return &Error_notice_s2c{
		MsgId: ID_Error_notice_s2c,
	}
}

func (this *Error_notice_s2c) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetUint32(buf, this.ErrorCode)

	return buf.Bytes()
}

func (this *Error_notice_s2c) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.ErrorCode = proto.GetUint32(buf)

}
