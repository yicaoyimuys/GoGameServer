package global

import (
	//	"runtime"
	. "tools"
	. "tools/gc"
)

//服务器启动
func Startup(serverName string, logFile string, stopServerFunc func()) {
	//	runtime.GOMAXPROCS(runtime.NumCPU())

	// 开启Log记录
	SetLogFile(logFile)
	SetLogPrefix(serverName)

	// 信号量监听
	go SignalProc(stopServerFunc)

	// 开启GC及系统环境信息监测
	SysRoutine()

	// 开启服务器
	INFO("Starting " + serverName)
}

// 保持进程
func Run() {
	temp := make(chan int32, 10)
	for {
		select {
		case <-temp:
		}
	}
}
