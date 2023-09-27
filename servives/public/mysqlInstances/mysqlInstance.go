package mysqlInstances

import (
	"github.com/yicaoyimuys/GoGameServer/core"
	"github.com/yicaoyimuys/GoGameServer/core/libs/mysql"
)

func Global() *mysql.Client {
	return core.Service.GetMysqlClient("global")
}

func User() *mysql.Client {
	return core.Service.GetMysqlClient("user")
}

func Log() *mysql.Client {
	return core.Service.GetMysqlClient("log")
}
