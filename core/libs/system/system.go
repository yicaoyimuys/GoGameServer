package system

import (
	"os"

	"github.com/jessevdk/go-flags"
)

var (
	Root string
	Args struct {
		Env       string `short:"e" long:"env" description:"环境" default:"local"`
		ServiceId int    `short:"s" long:"serviceId" description:"服务ID" default:"1"`
	}
)

func init() {
	initRootPath()
	initArgs()
}

func initRootPath() {
	Root, _ = os.Getwd()
}

func initArgs() error {
	_, err := flags.Parse(&Args)
	return err
}

// Run 保持进程
func Run() {
	temp := make(chan int32, 10)
	for {
		select {
		case <-temp:
		}
	}
}
