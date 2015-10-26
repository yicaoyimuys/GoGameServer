package model

import (
	"errors"
	"strconv"
	. "tools"
	"tools/db"
)

type DBUserModel struct {
	db    *db.Model
	ID    int32
	Name  string
	Money int32
}

func NewDBUser() *DBUserModel {
	dbUser := new(DBUserModel)
	dbUser.db = db.DBOrm
	return dbUser
}

func (this *DBUserModel) GetUserByUserName(userName string) error {
	data := this.db.SetTable("user").Where("name = '" + userName + "'").FindOne()
	if data == nil {
		return errors.New("sql user select fail")
	}

	if len(data) == 0 {
		addMoney := RandomInt31n(999)

		var value = make(map[string]interface{}) //设置map存储数据，map[key]value
		value["name"] = userName
		value["money"] = int(addMoney)

		id, err := this.db.Insert(value) // 插入数据，返回增加个数和错误信息。返回最后增长的id和错误信息
		if err != nil {
			return err
		}

		this.ID = int32(id)
		this.Name = userName
		this.Money = addMoney
	} else {
		//		db.Print(data)
		this.full(data[1])
	}
	return nil
}

func (this *DBUserModel) GetUser(userId int32) error {
	data := this.db.SetTable("user").Where("id = " + strconv.Itoa(int(userId))).FindOne()
	if data == nil {
		return errors.New("sql select fail")
	}

	if len(data) == 0 {
		return errors.New("user is not exists")
	}

	this.full(data[1])
	return nil
}

func (this *DBUserModel) full(data map[string]string) {
	id, _ := strconv.Atoi(data["id"])
	this.ID = int32(id)

	this.Name = data["name"]

	money, _ := strconv.Atoi(data["money"])
	this.Money = int32(money)
}
