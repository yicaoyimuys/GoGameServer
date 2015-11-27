package model

import (
	"github.com/funny/link"
)

type OnlineUserModel struct {
	Session  *link.Session
	UserID   uint64
	UserName string
}

type UserModel struct {
	DBUser *DBUserModel
}

func NewUserModel(dbUser *DBUserModel) *UserModel {
	return &UserModel{
		DBUser: dbUser,
	}
}
