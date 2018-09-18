package core

import (
	"core/argv"
	"core/config"
	. "core/libs"
	"core/libs/dict"
	"core/libs/logger"
	"runtime"
)

func NewService(serviceName string) {
	//错误捕获
	recoverErr()

	//初始化: 使用CPU数设置
	initMaxProcs()

	//初始化: 命令行参数
	initArgv(serviceName)

	//初始化: 配置文件
	initConfig()

	//初始化: log
	initLog()

	//系统环境输出
	printEnv()
}

func initMaxProcs() {
	//runtime.GOMAXPROCS(1)
}

func initArgv(serviceName string) {
	err := argv.Init(serviceName)
	checkError(err)
}

func initConfig() {
	config.Init()
}

func initLog() {
	logConfig := config.GetLog()

	logOpenDebug := dict.GetBool(logConfig, "debug")
	logOutput := dict.GetString(logConfig, "output")
	logFileName := argv.Values.ServiceName + "-" + NumToString(argv.Values.ServiceId)

	logger.SetLogFile(logFileName, logOutput)
	logger.SetLogDebug(logOpenDebug)
}

func printEnv() {
	INFO("使用CPU数量:" + NumToString(runtime.GOMAXPROCS(-1)))
	INFO("初始GoroutineNum:" + NumToString(runtime.NumGoroutine()))
	INFO("服务平台:" + argv.Values.Env)
	INFO("服务名称:" + argv.Values.ServiceName)
	INFO("服务ID:" + NumToString(argv.Values.ServiceId))
}

func recoverErr() {
	defer func() {
		if x := recover(); x != nil {
			ERR("caught panic in main()", x)
		}
	}()
}

func checkError(err error) {
	if err != nil {
		ERR("Fatal error: %v", err)
	}
}
