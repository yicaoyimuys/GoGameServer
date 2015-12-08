package debug

import (
	"tools/cfg"
	"time"
	. "tools"
)

var (
	debug_open bool
	logs map[string]int64
)

func init() {
	if cfg.GetValue("debug") == "true" {
		debug_open = true
	}

	logs = make(map[string]int64)
}

func Start(key string) {
	if !debug_open {
		return
	}
	logs[key] = time.Now().Unix()
}

func Stop(key string) int64 {
	var chaTime int64 = 0
	if !debug_open {
		return chaTime
	}

	if startTime, exists := logs[key]; exists {
		chaTime = time.Now().Unix() - startTime
		DEBUG("Debug Time [" + key + "]：", chaTime)
	}

	return chaTime
}