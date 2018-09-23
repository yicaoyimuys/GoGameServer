package dbModels

import (
	"core"
	"github.com/astaxie/beego/orm"
	"time"
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

//o := orm.NewOrm()
//
//user := User{Name: "slene"}
//
//// insert
//id, err := o.Insert(&user)
//
//// update
//user.Name = "astaxie"
//num, err := o.Update(&user)
//
//// read one
//u := User{Id: user.Id}
//err = o.Read(&u)
//
//// delete
//num, err = o.Delete(&u)

func AddUser(account string, money int32) *User {
	create_time := time.Now().Unix()

	userDb := core.Service.GetMysqlClient("user")
	user := User{
		Account:       account,
		Money:         money,
		CreateTime:    create_time,
		LastLoginTime: create_time,
	}

	// insert
	_, err := userDb.Insert(&user)
	if err != nil {
		return nil
	}
	return &user
}

func GetUser(account string) *User {
	userDb := core.Service.GetMysqlClient("user")
	user := User{Account: account}
	err := userDb.Read(&user, "Account")
	if err != nil {
		return nil
	}
	return &user
}
