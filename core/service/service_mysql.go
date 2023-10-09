package service

import (
	"github.com/yicaoyimuys/GoGameServer/core/config"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/mysql"
	"go.uber.org/zap"
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
			INFO("Mysql连接成功", zap.String("AliasName", key))
		}
	}
}

func (this *Service) GetMysqlClient(dbAliasName string) *mysql.Client {
	client, _ := this.mysqlClients[dbAliasName]
	return client
}
