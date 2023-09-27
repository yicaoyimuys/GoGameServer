package model

import "github.com/yicaoyimuys/GoGameServer/core/libs/sessions"

type ChatUser struct {
	Session  *sessions.BackSession
	UserID   uint64
	UserName string
}
