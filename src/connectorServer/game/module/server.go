package module

import (
	"connectorServer/game/cache"
	"connectorServer/global"
	"connectorServer/sessions"
	. "connectorServer/tools"
	"connectorServer/tools/timer"
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
