package caches

import (
	"github.com/funny/link"
	. "model"
	"module"
	"sync"
)

type CacheModule struct {
	onlineUsers        map[string]*OnlineUserModel
	onlineUsersID      map[uint64]string
	onlineUsersSession map[uint64]string
	onlineUsersNum     int32
	onlineUsersMutex   sync.RWMutex
}

// 在初始化的时候将模块注册到module包
func init() {
	module.Cache = &CacheModule{
		onlineUsers:        make(map[string]*OnlineUserModel),
		onlineUsersID:      make(map[uint64]string),
		onlineUsersSession: make(map[uint64]string),
		onlineUsersNum:     0,
	}
}

//添加在线用户缓存
func (this *CacheModule) AddOnlineUser(userName string, userID uint64, session *link.Session, gameServerID uint32) bool {
	this.onlineUsersMutex.Lock()
	defer this.onlineUsersMutex.Unlock()

	//同一账号只能对应一个Session
	//同一SessionID只能登陆一个账号
	_, exists1 := this.onlineUsers[userName]
	_, exists2 := this.onlineUsersSession[session.Id()]
	if !exists1 && !exists2 {
		model := &OnlineUserModel{
			Session:        session,
			UserID:        userID,
			UserName:        userName,
			GameServerID:    gameServerID,
		}
		this.onlineUsers[userName] = model
		this.onlineUsersID[userID] = userName
		this.onlineUsersSession[session.Id()] = userName
		this.onlineUsersNum += 1
		return true
	} else {
		return false
	}
}

//获取用户数据根据UserName
func (this *CacheModule) GetOnlineUserByUserName(userName string) *OnlineUserModel {
	this.onlineUsersMutex.Lock()
	defer this.onlineUsersMutex.Unlock()

	if user, exists1 := this.onlineUsers[userName]; exists1 {
		return user
	}
	return nil
}

//获取用户数据根据UserID
func (this *CacheModule) GetOnlineUserByUserID(userID uint64) *OnlineUserModel {
	this.onlineUsersMutex.Lock()
	defer this.onlineUsersMutex.Unlock()

	if userName, exists := this.onlineUsersID[userID]; exists {
		if user, exists1 := this.onlineUsers[userName]; exists1 {
			return user
		}
	}
	return nil
}

//获取用户数据根据SessionID
func (this *CacheModule) GetOnlineUserBySession(sessionID uint64) *OnlineUserModel {
	this.onlineUsersMutex.Lock()
	defer this.onlineUsersMutex.Unlock()

	if userName, exists := this.onlineUsersSession[sessionID]; exists {
		if user, exists1 := this.onlineUsers[userName]; exists1 {
			return user
		}
	}
	return nil
}

//移除一个在线用户
func (this *CacheModule) RemoveOnlineUser(sessionID uint64) {
	this.onlineUsersMutex.Lock()
	defer this.onlineUsersMutex.Unlock()

	if userName, exists := this.onlineUsersSession[sessionID]; exists {
		if onlineUser, exists1 := this.onlineUsers[userName]; exists1 {
			delete(this.onlineUsers, onlineUser.UserName)
			delete(this.onlineUsersID, onlineUser.UserID)
			delete(this.onlineUsersSession, onlineUser.Session.Id())
			this.onlineUsersNum -= 1
		}
	}
}

//获取在线用户数量
func (this *CacheModule) GetOnlineUsersNum() int32 {
	return this.onlineUsersNum
}
