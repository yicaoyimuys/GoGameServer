package gc

import (
	"runtime"
)

import (
	. "tools"
	"tools/timer"
)

const (
	GC_INTERVAL = 60 * 8
)

//系统routine
func SysRoutine() {
	timer.DoTimer(int64(GC_INTERVAL), onTimer)
}

func onTimer() {
	runtime.GC()
	INFO("GC executed")
	INFO("NumGoroutine", runtime.NumGoroutine())
	INFO("GC Summary:", GCSummary())
}
