package model

import (
	"github.com/funny/link"
)

type OnlineUserModel struct {
	Session  *link.Session
	UserID   int32
	UserName string
}

type UserModel struct {
	DBUser   *DBUserModel
	Session  *link.Session
	IsOnline int32
}

func NewUserModel() *UserModel {
	return &UserModel{
		DBUser:   NewDBUser(),
		IsOnline: 0,
	}
}
