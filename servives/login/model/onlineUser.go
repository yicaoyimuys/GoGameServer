package model

import "github.com/yicaoyimuys/GoGameServer/core/libs/sessions"

type OnlineUser struct {
	Session *sessions.BackSession
	UserID  uint64
	Account string
}
