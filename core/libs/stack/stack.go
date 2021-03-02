package stack

import (
	"GoGameServer/core/libs/common"
	"GoGameServer/core/libs/logger"
	"runtime"
)

func PrintPanicStack() {
	for i := 0; i < 10; i++ {
		funcName, file, line, ok := runtime.Caller(i)
		if ok {
			funcName := runtime.FuncForPC(funcName).Name()
			logger.Error("frame " + common.NumToString(i) + ":[func:" + funcName + ", file: " + file + ", line:" + common.NumToString(line) + "]")
		}
	}
}

func TryError() {
	if x := recover(); x != nil {
		logger.Error(x)
		PrintPanicStack()
	}
}

func CheckError(err error) {
	if err != nil {
		logger.Error("Fatal error: %v", err)
	}
}
