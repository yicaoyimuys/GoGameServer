package caches

import (
	"github.com/funny/link"
	. "model"
	"module"
	"sync"
	. "tools"
	"time"
	"tools/timer"
	"container/list"
	"proxys/redisProxy"
	"tools/debug"
)

type CacheModule struct {
	onlineUsers        map[string]*OnlineUserModel
	onlineUsersID      map[uint64]string
	onlineUsersSession map[uint64]string
	onlineUsersNum     int32
	onlineUsersMutex   sync.RWMutex
	offlineUsers       map[uint64]int64
	offlineUsersNum    int32
	offlineUsersMutex  sync.RWMutex
}

const (
	//处理下线用户数据间隔(5分钟)
	DEAL_OFFLINEUSER_INTERVAL = 5 * 60
	//用户下线后数据存在时长(3小时)
	OFFLINEUSER_TIME = 3 * 60 * 60
)

// 在初始化的时候将模块注册到module包
func init() {
	module.Cache = &CacheModule{
		onlineUsers:        make(map[string]*OnlineUserModel),
		onlineUsersID:      make(map[uint64]string),
		onlineUsersSession: make(map[uint64]string),
		onlineUsersNum:     0,
		offlineUsers:       make(map[uint64]int64),
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

//添加下线用户
func (this *CacheModule) AddOfflineUser(userID uint64) {
	this.offlineUsersMutex.Lock()
	defer this.offlineUsersMutex.Unlock()

	this.offlineUsers[userID] = time.Now().Unix()
	this.offlineUsersNum += 1
}

//移除下线用户
func (this *CacheModule) RemoveOfflineUser(userID uint64) {
	this.offlineUsersMutex.Lock()
	defer this.offlineUsersMutex.Unlock()

	if _, exists := this.offlineUsers[userID]; exists {
		delete(this.offlineUsers, userID)
		this.offlineUsersNum -= 1
	}
}

//开启处理用户下线
func (this *CacheModule) StartDealOfflineUser() {
	timer.DoTimer(int64(DEAL_OFFLINEUSER_INTERVAL), this.onDealOfflineUserTimer)
}

//定时处理用户下线
func (this *CacheModule) onDealOfflineUserTimer() {
	debug.Start("DealOfflineUserTimer")
	defer debug.Stop("DealOfflineUserTimer")

	this.offlineUsersMutex.Lock()
	defer this.offlineUsersMutex.Unlock()

	nowTime := time.Now().Unix()
	delUsers := list.New()
	for userID, offlineTime := range this.offlineUsers {
		if nowTime >= offlineTime + OFFLINEUSER_TIME {
			delUsers.PushBack(userID)
		}
	}

	for tmp := delUsers.Front(); tmp != nil; tmp = tmp.Next() {
		userID := tmp.Value.(uint64)
		delete(this.offlineUsers, userID)
		this.offlineUsersNum -= 1
		redisProxy.RemoveDBUser(userID)
	}
	INFO("Remove User Redis Cache Data Num：", delUsers.Len())
}
