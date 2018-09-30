package module

import (
	"core"
	"core/consts/serviceType"
	. "core/libs"
	"core/libs/sessions"
	"core/libs/timer"
	"servives/connector/cache"
)

func StartServerTimer() {
	initServerLogTimer()
}

func initServerLogTimer() {
	//每隔20秒记录一次
	timer.DoTimer(20*1000, func() {
		onlineUsersNum := sessions.FrontSessionLen()
		ip := core.Service.Ip()
		port := core.Service.Port(ServiceType.WS)
		cache.SetServerInfo(ip, port, onlineUsersNum)
		INFO("在线用户数量:" + NumToString(onlineUsersNum))
	})
}
