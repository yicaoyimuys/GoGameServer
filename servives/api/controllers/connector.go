package controllers

import (
	"github.com/yicaoyimuys/GoGameServer/core/consts"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/consul"

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

	serviceName := ""

	typeStr := this.GetString("type")
	if typeStr == "Socket" {
		serviceName = packageServiceName(consts.ServiceType_Socket, consts.Service_Connector)
	} else if typeStr == "WebSocket" {
		serviceName = packageServiceName(consts.ServiceType_WebSocket, consts.Service_Connector)
	}

	services := consulClient.GetServices(serviceName)

	this.Data["json"] = services
	this.ServeJSON()
}
