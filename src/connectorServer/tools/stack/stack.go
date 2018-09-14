package stack

import (
	. "connectorServer/tools"
	"runtime"
)

func PrintPanicStack() {
	for i := 0; i < 10; i++ {
		funcName, file, line, ok := runtime.Caller(i)
		if ok {
			funcName := runtime.FuncForPC(funcName).Name()
			ERR("frame " + NumToString(i) + ":[func:" + funcName + ", file: " + file + ", line:" + NumToString(line) + "]")
		}
	}
}

func PrintPanicStackError() {
	if x := recover(); x != nil {
		ERR(x)
		PrintPanicStack()
	}
}
