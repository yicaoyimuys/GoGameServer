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
	. "tools"
	"proxys/logProxy"
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
		success := this.LoginSuccess(session, userName, userID, 0)
		if success {
			//登录成功后处理
			this.dealLoginSuccess(session, userName, userID)
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
			success := this.LoginSuccess(session, onlineUser.UserName, onlineUser.UserID, 0)
			if success {
				//登录成功后处理
				this.dealLoginSuccess(session, userName, onlineUser.UserID)
			} else {
				module.SendLoginResult(0, session)
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
	module.SendLoginResult(userID, session)
	//如果用户在下线列表中，则移除
	module.Cache.RemoveOfflineUser(userID)
	//用户下线时处理
	session.AddCloseCallback(session, func() {
		//记录用户下线时间
		module.Cache.AddOfflineUser(userID)
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
			//记录用户下线Log
			if global.IsWorldServer() {
				logProxy.UserOffLine(userID)
			}
		})
		DEBUG("用户上线：当前在线人数", module.Cache.GetOnlineUsersNum())
		return true
	} else {
		ERR("what????")
		return false
	}
}

//用户上线
func (this UserModule) Online(session *link.Session) {
	global.AddSession(session)
}

//用户下线
func (this UserModule) Offline(session *link.Session) {
	session.Close()
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
