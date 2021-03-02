package libs

import (
	"GoGameServer/core/libs/common"
	"GoGameServer/core/libs/logger"
	"GoGameServer/core/libs/stack"
	"GoGameServer/core/libs/system"
)

func init() {
}

func ERR(v ...interface{}) {
	logger.Error(v)
}

func WARN(v ...interface{}) {
	logger.Warn(v)
}

func INFO(v ...interface{}) {
	logger.Info(v)
}

func NOTICE(v ...interface{}) {
	logger.Notice(v)
}

func DEBUG(v ...interface{}) {
	logger.Debug(v)
}

func Run() {
	system.Run()
}

func NumToString(num interface{}) string {
	return common.NumToString(num)
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	return common.If(condition, trueVal, falseVal)
}

func CheckError(err error) {
	stack.CheckError(err)
}
