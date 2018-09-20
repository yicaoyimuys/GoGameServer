package argv

import (
	. "core/libs"
	"github.com/jessevdk/go-flags"
)

var Values struct {
	Env       string `short:"e" long:"env" description:"环境" default:"local"`
	ServiceId int    `short:"s" long:"serviceId" description:"服务ID" default:"1"`
}

func Init() error {
	_, err := flags.Parse(&Values)
	DEBUG("启动参数：", Values)
	return err
}
