package msg

import "connectorServer/proto"

const (
	ID_System_userOffline_c2s uint16 = 1;

)

func init() {
	proto.SetMsgByName("System_userOffline_c2s", System_userOffline_c2s{})
	proto.SetMsgById(1, System_userOffline_c2s{})

}