package msg

import "connectorServer/proto"
import "bytes"



type Client_network_c2s struct {
	MsgId	uint16
	Time	uint32

}


func NewClient_network_c2s() *Client_network_c2s {
	return &Client_network_c2s{
		MsgId: 	ID_Client_network_c2s,
	}
}


func (this *Client_network_c2s) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetUint32(buf, this.Time)

	return buf.Bytes()
}

func (this *Client_network_c2s) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.Time = proto.GetUint32(buf)

}