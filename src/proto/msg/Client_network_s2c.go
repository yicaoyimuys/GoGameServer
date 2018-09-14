package msg

import "proto"
import "bytes"

type Client_network_s2c struct {
	MsgId uint16
	Time  uint32
}

func NewClient_network_s2c() *Client_network_s2c {
	return &Client_network_s2c{
		MsgId: ID_Client_network_s2c,
	}
}

func (this *Client_network_s2c) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetUint32(buf, this.Time)

	return buf.Bytes()
}

func (this *Client_network_s2c) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.Time = proto.GetUint32(buf)

}
