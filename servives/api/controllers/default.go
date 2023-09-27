package controllers

import (
	"github.com/yicaoyimuys/GoGameServer/core/libs/guid"
	"github.com/yicaoyimuys/GoGameServer/core/libs/random"
	"github.com/yicaoyimuys/GoGameServer/servives/public/mongoModels"

	"github.com/astaxie/beego"
	"github.com/spf13/cast"
)

type DefaultController struct {
	beego.Controller
}

var (
	g *guid.Guid
)

func init() {
	g = guid.NewGuid(1)
}

func (this *DefaultController) Get() {
	id := g.NewID()
	account := this.GetString("name")
	if len(account) == 0 {
		account = "ys_" + cast.ToString(id)
	}
	money := int32(random.RandIntRange(1000, 9999))

	dbUser := mongoModels.AddUser(id, account, money)
	if dbUser == nil {
		this.Ctx.WriteString("mongo insert fail")
	} else {
		money = int32(random.RandIntRange(100, 999))
		if !mongoModels.UpdateUserMoney(id, money) {
			this.Ctx.WriteString("mongo update fail")
		} else {
			dbUser = mongoModels.GetUser(id)
			if dbUser == nil {
				this.Ctx.WriteString("mongo select fail")
			} else {
				this.Data["json"] = dbUser
				this.ServeJSON()
			}
		}
	}
}
