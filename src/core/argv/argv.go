package argv

import (
	. "core/libs"
	"github.com/jessevdk/go-flags"
)

var Values struct {
	Env         string `short:"e" long:"env" description:"环境" default:"local"`
	ServiceId   int    `short:"s" long:"serviceId" description:"服务ID" default:"1"`
	ServiceName string `short:"n" long:"serviceName" description:"服务名称" default:""`
}

func Init(serviceName string) error {
	_, err := flags.Parse(&Values)
	Values.ServiceName = serviceName
	DEBUG("启动参数：", Values)
	return err
}
