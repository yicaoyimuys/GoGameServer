package module

import (
	"time"

	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/protos"
	"github.com/yicaoyimuys/GoGameServer/core/libs/random"
	"github.com/yicaoyimuys/GoGameServer/core/libs/sessions"
	"github.com/yicaoyimuys/GoGameServer/servives/login/cache"
	"github.com/yicaoyimuys/GoGameServer/servives/public"
	"github.com/yicaoyimuys/GoGameServer/servives/public/gameProto"
	"github.com/yicaoyimuys/GoGameServer/servives/public/mysqlModels"
	"github.com/yicaoyimuys/GoGameServer/servives/public/redisCaches"
	"go.uber.org/zap"

	"google.golang.org/protobuf/proto"
)

// 登录
func Login(clientSession *sessions.BackSession, msgData proto.Message) {
	data := msgData.(*gameProto.UserLoginC2S)
	account := data.GetAccount()

	onlineUser := cache.GetOnlineUserByAccount(account)
	if onlineUser != nil {
		oldClientSession := onlineUser.Session
		if oldClientSession.ID() != clientSession.ID() {
			//当前在线，但是连接不同，其他客户端连接，需通知当前客户端下线
			sendOtherLogin(oldClientSession)
			//替换Session
			cache.RemoveOnlineUser(oldClientSession.ID())
			//登录成功后处理
			loginSuccess(clientSession, onlineUser.Account, onlineUser.UserID)
		}
	} else {
		//进行DB登录
		dbUser := login(account)
		//登录成功后处理
		loginSuccess(clientSession, dbUser.Account, dbUser.Id)
	}
}

func login(account string) *mysqlModels.User {
	//db中获取用户数据
	dbUser := mysqlModels.GetUser(account)
	if dbUser == nil {
		//注册新用户
		addMoney := random.RandomInt31n(999)
		dbUser = mysqlModels.AddUser(account, addMoney)
	} else {
		//更新用户最后登录时间
		dbUser.LastLoginTime = time.Now().Unix()
		mysqlModels.UpdateUserLoginTime(dbUser.Id, dbUser.LastLoginTime)
	}
	//加入redis缓存
	redisCaches.SetUser(dbUser)
	return dbUser
}

// 登录成功后处理
func loginSuccess(clientSession *sessions.BackSession, account string, userID uint64) {
	//缓存用户在线数据
	cache.AddOnlineUser(userID, account, clientSession)
	clientSession.AddCloseCallback(nil, "user.loginSuccess", func() {
		cache.RemoveOnlineUser(clientSession.ID())
		DEBUG("用户下线", zap.Int32("OnlineUsersNum", cache.GetOnlineUsersNum()))
	})
	DEBUG("用户上线", zap.Int32("OnlineUsersNum", cache.GetOnlineUsersNum()))

	//返回客户端数据
	token := public.CreateToken(userID)
	sendMsg := &gameProto.UserLoginS2C{
		Token: protos.String(token),
	}
	public.SendMsgToClient(clientSession, sendMsg)
}

func sendOtherLogin(clientSession *sessions.BackSession) {
	sendMsg := &gameProto.UserOtherLoginNoticeS2C{}
	public.SendMsgToClient(clientSession, sendMsg)
}
