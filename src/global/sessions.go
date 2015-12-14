package global

import (
	"github.com/funny/link"
	"sync"
)

var (
	sessions 		map[uint64]*link.Session = make(map[uint64]*link.Session)
	sessionMutex 	sync.Mutex
	sessionNum 		uint32 = 0
)

func AddSession(session *link.Session) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	sessionNum += 1
	sessions[session.Id()] = session
	session.AddCloseCallback(session, func() {
		RemoveSession(session.Id())
		sessionNum -= 1
	})
}

func RemoveSession(key uint64) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	if _, exists := sessions[key]; exists {
		delete(sessions, key)
	}
}

func GetSession(key uint64) *link.Session {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	if session, exists := sessions[key]; exists {
		return session
	}
	return nil
}
