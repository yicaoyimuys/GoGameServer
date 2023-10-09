package libs

import (
	"github.com/yicaoyimuys/GoGameServer/core/libs/common"
	"github.com/yicaoyimuys/GoGameServer/core/libs/logger"
	"github.com/yicaoyimuys/GoGameServer/core/libs/stack"
	"github.com/yicaoyimuys/GoGameServer/core/libs/system"
	"go.uber.org/zap"
)

func init() {
}

func ERR(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func WARN(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func INFO(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func DEBUG(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Run() {
	system.Run()
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	return common.If(condition, trueVal, falseVal)
}

func CheckError(err error) {
	stack.CheckError(err)
}
