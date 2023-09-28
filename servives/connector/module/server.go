package module

import (
	"github.com/yicaoyimuys/GoGameServer/core"
	"github.com/yicaoyimuys/GoGameServer/core/consts"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/sessions"
	"github.com/yicaoyimuys/GoGameServer/core/libs/timer"
	"github.com/yicaoyimuys/GoGameServer/servives/connector/cache"
)

func StartServerTimer() {
	initServerLogTimer()
}

func initServerLogTimer() {
	//每隔20秒记录一次
	timer.DoTimer(20*1000, func() {
		onlineUsersNum := sessions.FrontSessionLen()
		ip := core.Service.Ip()
		port := core.Service.Port(consts.ServiceType_WebSocket)
		cache.SetServerInfo(ip, port, onlineUsersNum)
		INFO("当前在线用户数量:", onlineUsersNum)
	})
}
