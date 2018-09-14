package msg

import "proto"
import "bytes"

type Client_ping_c2s struct {
	MsgId uint16
}

func NewClient_ping_c2s() *Client_ping_c2s {
	return &Client_ping_c2s{
		MsgId: ID_Client_ping_c2s,
	}
}

func (this *Client_ping_c2s) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)

	return buf.Bytes()
}

func (this *Client_ping_c2s) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)

}
