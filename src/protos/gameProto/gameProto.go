package gameProto

import (
	"protos"
)

//初始化消息ID和消息类型的对应关系
func init() {
	protos.SetMsg(ID_client_ping_c2s, ClientPingC2S{})
	protos.SetMsg(ID_error_notice_s2c, ErrorNoticeS2C{})

	protos.SetMsg(ID_user_login_c2s, UserLoginC2S{})
	protos.SetMsg(ID_user_login_s2c, UserLoginS2C{})

	protos.SetMsg(ID_user_getInfo_c2s, UserGetInfoC2S{})
	protos.SetMsg(ID_user_getInfo_s2c, UserGetInfoS2C{})
}
