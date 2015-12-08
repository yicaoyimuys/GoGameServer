package tools

import (
	"fmt"
	"log"
	"strings"
)

import (
	"tools/cfg"
	"tools/logger"
)

var (
	debug_open bool
)

func init() {
	if cfg.GetValue("debug") == "true" {
		debug_open = true
	}
}

//------------------------------------------------ 严重程度由高到低
func ERR(v ...interface{}) {
	log.Printf("\033[1;4;31m[ERROR] %v \033[0m\n", strings.TrimRight(fmt.Sprintln(v...), "\n"))
}

func WARN(v ...interface{}) {
	log.Printf("\033[1;33m[WARN] %v \033[0m\n", strings.TrimRight(fmt.Sprintln(v...), "\n"))
}

func INFO(v ...interface{}) {
	log.Printf("\033[32m[INFO] %v \033[0m\n", strings.TrimRight(fmt.Sprintln(v...), "\n"))
}

func NOTICE(v ...interface{}) {
	log.Printf("[NOTICE] %v\n", strings.TrimRight(fmt.Sprintln(v...), "\n"))
}

func DEBUG(v ...interface{}) {
	if debug_open {
		log.Printf("\033[1;35m[DEBUG] %v \033[0m\n", strings.TrimRight(fmt.Sprintln(v...), "\n"))
	}
}

func SetLogFile(fileName string) {
	config := cfg.Get()
	if config[fileName] != "" {
		logger.StartLogger(config[fileName])
	}
}

func SetLogPrefix(prefix string) {
	log.SetPrefix("[" + prefix + "] ")
	//	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
