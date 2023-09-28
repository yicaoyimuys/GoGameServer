package service

import (
	"github.com/astaxie/beego"
	"github.com/spf13/cast"
	"github.com/yicaoyimuys/GoGameServer/core/config"
	"github.com/yicaoyimuys/GoGameServer/core/consts"
)

func (this *Service) StartHttpServer() {
	//Api服务配置
	serviceConfig := config.GetService("api")
	serviceNodeConfig := serviceConfig.ServiceNodes[this.id]
	port := serviceNodeConfig.ClientPort
	useSSL := serviceNodeConfig.UseSSL

	//Http服务配置
	if useSSL {
		tslCrt := serviceConfig.TslCrt
		tslKey := serviceConfig.TslKey

		beego.BConfig.Listen.EnableHTTPS = true
		beego.BConfig.Listen.HTTPSCertFile = tslCrt
		beego.BConfig.Listen.HTTPSKeyFile = tslKey
		beego.BConfig.Listen.HTTPSPort = cast.ToInt(port)
	} else {
		beego.BConfig.Listen.HTTPPort = cast.ToInt(port)
	}
	beego.BConfig.RunMode = beego.PROD

	//启动http服务
	go beego.Run()

	//服务注册
	this.registerService(consts.ServiceType_Http, port)
}

func (this *Service) RegisterHttpRouter(rootPath string, controller beego.ControllerInterface) {
	beego.Router(rootPath, controller)
}
