package caches

import (
	"github.com/funny/link"
	. "model"
	"module"
	"sync"
	//	. "tools"
)

type CacheModule struct {
	onlineUsers        map[string]OnlineUserModel
	onlineUsersSession map[uint64]string
	onlineUsersNum     int32
	onlineUsersMutex   sync.RWMutex

	users        map[int32]*UserModel
	userNames    map[string]int32
	userSessions map[uint64]int32
}

// 在初始化的时候将模块注册到module包
func init() {
	module.Cache = &CacheModule{
		onlineUsers:        make(map[string]OnlineUserModel),
		onlineUsersSession: make(map[uint64]string),
		onlineUsersNum:     0,
		users:              make(map[int32]*UserModel),
		userNames:          make(map[string]int32),
		userSessions:       make(map[uint64]int32),
	}
}

//添加在线用户缓存
func (this *CacheModule) AddOnlineUser(userName string, userID int32, session *link.Session) bool {
	this.onlineUsersMutex.Lock()
	defer this.onlineUsersMutex.Unlock()

	//同一账号只能对应一个Session
	//同一SessionID只能登陆一个账号
	_, exists1 := this.onlineUsers[userName]
	_, exists2 := this.onlineUsersSession[session.Id()]
	if !exists1 && !exists2 {
		model := OnlineUserModel{
			Session:  session,
			UserID:   userID,
			UserName: userName,
		}
		this.onlineUsers[userName] = model
		this.onlineUsersSession[session.Id()] = userName
		this.onlineUsersNum += 1
		return true
	} else {
		return false
	}
}

//判定一个用户是否在线
func (this *CacheModule) UserIsOnline(userName string) (OnlineUserModel, bool) {
	this.onlineUsersMutex.Lock()
	defer this.onlineUsersMutex.Unlock()

	value, exists := this.onlineUsers[userName]
	return value, exists
}

//移除一个在线用户
func (this *CacheModule) RemoveOnlineUser(sessionID uint64) {
	this.onlineUsersMutex.Lock()
	defer this.onlineUsersMutex.Unlock()

	if userName, exists := this.onlineUsersSession[sessionID]; exists {
		delete(this.onlineUsers, userName)
		delete(this.onlineUsersSession, sessionID)
		this.onlineUsersNum -= 1
	}
}

//获取在线用户数量
func (this *CacheModule) GetOnlineUsersNum() int32 {
	return this.onlineUsersNum
}

//缓存用户详细信息
func (this *CacheModule) AddUser(user *UserModel) {
	this.users[user.DBUser.ID] = user
	this.userNames[user.DBUser.Name] = user.DBUser.ID
	this.userSessions[user.Session.Id()] = user.DBUser.ID
}

//移除用户缓存用户
func (this *CacheModule) RemoveUser(user *UserModel) {
	if _, exists := this.users[user.DBUser.ID]; exists {
		delete(this.users, user.DBUser.ID)
		delete(this.userNames, user.DBUser.Name)
		delete(this.userSessions, user.Session.Id())
	}
}

//获取用户详细信息
func (this *CacheModule) GetUser(userId int32) *UserModel {
	if u, exists := this.users[userId]; exists {
		return u
	}
	return nil
}

//根据用户名获取用户信息
func (this *CacheModule) GetUserByName(userName string) *UserModel {
	if userID, exists := this.userNames[userName]; exists {
		return this.GetUser(userID)
	}
	return nil
}

//根据SessionID获取用户信息
func (this *CacheModule) GetUserBySession(sessionID uint64) *UserModel {
	if userID, exists := this.userSessions[sessionID]; exists {
		return this.GetUser(userID)
	}
	return nil
}
