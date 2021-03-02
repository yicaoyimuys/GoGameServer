package cache

import (
	"GoGameServer/core/libs/sessions"
	. "GoGameServer/servives/login/model"
	"sync"
)

var (
	onlineUsers              = make(map[string]*OnlineUser)
	onlineUserIds            = make(map[uint64]string)
	onlineUserSessions       = make(map[string]string)
	onlineUsersNum     int32 = 0
	onlineUsersMutex   sync.RWMutex
)

//添加在线用户缓存
func AddOnlineUser(userID uint64, account string, session *sessions.BackSession) bool {
	onlineUsersMutex.Lock()
	defer onlineUsersMutex.Unlock()

	//同一账号只能对应一个Session
	//同一SessionID只能登陆一个账号
	_, exists1 := onlineUsers[account]
	_, exists2 := onlineUserSessions[session.ID()]
	if !exists1 && !exists2 {
		model := &OnlineUser{
			Session: session,
			UserID:  userID,
			Account: account,
		}
		onlineUsers[account] = model
		onlineUserIds[userID] = account
		onlineUserSessions[session.ID()] = account
		onlineUsersNum += 1
		return true
	} else {
		return false
	}
}

//获取用户数据根据Account
func GetOnlineUserByAccount(account string) *OnlineUser {
	onlineUsersMutex.Lock()
	defer onlineUsersMutex.Unlock()

	if user, exists1 := onlineUsers[account]; exists1 {
		return user
	}
	return nil
}

//获取用户数据根据UserID
func GetOnlineUserByUserID(userID uint64) *OnlineUser {
	onlineUsersMutex.Lock()
	defer onlineUsersMutex.Unlock()

	if account, exists := onlineUserIds[userID]; exists {
		if user, exists1 := onlineUsers[account]; exists1 {
			return user
		}
	}
	return nil
}

//获取用户数据根据SessionID
func GetOnlineUserBySession(sessionID string) *OnlineUser {
	onlineUsersMutex.Lock()
	defer onlineUsersMutex.Unlock()

	if account, exists := onlineUserSessions[sessionID]; exists {
		if user, exists1 := onlineUsers[account]; exists1 {
			return user
		}
	}
	return nil
}

//移除一个在线用户
func RemoveOnlineUser(sessionID string) {
	onlineUsersMutex.Lock()
	defer onlineUsersMutex.Unlock()

	if account, exists := onlineUserSessions[sessionID]; exists {
		if onlineUser, exists1 := onlineUsers[account]; exists1 {
			delete(onlineUsers, onlineUser.Account)
			delete(onlineUserIds, onlineUser.UserID)
			delete(onlineUserSessions, onlineUser.Session.ID())
			onlineUsersNum -= 1
		}
	}
}

//获取在线用户数量
func GetOnlineUsersNum() int32 {
	return onlineUsersNum
}
