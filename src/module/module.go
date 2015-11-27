package module

import (
	"github.com/funny/link"
	. "model"
)

type CacheModule interface {
	AddOnlineUser(userName string, userID uint64, session *link.Session) bool
	GetOnlineUserByUserName(userName string) *OnlineUserModel
	GetOnlineUserByUserID(userID uint64) *OnlineUserModel
	GetOnlineUserBySession(sessionID uint64) *OnlineUserModel
	RemoveOnlineUser(sessionID uint64)
	GetOnlineUsersNum() int32
}

type ConfigModule interface {
	Load()
}

type UserModule interface {
	UserLoginHandle(session *link.Session, userName string, userID uint64)

	Login(userName string, session *link.Session)
	AgainConnect(oldSessionID uint64, session *link.Session) uint64
	GetUserInfo(userID uint64, session *link.Session)
}

// 这些是接口的具体实现，等待外部主动注册进来，
// 这样module包永远是被引用的，不会出现递归引用问题。
var (
	Cache  CacheModule
	Config ConfigModule
	User   UserModule
)
