package module

import (
	"core"
	. "core/libs"
	"core/libs/sessions"
	"core/libs/timer"
	"runtime"
	"servives/connector/cache"
)

func StartServerTimer() {
	initServerLogTimer()
}

func initServerLogTimer() {
	//每隔20秒记录一次
	timer.DoTimer(20*1000, func() {
		onlineUsersNum := sessions.FrontSessionLen()
		localIp := GetLocalIp()
		cache.SetServerInfo(localIp, core.Service.Port(), onlineUsersNum)
		INFO("在线用户数量:" + NumToString(onlineUsersNum) + "   GoroutineNum:" + NumToString(runtime.NumGoroutine()))
	})
}
