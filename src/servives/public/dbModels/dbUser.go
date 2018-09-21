package dbModels

import (
	"core"
	. "core/libs"
	"time"
)

type DbUser struct {
	Id            uint64
	Account       string
	Money         int32
	CreateTime    int64
	LastLoginTime int64
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

func AddDbUser(account string, money int32) *DbUser {
	create_time := time.Now().Unix()

	userDb := core.Service.GetMysqlClient("user")
	user := DbUser{
		Account:       account,
		Money:         money,
		CreateTime:    create_time,
		LastLoginTime: create_time,
	}

	// insert
	id, err := userDb.Insert(&user)
	if err != nil {
		return nil
	}
	DEBUG(user, id)
	return &user
}

func GetDbUser(account string) *DbUser {
	userDb := core.Service.GetMysqlClient("user")
	user := DbUser{Account: account}
	err := userDb.Read(&user)
	if err != nil {
		return nil
	}
	return &user
}
