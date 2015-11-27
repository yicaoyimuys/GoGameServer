package module_db

import (
	"errors"
	. "model"
	"strconv"
	"time"
	. "tools"
	"tools/cfg"
	"tools/db"
	"tools/guid"
)

var userGuid *guid.Guid = guid.NewGuid()

//修改用户的最后登录时间
func UpdateUserLastLoginTime(userId uint64, time int64) error {
	var value = make(map[string]interface{})
	value["last_login_time"] = time
	_, err := db.DBOrm.SetTable("user").Where("id = " + strconv.FormatUint(userId, 10)).Update(value)
	return err
}

//根据用户名获取用户数据
func GetUserByUserName(userName string) (*DBUserModel, error) {
	data := db.DBOrm.SetTable("user").Where("name = '" + userName + "'").FindOne()
	if data == nil {
		return nil, errors.New("sql user select fail")
	}

	if len(data) != 0 {
		return fullDBUserModel(data[1]), nil
	} else {
		return insertUser(userName, userGuid.NewID(cfg.GetUint16("server_id")))
	}

	return nil, nil
}

//注册新用户
func insertUser(userName string, userId uint64) (*DBUserModel, error) {
	addMoney := RandomInt31n(999)

	nowTime := time.Now().Unix()
	var value = make(map[string]interface{}) //设置map存储数据，map[key]value
	value["id"] = userId
	value["name"] = userName
	value["money"] = int(addMoney)
	value["create_time"] = nowTime
	value["last_login_time"] = nowTime

	_, err := db.DBOrm.Insert(value) // 插入数据，返回增加个数和错误信息。返回最后增长的id和错误信息
	if err != nil {
		return nil, err
	}

	model := NewDBUserModel()
	model.ID = userId
	model.Name = userName
	model.Money = addMoney
	model.CreateTime = nowTime
	model.LastLoginTime = nowTime

	return model, nil
}

//填充DBUserModel
func fullDBUserModel(data map[string]string) *DBUserModel {
	model := NewDBUserModel()

	id, _ := strconv.ParseUint(data["id"], 10, 64)
	model.ID = id

	model.Name = data["name"]

	money, _ := strconv.Atoi(data["money"])
	model.Money = int32(money)

	create_time, _ := strconv.ParseInt(data["create_time"], 10, 64)
	model.CreateTime = create_time

	last_login_time, _ := strconv.ParseInt(data["last_login_time"], 10, 64)
	model.LastLoginTime = int64(last_login_time)

	return model
}
