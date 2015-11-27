package tools

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

import (
	"tools/cfg"
)

type SignalHandler func(s os.Signal, arg interface{})

type SignalSet struct {
	m map[os.Signal]SignalHandler
}

func SignalSetNew() *SignalSet {
	ss := new(SignalSet)
	ss.m = make(map[os.Signal]SignalHandler)
	return ss
}

func (set *SignalSet) Register(s os.Signal, handler SignalHandler) {
	if _, found := set.m[s]; !found {
		set.m[s] = handler
	}
}

func (set *SignalSet) Handle(sig os.Signal, arg interface{}) (err error) {
	if _, found := set.m[sig]; found {
		set.m[sig](sig, arg)
		return nil
	} else {
		return fmt.Errorf("No handler available for signal %v", sig)
	}
}

var (
	stopServerFunc func()
)

//处理信号
func SignalProc(stopServerCallback func()) {
	//http://blog.csdn.net/trojanpizza/article/details/6656321

	stopServerFunc = stopServerCallback

	sigHandler := SignalSetNew()
	sigHandler.Register(syscall.SIGWINCH, sigHandlerFunc)
	sigHandler.Register(syscall.SIGPIPE, sigHandlerFunc)

	sigHandler.Register(syscall.SIGHUP, sigHandlerFunc)
	sigHandler.Register(syscall.SIGINT, sigHandlerFunc)
	sigHandler.Register(syscall.SIGTERM, sigHandlerFunc)

	sigChan := make(chan os.Signal, 10)
	signal.Notify(sigChan)

	for true {
		select {
		case sig := <-sigChan:
			err := sigHandler.Handle(sig, nil)
			if err != nil {
				INFO("unknown signal received:", sig)
				//				os.Exit(1)
			}
		}
	}
}

func sigHandlerFunc(s os.Signal, arg interface{}) {
	switch s {
	case syscall.SIGHUP:
		cfg.Reload()
		INFO("ReloadConfig")
	case syscall.SIGINT:
		stopServer()
	case syscall.SIGTERM:
		stopServer()
	case syscall.SIGWINCH:
	case syscall.SIGPIPE:
	}
}

func stopServer() {
	if stopServerFunc != nil {
		stopServerFunc()
	}
	INFO("StopServer")
	os.Exit(1)
}
