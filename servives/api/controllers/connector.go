package controllers

import (
	"GoGameServer/core/consts/Service"
	"GoGameServer/core/consts/ServiceType"
	. "GoGameServer/core/libs"
	"GoGameServer/core/libs/consul"

	"github.com/astaxie/beego"
)

type ConnectorController struct {
	beego.Controller
}

func init() {

}

func packageServiceName(serviceType string, serviceName string) string {
	return "<" + serviceType + ">" + serviceName
}

func (this *ConnectorController) Get() {
	consulClient, err := consul.NewClient()
	CheckError(err)

	serviceName := packageServiceName(ServiceType.WS, Service.Connector)
	services := consulClient.GetServices(serviceName)

	this.Data["json"] = services
	this.ServeJSON()
}
