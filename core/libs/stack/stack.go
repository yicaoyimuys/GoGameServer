package stack

import (
	"runtime"
	"strconv"

	"github.com/yicaoyimuys/GoGameServer/core/libs/logger"
	"go.uber.org/zap"
)

// PrintPanicStack 输出错误堆栈信息
func PrintPanicStack() {
	for i := 0; i < 10; i++ {
		funcName, file, line, ok := runtime.Caller(i)
		if ok {
			funcName := runtime.FuncForPC(funcName).Name()
			logger.Error("错误堆栈", zap.String("Caller", "frame"+strconv.Itoa(i)+": [func:"+funcName+", file:"+file+", line:"+strconv.Itoa(line)+"]"))
		}
	}
}

// GetCallStack 获取调用堆栈
func GetCallStack() []string {
	result := make([]string, 10)
	for i := 0; i < 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if ok {
			funcName := runtime.FuncForPC(pc).Name()
			result[i] = "frame" + strconv.Itoa(i) + ": [func:" + funcName + ", file:" + file + ", line:" + strconv.Itoa(line) + "]"
		}
	}
	return result
}

// TryError 捕获异常
func TryError() {
	if x := recover(); x != nil {
		logger.Error("Error", zap.Any("Recover", x))
		PrintPanicStack()
	}
}

// CheckError 检查Error
func CheckError(err error) {
	if err != nil {
		logger.Error("Fatal error: %v", zap.Error(err))
	}
}
