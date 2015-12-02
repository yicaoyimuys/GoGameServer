package user

import (
	"github.com/funny/link"
	"global"
	. "model"
	"module"
	"protos/gameProto"
	"proxys/dbProxy"
	"proxys/redisProxy"
	"proxys/transferProxy"
	"proxys/worldProxy"
	"time"
	. "tools"
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
		module.SendLoginResult(0, session)
	} else {
		//登录成功处理
		success := this.LoginSuccess(session, userName, userID)
		if success {
			//发送登录成功消息
			module.SendLoginResult(userID, session)
			//更新最后一次登录时间
			this.updateLastLoginTime(session, userID)
		} else {
			module.SendLoginResult(0, session)
		}
	}
}

//登录
func (this UserModule) Login(userName string, session *link.Session) {
	onlineUser := module.Cache.GetOnlineUserByUserName(userName)
	if onlineUser != nil {
		if onlineUser.Session.Id() != session.Id() {
			//当前在线，但是连接不同，其他客户端连接，需通知当前客户端下线
			module.SendOtherLogin(onlineUser.Session)
			//替换Session
			module.Cache.RemoveOnlineUser(onlineUser.Session.Id())
			//登录成功处理
			success := this.LoginSuccess(session, onlineUser.UserName, onlineUser.UserID)
			if success {
				//发送登录成功消息
				module.SendLoginResult(onlineUser.UserID, session)
				//更新最后一次登录时间
				this.updateLastLoginTime(session, onlineUser.UserID)
			} else {
				module.SendLoginResult(0, session)
			}
		}
	} else {
		cacheDbUser := redisProxy.GetDBUserByUserName(userName)
		if cacheDbUser != nil {
			this.UserLoginHandle(session, cacheDbUser.Name, cacheDbUser.ID)
		} else {
			dbProxy.UserLogin(session.Id(), userName)
		}
	}
}

//更新用户最后上线时间，更新内存和数据库
func (this UserModule) updateLastLoginTime(session *link.Session, userID uint64) {
	nowTime := time.Now().Unix()
	redisProxy.UpdateUserLastLoginTime(userID, nowTime)
	dbProxy.UpdateUserLastLoginTime(session.Id(), userID, nowTime)
}

//用户登录成功处理
func (this UserModule) LoginSuccess(session *link.Session, userName string, userID uint64) bool {
	cacheSuccess := module.Cache.AddOnlineUser(userName, userID, session)
	if cacheSuccess {
		session.AddCloseCallback(session, func() {
			module.Cache.RemoveOnlineUser(session.Id())
			DEBUG("下线：在线人数", module.Cache.GetOnlineUsersNum())
		})
		DEBUG("上线：在线人数", module.Cache.GetOnlineUsersNum())

		//如果当前是LoginServer
		if global.IsLoginServer() {
			//通知GameServer登录成功
			transferProxy.SetClientLoginSuccess(userName, userID, session.Id())
		}

		//如果当前是GameServer
		if global.IsGameServer() {
			//通知WroldServer登录成功
			worldProxy.SetClientLoginSuccess(userName, userID, session.Id())
		}

		return true
	} else {
		ERR("what????")
		return false
	}
}

//用户上线
func (this UserModule) Online(session *link.Session) {
	global.AddSession(session)
	//如果当前是GameServer
	if global.IsGameServer() {
		//通知WroldServer用户上线
		worldProxy.SetClientSessionOnline(session)
	}
}

//用户下线
func (this UserModule) Offline(session *link.Session) {
	session.Close()
	//如果当前是GameServer
	if global.IsGameServer() {
		//通知WroldServer用户下线
		worldProxy.SetClientSessionOffline(session.Id())
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
			module.SendGetUserInfoResult(0, userModel, session)
		} else {
			module.SendGetUserInfoResult(gameProto.User_Not_Exists, nil, session)
		}
	} else {
		module.SendGetUserInfoResult(gameProto.User_Login_Fail, nil, session)
	}
}
