package global

import (
	"github.com/funny/link"
	"strings"
	"sync"
	//	. "tools"
)

var ServerName string
var ServerID uint32 = 0

var sessions map[uint64]*link.Session = make(map[uint64]*link.Session)
var sessionMutex sync.Mutex
var sessionNum uint32 = 0

func GetTrueServerName() string {
	return strings.Split(ServerName, "[")[0]
}

func LocalServer() bool {
	return GetTrueServerName() == "LocalServer"
}

func IsWorldServer() bool {
	return GetTrueServerName() == "WorldServer"
}

func IsGameServer() bool {
	return GetTrueServerName() == "GameServer"
}

func IsLoginServer() bool {
	return GetTrueServerName() == "LoginServer"
}

func AddSession(session *link.Session) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	sessionNum += 1
	sessions[session.Id()] = session
	session.AddCloseCallback(session, func() {
		RemoveSession(session)
		sessionNum -= 1
	})
}

func RemoveSession(session *link.Session) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	delete(sessions, session.Id())
}

func GetSession(sessionId uint64) *link.Session {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	if session, exists := sessions[sessionId]; exists {
		return session
	}
	return nil
}
