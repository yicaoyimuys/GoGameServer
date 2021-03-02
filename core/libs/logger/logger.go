package logger

import (
	"GoGameServer/core/libs/system"
	"fmt"
	"strings"

	"github.com/astaxie/beego/logs"
)

var (
	prefix = "log"
	log    = logs.GetBeeLogger()
)

func init() {
	logs.Async(10000)
}

func SetLogDebug(debug bool) {
	if debug {
		logs.SetLevel(logs.LevelDebug)
	} else {
		logs.SetLevel(logs.LevelInfo)
	}
}

func SetLogFile(fileName string, log_output string) {
	prefix = fileName
	startLogger(fileName+".log", log_output)
}

func startLogger(path string, log_output string) {
	path = system.Root + "/logs/" + path

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

func getLogMsg(v ...interface{}) string {
	return "[" + prefix + "] " + strings.TrimRight(fmt.Sprintln(v...), "\n")
}

func Error(v ...interface{}) {
	log.Error(getLogMsg(v))
}

func Warn(v ...interface{}) {
	log.Warn(getLogMsg(v))
}

func Info(v ...interface{}) {
	log.Info(getLogMsg(v))
}

func Notice(v ...interface{}) {
	log.Notice(getLogMsg(v))
}

func Debug(v ...interface{}) {
	log.Debug(getLogMsg(v))
}
