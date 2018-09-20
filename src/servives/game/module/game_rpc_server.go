package module

type GameRpcServer struct {
}

type RpcCustomRoomEventReq struct {
	RoomId string
	UserId string
	Event  string
}

type RpcCustomRoomEventRes struct {
	Code uint16
}

func (this *GameRpcServer) CustomRoomEvent(args *RpcCustomRoomEventReq, reply *RpcCustomRoomEventRes) error {
	//roomInfo := datas.GetRoom(args.RoomId)
	//if roomInfo == nil {
	//	return errors.New("房间不存在")
	//}
	//
	//roomUser := getRoomUser(roomInfo, args.UserId)
	//if roomUser == nil {
	//	return errors.New("用户不存在")
	//}
	//
	////通知前端
	//var sendMsg = msg.NewGame_event_notice_s2c()
	//sendMsg.UserId = args.UserId
	//sendMsg.Event = args.Event
	//NoticeRoomMsg(roomInfo, sendMsg, "")

	return nil
}
