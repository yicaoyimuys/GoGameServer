package sessions

import (
	"GoGameServer/core/libs/common"
	"sync"
)

var (
	backSessions     = make(map[string]*BackSession)
	backSessionMutex sync.Mutex
)

func SetBackSession(session *BackSession) {
	backSessionMutex.Lock()
	defer backSessionMutex.Unlock()

	if oldSession, ok := backSessions[session.id]; ok {
		oldSession.RemoveCloseCallback(nil, "RemoveBackSession")
		oldSession.Close()
	}

	backSessions[session.id] = session
	session.AddCloseCallback(nil, "RemoveBackSession", func() {
		RemoveBackSession(session.id)
	})
}

func GetBackSession(key string) *BackSession {
	backSessionMutex.Lock()
	defer backSessionMutex.Unlock()

	session, _ := backSessions[key]
	return session
}

func RemoveBackSession(key string) {
	backSessionMutex.Lock()
	defer backSessionMutex.Unlock()

	delete(backSessions, key)
}

func BackSessionLen() int {
	backSessionMutex.Lock()
	defer backSessionMutex.Unlock()

	return len(backSessions)
}

func CreateBackSessionId(serviceIdentify string, userSessionId uint64) string {
	return serviceIdentify + "_" + common.NumToString(userSessionId)
}
