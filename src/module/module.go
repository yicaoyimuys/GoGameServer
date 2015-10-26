package module

import (
	"github.com/funny/link"
	. "model"
)

type CacheModule interface {
	AddOnlineUser(userName string, userID int32, session *link.Session) bool
	UserIsOnline(userName string) (OnlineUserModel, bool)
	RemoveOnlineUser(sessionID uint64)
	GetOnlineUsersNum() int32

	AddUser(user *UserModel)
	RemoveUser(user *UserModel)
	GetUser(userId int32) *UserModel
	GetUserByName(userName string) *UserModel
	GetUserBySession(sessionID uint64) *UserModel
}

type ConfigModule interface {
	Load()
}

type UserModule interface {
	Login(userName string, session *link.Session) int32
	AgainConnect(oldSessionID uint64, session *link.Session) uint64
	GetUserInfo(userID int32, session *link.Session) *UserModel
}

// 这些是接口的具体实现，等待外部主动注册进来，
// 这样module包永远是被引用的，不会出现递归引用问题。
var (
	Cache  CacheModule
	Config ConfigModule
	User   UserModule
)
