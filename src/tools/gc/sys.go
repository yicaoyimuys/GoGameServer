package gc

import (
	"runtime"
	"time"
)

import (
	. "tools"
	"tools/timer"
)

const (
	GC_INTERVAL = 300
)

//系统routine
func SysRoutine() {
	// timer
	gc_timer := make(chan int32, 10)
	gc_timer <- 1

	for {
		select {
		case <-gc_timer:
			// gc work
			runtime.GC()
			INFO("GC executed")
			INFO("NumGoroutine", runtime.NumGoroutine())
			INFO("GC Summary:", GCSummary())
			timer.Add(0, time.Now().Unix()+int64(GC_INTERVAL), gc_timer)
		}
	}
}
