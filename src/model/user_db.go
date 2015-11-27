package model

type DBUserModel struct {
	ID            uint64
	Name          string
	Money         int32
	CreateTime    int64
	LastLoginTime int64
}

func NewDBUserModel() *DBUserModel {
	dbUser := new(DBUserModel)
	return dbUser
}
