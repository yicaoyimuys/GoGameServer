package module

import (
	"core/libs/sessions"
	"core/protos"
	"core/protos/gameProto"
	"github.com/golang/protobuf/proto"
	"servives/public"
	"servives/public/errCodes"
	"servives/public/redisCaches"
)

//获取用户信息
func GetInfo(clientSession *sessions.BackSession, msgData proto.Message) {
	data := msgData.(*gameProto.UserGetInfoC2S)
	token := data.GetToken()
	userId := public.GetUserIdByToken(token)
	if userId == 0 {
		public.SendErrorMsgToClient(clientSession, errCodes.PARAM_ERROR)
		return
	}

	//获取缓存中用户数据
	dbUser := redisCaches.GetUser(userId)
	if dbUser == nil {
		public.SendErrorMsgToClient(clientSession, errCodes.PARAM_ERROR)
		return
	}

	//返回客户端消息
	sendMsg := &gameProto.UserGetInfoS2C{
		Data: &gameProto.UserInfo{
			Id:    protos.Uint64(dbUser.Id),
			Name:  protos.String(dbUser.Account),
			Money: protos.Int32(dbUser.Money),
		},
	}
	public.SendMsgToClient(clientSession, sendMsg)
}
