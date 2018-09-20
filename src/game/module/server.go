package module

import (
	"core"
	. "core/libs"
	"core/libs/common"
	"core/libs/timer"
	"core/sessions"
	"game/cache"
	"runtime"
)

func StartServerTimer() {
	initServerLogTimer()
}

func initServerLogTimer() {
	//每隔20秒记录一次
	timer.DoTimer(20*1000, func() {
		onlineUsersNum := sessions.FrontSessionLen()
		localIp := common.GetLocalIp()
		cache.SetServerInfo(localIp, core.Service.Port(), onlineUsersNum)
		INFO("在线用户数量:" + NumToString(onlineUsersNum) + "   GoroutineNum:" + NumToString(runtime.NumGoroutine()))
	})
}
