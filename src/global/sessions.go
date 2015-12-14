package global

import (
	"github.com/funny/link"
	"sync"
)

var (
	sessions 		map[uint64]*link.Session = make(map[uint64]*link.Session)
	sessionMutex 	sync.Mutex
)

func AddSession(session *link.Session) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	sessions[session.Id()] = session
	session.AddCloseCallback(session, func() {
		RemoveSession(session.Id())
	})
}

func RemoveSession(key uint64) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	delete(sessions, key)
}

func GetSession(key uint64) *link.Session {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	session, _ := sessions[key]
	return session
}

func SessionLen() int {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	return len(sessions)
}

func FetchSession(callback func(*link.Session)) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	for _, session := range sessions {
		callback(session)
	}
}