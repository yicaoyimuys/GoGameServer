package msg

import "proto"
import "bytes"

type UserInfo struct {
	Uuid         string
	UserId       string
	UserName     string
	UserPic      string
	UserSex      uint8
	UserProgress uint16
}

func NewUserInfo() *UserInfo {
	return &UserInfo{}
}

func (this *UserInfo) Encode() []byte {
	buf := new(bytes.Buffer)
	proto.SetString(buf, this.Uuid)
	proto.SetString(buf, this.UserId)
	proto.SetString(buf, this.UserName)
	proto.SetString(buf, this.UserPic)
	proto.SetUint8(buf, this.UserSex)
	proto.SetUint16(buf, this.UserProgress)

	return buf.Bytes()
}

func (this *UserInfo) Decode(msg []byte) {
	buf := bytes.NewBuffer(msg)
	this.Uuid = proto.GetString(buf)
	this.UserId = proto.GetString(buf)
	this.UserName = proto.GetString(buf)
	this.UserPic = proto.GetString(buf)
	this.UserSex = proto.GetUint8(buf)
	this.UserProgress = proto.GetUint16(buf)

}
