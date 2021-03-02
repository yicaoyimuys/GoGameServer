package mongoModels

import (
	"GoGameServer/servives/public/mongoInstances"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id         uint64 `bson:"_id" json:"id"`
	Account    string `bson:"account" json:"account"`
	Money      int32  `bson:"money" json:"money"`
	CreateTime int64  `bson:"create_time" json:"create_time"`
}

var collection string

func init() {
	collection = "users"
}

func AddUser(id uint64, account string, money int32) *User {
	create_time := time.Now().Unix()

	user := User{
		Id:         id,
		Account:    account,
		Money:      money,
		CreateTime: create_time,
	}

	// insert
	err := mongoInstances.Global().Insert(collection, user)
	if err != nil {
		return nil
	}
	return &user
}

func GetUser(id uint64) *User {
	var user User
	err := mongoInstances.Global().FindOne(collection, bson.M{"_id": id}, nil, &user)
	if err != nil {
		return nil
	}
	return &user
}

func UpdateUserMoney(id uint64, money int32) bool {
	err := mongoInstances.Global().Update(collection, bson.M{"_id": id}, bson.M{"$set": bson.M{"money": money}})
	if err != nil {
		return false
	}
	return true
}
