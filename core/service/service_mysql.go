package service

import (
	"github.com/yicaoyimuys/GoGameServer/core/config"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/mysql"
)

func (this *Service) StartMysql() {
	this.mysqlClients = make(map[string]*mysql.Client)

	mysqlConfigs := config.GetMysqlConfig()
	index := 0
	for key, mysqlConfig := range mysqlConfigs {
		dbAliasName := key
		if index == 0 {
			dbAliasName = "default"
		}
		index++

		client, err := mysql.NewClient(dbAliasName, mysqlConfig)
		CheckError(err)

		if client != nil {
			this.mysqlClients[key] = client
			INFO("mysql_" + key + "连接成功...")
		}
	}
}

func (this *Service) GetMysqlClient(dbAliasName string) *mysql.Client {
	client, _ := this.mysqlClients[dbAliasName]
	return client
}
