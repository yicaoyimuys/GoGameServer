package msg

import "proto"

const (
	ID_Client_ping_c2s                   uint16 = 1000
	ID_Client_network_c2s                uint16 = 1001
	ID_Client_network_s2c                uint16 = 1002
	ID_Error_notice_s2c                  uint16 = 1999
	ID_Platform_login_c2s                uint16 = 2001
	ID_Platform_login_s2c                uint16 = 2002
	ID_Game_join_c2s                     uint16 = 3101
	ID_Game_join_s2c                     uint16 = 3102
	ID_Game_otherUserJoin_notice_s2c     uint16 = 3104
	ID_Game_loadProgress_c2s             uint16 = 3001
	ID_Game_loadProgress_notice_s2c      uint16 = 3002
	ID_Game_start_notice_s2c             uint16 = 3004
	ID_Game_exit_c2s                     uint16 = 3005
	ID_Game_result_c2s                   uint16 = 3007
	ID_Game_result_notice_s2c            uint16 = 3008
	ID_Game_event_c2s                    uint16 = 3009
	ID_Game_event_notice_s2c             uint16 = 3010
	ID_Game_trunNumNotice_s2c            uint16 = 3012
	ID_Game_matching_c2s                 uint16 = 3013
	ID_Game_matching_s2c                 uint16 = 3014
	ID_Game_matching_notice_s2c          uint16 = 3016
	ID_Game_cancelMatching_c2s           uint16 = 3017
	ID_Game_cancelMatching_s2c           uint16 = 3018
	ID_Game_createReadyRoom_c2s          uint16 = 4101
	ID_Game_createReadyRoom_s2c          uint16 = 4102
	ID_Game_joinReadyRoom_c2s            uint16 = 4103
	ID_Game_joinReadyRoom_s2c            uint16 = 4104
	ID_Game_joinReadyRoom_notice_s2c     uint16 = 4106
	ID_Game_leaveReadyRoom_c2s           uint16 = 4107
	ID_Game_leaveReadyRoom_notice_s2c    uint16 = 4108
	ID_Game_dissolveReadyRoom_c2s        uint16 = 4109
	ID_Game_dissolveReadyRoom_notice_s2c uint16 = 4110
	ID_Game_startByReadyRoom_c2s         uint16 = 4111
	ID_Game_startByReadyRoom_notice_s2c  uint16 = 4112
	ID_Game_refuseReadyRoom_c2s          uint16 = 4113
	ID_Game_refuseReadyRoom_notice_s2c   uint16 = 4114
	ID_Game_againJoinReadyRoom_c2s       uint16 = 4115
)

func init() {
	proto.SetMsgByName("Client_ping_c2s", Client_ping_c2s{})
	proto.SetMsgById(1000, Client_ping_c2s{})
	proto.SetMsgByName("Client_network_c2s", Client_network_c2s{})
	proto.SetMsgById(1001, Client_network_c2s{})
	proto.SetMsgByName("Client_network_s2c", Client_network_s2c{})
	proto.SetMsgById(1002, Client_network_s2c{})
	proto.SetMsgByName("Error_notice_s2c", Error_notice_s2c{})
	proto.SetMsgById(1999, Error_notice_s2c{})
	proto.SetMsgByName("Platform_login_c2s", Platform_login_c2s{})
	proto.SetMsgById(2001, Platform_login_c2s{})
	proto.SetMsgByName("Platform_login_s2c", Platform_login_s2c{})
	proto.SetMsgById(2002, Platform_login_s2c{})
	proto.SetMsgByName("UserInfo", UserInfo{})
	proto.SetMsgByName("Game_join_c2s", Game_join_c2s{})
	proto.SetMsgById(3101, Game_join_c2s{})
	proto.SetMsgByName("Game_join_s2c", Game_join_s2c{})
	proto.SetMsgById(3102, Game_join_s2c{})
	proto.SetMsgByName("Game_otherUserJoin_notice_s2c", Game_otherUserJoin_notice_s2c{})
	proto.SetMsgById(3104, Game_otherUserJoin_notice_s2c{})
	proto.SetMsgByName("Game_loadProgress_c2s", Game_loadProgress_c2s{})
	proto.SetMsgById(3001, Game_loadProgress_c2s{})
	proto.SetMsgByName("Game_loadProgress_notice_s2c", Game_loadProgress_notice_s2c{})
	proto.SetMsgById(3002, Game_loadProgress_notice_s2c{})
	proto.SetMsgByName("Game_start_notice_s2c", Game_start_notice_s2c{})
	proto.SetMsgById(3004, Game_start_notice_s2c{})
	proto.SetMsgByName("Game_exit_c2s", Game_exit_c2s{})
	proto.SetMsgById(3005, Game_exit_c2s{})
	proto.SetMsgByName("Game_result_c2s", Game_result_c2s{})
	proto.SetMsgById(3007, Game_result_c2s{})
	proto.SetMsgByName("Game_result_notice_s2c", Game_result_notice_s2c{})
	proto.SetMsgById(3008, Game_result_notice_s2c{})
	proto.SetMsgByName("Game_event_c2s", Game_event_c2s{})
	proto.SetMsgById(3009, Game_event_c2s{})
	proto.SetMsgByName("Game_event_notice_s2c", Game_event_notice_s2c{})
	proto.SetMsgById(3010, Game_event_notice_s2c{})
	proto.SetMsgByName("Game_trunNumNotice_s2c", Game_trunNumNotice_s2c{})
	proto.SetMsgById(3012, Game_trunNumNotice_s2c{})
	proto.SetMsgByName("Game_matching_c2s", Game_matching_c2s{})
	proto.SetMsgById(3013, Game_matching_c2s{})
	proto.SetMsgByName("Game_matching_s2c", Game_matching_s2c{})
	proto.SetMsgById(3014, Game_matching_s2c{})
	proto.SetMsgByName("Game_matching_notice_s2c", Game_matching_notice_s2c{})
	proto.SetMsgById(3016, Game_matching_notice_s2c{})
	proto.SetMsgByName("Game_cancelMatching_c2s", Game_cancelMatching_c2s{})
	proto.SetMsgById(3017, Game_cancelMatching_c2s{})
	proto.SetMsgByName("Game_cancelMatching_s2c", Game_cancelMatching_s2c{})
	proto.SetMsgById(3018, Game_cancelMatching_s2c{})
	proto.SetMsgByName("Game_createReadyRoom_c2s", Game_createReadyRoom_c2s{})
	proto.SetMsgById(4101, Game_createReadyRoom_c2s{})
	proto.SetMsgByName("Game_createReadyRoom_s2c", Game_createReadyRoom_s2c{})
	proto.SetMsgById(4102, Game_createReadyRoom_s2c{})
	proto.SetMsgByName("Game_joinReadyRoom_c2s", Game_joinReadyRoom_c2s{})
	proto.SetMsgById(4103, Game_joinReadyRoom_c2s{})
	proto.SetMsgByName("Game_joinReadyRoom_s2c", Game_joinReadyRoom_s2c{})
	proto.SetMsgById(4104, Game_joinReadyRoom_s2c{})
	proto.SetMsgByName("Game_joinReadyRoom_notice_s2c", Game_joinReadyRoom_notice_s2c{})
	proto.SetMsgById(4106, Game_joinReadyRoom_notice_s2c{})
	proto.SetMsgByName("Game_leaveReadyRoom_c2s", Game_leaveReadyRoom_c2s{})
	proto.SetMsgById(4107, Game_leaveReadyRoom_c2s{})
	proto.SetMsgByName("Game_leaveReadyRoom_notice_s2c", Game_leaveReadyRoom_notice_s2c{})
	proto.SetMsgById(4108, Game_leaveReadyRoom_notice_s2c{})
	proto.SetMsgByName("Game_dissolveReadyRoom_c2s", Game_dissolveReadyRoom_c2s{})
	proto.SetMsgById(4109, Game_dissolveReadyRoom_c2s{})
	proto.SetMsgByName("Game_dissolveReadyRoom_notice_s2c", Game_dissolveReadyRoom_notice_s2c{})
	proto.SetMsgById(4110, Game_dissolveReadyRoom_notice_s2c{})
	proto.SetMsgByName("Game_startByReadyRoom_c2s", Game_startByReadyRoom_c2s{})
	proto.SetMsgById(4111, Game_startByReadyRoom_c2s{})
	proto.SetMsgByName("Game_startByReadyRoom_notice_s2c", Game_startByReadyRoom_notice_s2c{})
	proto.SetMsgById(4112, Game_startByReadyRoom_notice_s2c{})
	proto.SetMsgByName("Game_refuseReadyRoom_c2s", Game_refuseReadyRoom_c2s{})
	proto.SetMsgById(4113, Game_refuseReadyRoom_c2s{})
	proto.SetMsgByName("Game_refuseReadyRoom_notice_s2c", Game_refuseReadyRoom_notice_s2c{})
	proto.SetMsgById(4114, Game_refuseReadyRoom_notice_s2c{})
	proto.SetMsgByName("Game_againJoinReadyRoom_c2s", Game_againJoinReadyRoom_c2s{})
	proto.SetMsgById(4115, Game_againJoinReadyRoom_c2s{})

}
