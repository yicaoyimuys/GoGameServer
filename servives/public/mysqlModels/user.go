package mysqlModels

import (
	"GoGameServer/servives/public/mysqlInstances"
	"time"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id            uint64
	Account       string
	Money         int32
	CreateTime    int64
	LastLoginTime int64
}

func init() {
	orm.RegisterModel(new(User))
}

func AddUser(account string, money int32) *User {
	create_time := time.Now().Unix()

	user := User{
		Account:       account,
		Money:         money,
		CreateTime:    create_time,
		LastLoginTime: create_time,
	}

	// insert
	_, err := mysqlInstances.User().Insert(&user)
	if err != nil {
		return nil
	}
	return &user
}

func GetUser(account string) *User {
	user := User{Account: account}
	err := mysqlInstances.User().Read(&user, "Account")
	if err != nil {
		return nil
	}
	return &user
}

func UpdateUser(dbUser *User) bool {
	_, err := mysqlInstances.User().Update(dbUser)
	if err != nil {
		return false
	}
	return true
}

func UpdateUserLoginTime(userId uint64, loginTime int64) bool {
	dbUser := User{
		Id:            userId,
		LastLoginTime: loginTime,
	}
	_, err := mysqlInstances.User().Update(&dbUser, "LastLoginTime")
	if err != nil {
		return false
	}
	return true
}
