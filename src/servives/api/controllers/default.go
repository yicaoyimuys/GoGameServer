package controllers

import (
	"github.com/astaxie/beego"
)

type DefaultController struct {
	beego.Controller
}

func (this *DefaultController) Get() {
	this.Ctx.WriteString("GoGameServer: hello world")
}
