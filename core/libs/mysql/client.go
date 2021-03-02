package mysql

import (
	"GoGameServer/core/libs/dict"
	"GoGameServer/core/libs/logger"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		logger.Error("Mysql_注册失败", err)
		return
	}
	//开启debug调试
	orm.Debug = true
}

type Client struct {
	orm.Ormer
}

func NewClient(dbAliasName string, mysqlConfig map[string]interface{}) (*Client, error) {
	dbUser := dict.GetString(mysqlConfig, "user")
	dbPassword := dict.GetString(mysqlConfig, "password")
	dbHost := dict.GetString(mysqlConfig, "host")
	dbPort := dict.GetString(mysqlConfig, "port")
	dbName := dict.GetString(mysqlConfig, "db")
	dbCharset := dict.GetString(mysqlConfig, "charset")

	//数据库连接
	dataSource := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=" + dbCharset
	err := orm.RegisterDataBase(dbAliasName, "mysql", dataSource)
	if err != nil {
		return nil, err
	}

	//连接池设置
	orm.SetMaxIdleConns(dbAliasName, 30)
	orm.SetMaxOpenConns(dbAliasName, 30)

	//创建Orm对象
	o := orm.NewOrm()
	err = o.Using(dbAliasName)
	if err != nil {
		return nil, err
	}

	//返回数据
	client := &Client{o}
	return client, nil
}
