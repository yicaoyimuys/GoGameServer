package user

import (
	"github.com/funny/link"
	. "model"
	"module"
	"protos/gameProto"
	"proxys/dbProxy"
	"proxys/redisProxy"
	"proxys/transferProxy"
	. "tools"
	"proxys/logProxy"
	"proxys/gameProxy"
	"tools/timer"
	"tools/debug"
	"time"
	"container/list"
	"strconv"
)

const (
	//处理下线用户数据间隔(20分钟)
	DEAL_OFFLINEUSER_INTERVAL = 20 * 60
	//用户下线后数据存在时长(5小时)
	OFFLINEUSER_TIME = 5 * 60 * 60
)

type UserModule struct {
}

// 在初始化的时候将模块注册到module包
func init() {
	module.User = UserModule{}
}

//用户DB登录返回
func (this UserModule) UserLoginHandle(session *link.Session, userName string, userID uint64) {
	if userID == 0 {
		gameProxy.SendLoginResult(session, 0)
	} else {
		//登录成功处理
		success := this.LoginSuccess(session, userName, userID, 0)
		if success {
			//登录成功后处理
			this.dealLoginSuccess(session, userName, userID)
		} else {
			gameProxy.SendLoginResult(session, 0)
		}
	}
}

//登录
func (this UserModule) Login(userName string, session *link.Session) {
	onlineUser := module.Cache.GetOnlineUserByUserName(userName)
	if onlineUser != nil {
		if onlineUser.Session.Id() != session.Id() {
			//当前在线，但是连接不同，其他客户端连接，需通知当前客户端下线
			gameProxy.SendOtherLogin(onlineUser.Session)
			//替换Session
			module.Cache.RemoveOnlineUser(onlineUser.Session.Id())
			//登录成功处理
			success := this.LoginSuccess(session, onlineUser.UserName, onlineUser.UserID, 0)
			if success {
				//登录成功后处理
				this.dealLoginSuccess(session, userName, onlineUser.UserID)
			} else {
				gameProxy.SendLoginResult(session, 0)
			}
		}
	} else {
		dbProxy.UserLogin(session.Id(), userName)
	}
}

//登录成功后处理
func (this UserModule) dealLoginSuccess(session *link.Session, userName string, userID uint64){
	//通知GameServer登录成功
	transferProxy.SetClientLoginSuccess(userName, userID, session)
	//发送登录成功消息
	gameProxy.SendLoginResult(session, userID)
	//用户下线时处理
	session.AddCloseCallback(session, func() {
		//记录用户下线Log
		logProxy.UserOffLine(userID)
	})
	//记录用户登录Log
	logProxy.UserLogin(userID)
}

//用户登录成功处理
func (this UserModule) LoginSuccess(session *link.Session, userName string, userID uint64, gameServerID uint32) bool {
	cacheSuccess := module.Cache.AddOnlineUser(userName, userID, session, gameServerID)
	if cacheSuccess {
		session.AddCloseCallback(session, func() {
			module.Cache.RemoveOnlineUser(session.Id())
			DEBUG("用户下线：当前在线人数", module.Cache.GetOnlineUsersNum())
		})
		DEBUG("用户上线：当前在线人数", module.Cache.GetOnlineUsersNum())
		return true
	} else {
		ERR("what????", userName)
		return false
	}
}

//重新连接
func (this UserModule) AgainConnect(oldSessionID uint64, session *link.Session) uint64 {
	//	if oldSessionID == session.Id() {
	//		return 0
	//	}

	//	user := module.Cache.GetOnlineUserBySession(oldSessionID)
	//	if user == nil {
	//		return 0
	//	}

	//	module.Cache.RemoveOnlineUser(oldSessionID)

	//	cacheSuccess := module.Cache.AddOnlineUser(user.UserName, user.UserID, session)
	//	if cacheSuccess {
	//		return session.Id()
	//	}
	return 0
}

//获取用户详细信息
func (this UserModule) GetUserInfo(userID uint64, session *link.Session) {
	onlineUser := module.Cache.GetOnlineUserByUserID(userID)
	if onlineUser != nil {
		dbUser := redisProxy.GetDBUser(userID)
		if dbUser != nil {
			userModel := NewUserModel(dbUser)
			gameProxy.SendGetUserInfoResult(session, 0, userModel)
		} else {
			gameProxy.SendGetUserInfoResult(session, gameProto.User_Not_Exists, nil)
		}
	} else {
		gameProxy.SendGetUserInfoResult(session, gameProto.User_Login_Fail, nil)
	}
}

//开启处理用户下线
func (this UserModule) StartDealOfflineUser() {
	this.onDealOfflineUserTimer()
	timer.DoTimer(int64(DEAL_OFFLINEUSER_INTERVAL), this.onDealOfflineUserTimer)
}

//定时处理用户下线
func (this UserModule) onDealOfflineUserTimer() {
	debug.Start("DealOfflineUserTimer")
	defer debug.Stop("DealOfflineUserTimer")

	users := redisProxy.GetAllUserLastLoginTime()
	INFO("Deal Remove User Redis Cache Data Num：", len(users))

	nowTime := time.Now().Unix()
	delUsers := list.New()
	for userID, lastLoginTime := range users {
		if nowTime >= lastLoginTime + OFFLINEUSER_TIME {
			delUsers.PushBack(userID)
		}
	}

	for tmp := delUsers.Front(); tmp != nil; tmp = tmp.Next() {
		userID, _ := strconv.ParseUint(tmp.Value.(string), 10, 64)
		//超时并且不在线
		if module.Cache.GetOnlineUserByUserID(userID) == nil {
			redisProxy.RemoveDBUser(userID)
		}
	}
	INFO("Remove User Redis Cache Data Num：", delUsers.Len())
}
