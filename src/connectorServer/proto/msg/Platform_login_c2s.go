package msg

import "connectorServer/proto"
import "bytes"



type Platform_login_c2s struct {
	MsgId	uint16
	GameId	uint16
	PlatformId	uint16
	PlatformData	string

}


func NewPlatform_login_c2s() *Platform_login_c2s {
	return &Platform_login_c2s{
		MsgId: 	ID_Platform_login_c2s,
	}
}


func (this *Platform_login_c2s) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetUint16(buf, this.MsgId)
	proto.SetUint16(buf, this.GameId)
	proto.SetUint16(buf, this.PlatformId)
	proto.SetString(buf, this.PlatformData)

	return buf.Bytes()
}

func (this *Platform_login_c2s) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.MsgId = proto.GetUint16(buf)
	this.GameId = proto.GetUint16(buf)
	this.PlatformId = proto.GetUint16(buf)
	this.PlatformData = proto.GetString(buf)

}