package mysql

import (
	"github.com/yicaoyimuys/GoGameServer/core/config"
	"github.com/yicaoyimuys/GoGameServer/core/libs/logger"
	"go.uber.org/zap"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		logger.Error("Mysql_注册失败", zap.Error(err))
		return
	}
}

type Client struct {
	orm.Ormer
}

func NewClient(dbAliasName string, mysqlConfig config.MysqlConfig) (*Client, error) {
	dbHost := mysqlConfig.Host
	dbPort := mysqlConfig.Port
	dbUser := mysqlConfig.User
	dbPassword := mysqlConfig.Password
	dbName := mysqlConfig.Db
	dbCharset := mysqlConfig.Charset

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
