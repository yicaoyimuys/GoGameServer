package sessions

import (
	"GoGameServer/core/libs/timer"
	"sync"
	"time"
)

var (
	frontSessions     = make(map[uint64]*FrontSession)
	frontSessionMutex sync.Mutex
)

func AddFrontSession(session *FrontSession) {
	frontSessionMutex.Lock()
	defer frontSessionMutex.Unlock()

	frontSessions[session.ID()] = session
	session.AddCloseCallback(nil, "sessionService.RemoveSession", func() {
		RemoveFrontSession(session.ID())
	})
}

func RemoveFrontSession(key uint64) {
	frontSessionMutex.Lock()
	defer frontSessionMutex.Unlock()

	delete(frontSessions, key)
}

func GetFrontSession(key uint64) *FrontSession {
	frontSessionMutex.Lock()
	defer frontSessionMutex.Unlock()

	session, _ := frontSessions[key]
	return session
}

func FrontSessionLen() int {
	frontSessionMutex.Lock()
	defer frontSessionMutex.Unlock()

	return len(frontSessions)
}

func FetchFrontSession(callback func(*FrontSession)) {
	frontSessionMutex.Lock()
	defer frontSessionMutex.Unlock()

	for _, session := range frontSessions {
		callback(session)
	}
}

func FrontSessionOpenPing(overTimeSec int64) {
	//2秒钟检测一次
	timer.DoTimer(2*1000, func() {
		frontSessionMutex.Lock()
		nowTime := time.Now().Unix()
		closeSessions := []*FrontSession{}
		for _, session := range frontSessions {
			//超过15秒
			cha := nowTime - session.PingTime()
			if cha >= overTimeSec {
				closeSessions = append(closeSessions, session)
			}
		}
		frontSessionMutex.Unlock()

		//关闭Session
		for _, session := range closeSessions {
			session.Close()
		}
	})
}
