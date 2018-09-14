package module

import (
	. "core/libs"
	"core/libs/timer"
	"core/sessions"
	"game/cache"
	"global"
	"runtime"
)

func StartServerTimer() {
	initServerLogTimer()
}

func initServerLogTimer() {
	//每隔20秒记录一次
	timer.DoTimer(20*1000, func() {
		onlineUsersNum := sessions.FrontSessionLen()
		localIp := GetLocalIp()
		cache.SetServerInfo(localIp, global.ServerPort, onlineUsersNum)
		INFO("在线用户数量:" + NumToString(onlineUsersNum) + "   GoroutineNum:" + NumToString(runtime.NumGoroutine()))
	})
}
