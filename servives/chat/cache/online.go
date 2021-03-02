package cache

import (
	"GoGameServer/core/libs/sessions"
	"GoGameServer/servives/chat/model"
	"sync"
)

var (
	onlineUsers            = make(map[uint64]*model.ChatUser)
	onlineUsersNum   int32 = 0
	onlineUsersMutex sync.RWMutex
)

func AddUser(userID uint64, userName string, session *sessions.BackSession) {
	onlineUsersMutex.Lock()
	defer onlineUsersMutex.Unlock()

	user := &model.ChatUser{
		Session:  session,
		UserID:   userID,
		UserName: userName,
	}

	if _, ok := onlineUsers[userID]; !ok {
		onlineUsersNum++
	}
	onlineUsers[userID] = user
}

func RemoveUser(userID uint64) {
	onlineUsersMutex.Lock()
	defer onlineUsersMutex.Unlock()

	if _, ok := onlineUsers[userID]; ok {
		onlineUsersNum--
		delete(onlineUsers, userID)
	}
}

func GetUser(userID uint64) *model.ChatUser {
	onlineUsersMutex.Lock()
	defer onlineUsersMutex.Unlock()

	user, _ := onlineUsers[userID]
	return user
}

func GetOnlineUsersNum() int32 {
	return onlineUsersNum
}
