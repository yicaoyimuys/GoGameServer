package global

import (
	"github.com/funny/link"
	"sync"
	//	. "tools"
)

var ServerName string

var sessions map[uint64]*link.Session = make(map[uint64]*link.Session)
var sessionMutex sync.Mutex
var sessionNum uint32 = 0

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
