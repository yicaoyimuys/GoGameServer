package user

import (
	"github.com/funny/link"
	. "model"
	"module"
	"protos"
	. "tools"
)

type UserModule struct {
}

// 在初始化的时候将模块注册到module包
func init() {
	module.User = UserModule{}
}

//登录
func (this UserModule) Login(userName string, session *link.Session) int32 {
	var result int32 = -1
	onlineUser, isOnline := module.Cache.UserIsOnline(userName)
	if isOnline {
		if onlineUser.Session.Id() != session.Id() {
			//当前在线，但是连接不同，其他客户端连接，需通知当前客户端下线
			simple.SendOtherLogin(onlineUser.Session)
			//替换Session
			module.Cache.RemoveOnlineUser(onlineUser.Session.Id())
			module.Cache.AddOnlineUser(onlineUser.UserName, onlineUser.UserID, session)
		}
		result = onlineUser.UserID
	} else {
		newUser := NewUserModel()
		if err := newUser.DBUser.GetUserByUserName(userName); err == nil {
			cacheSuccess := module.Cache.AddOnlineUser(newUser.DBUser.Name, newUser.DBUser.ID, session)
			if cacheSuccess {
				result = newUser.DBUser.ID
			}
		} else {
			DEBUG(err)
		}
	}

	if result != -1 {
		//Session断线处理
		session.AddCloseCallback(session, func() {
			module.Cache.RemoveOnlineUser(session.Id())
			DEBUG("下线：在线人数", module.Cache.GetOnlineUsersNum())
		})
	}

	DEBUG("上线：在线人数", module.Cache.GetOnlineUsersNum())
	return result
}

//重新连接
func (this UserModule) AgainConnect(oldSessionID uint64, session *link.Session) uint64 {
	if oldSessionID == session.Id() {
		return 0
	}

	user := module.Cache.GetUserBySession(oldSessionID)
	if user == nil {
		return 0
	}

	this.setOnline(user, session)
	return session.Id()
}

//获取用户详细信息
func (this UserModule) GetUserInfo(userID int32, session *link.Session) *UserModel {
	result := module.Cache.GetUser(userID)
	if result == nil {
		//不存在缓存用户
		newUser := NewUserModel()
		if err := newUser.DBUser.GetUser(userID); err == nil {
			result = newUser
			this.setOnline(result, session)
		} else {
			DEBUG(err)
		}
	} else {
		//		INFO("是否在线：", result.IsOnline)
		this.setOnline(result, session)
	}
	return result
}

//设置用户在线
func (this UserModule) setOnline(user *UserModel, session *link.Session) {
	module.Cache.RemoveUser(user)

	user.Session = session
	user.IsOnline = 1
	module.Cache.AddUser(user)

	session.AddCloseCallback(this, func() {
		this.setOffline(session)
	})
}

//设置用户下线
func (this UserModule) setOffline(session *link.Session) {
	u := module.Cache.GetUserBySession(session.Id())
	if u != nil {
		u.IsOnline = 0
	}
}
