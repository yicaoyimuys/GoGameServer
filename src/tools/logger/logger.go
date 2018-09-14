package logger

import (
	"github.com/astaxie/beego/logs"
	"tools/cfg"
)

func init() {
	logs.Async(10000)
}

//系统日志
func StartLogger(path string, log_output string) {
	path = cfg.ROOT + "/logs/" + path

	switch log_output {
	case "both":
		logs.SetLogger("console", "")
	case "file":
		logs.SetLogger("file", `{"filename":"`+path+`"}`)
	case "both&file":
		logs.SetLogger("console", "")
		logs.SetLogger("file", `{"filename":"`+path+`"}`)
	}
}
