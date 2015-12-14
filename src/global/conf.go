package global

import (
	"strings"
	//	. "tools"
)

var (
	ServerName 	string
	ServerID 	uint32 = 0
)

func GetTrueServerName() string {
	return strings.Split(ServerName, "[")[0]
}

func LocalServer() bool {
	return GetTrueServerName() == "LocalServer"
}

func IsWorldServer() bool {
	return GetTrueServerName() == "WorldServer"
}

func IsGameServer() bool {
	return GetTrueServerName() == "GameServer"
}

func IsLoginServer() bool {
	return GetTrueServerName() == "LoginServer"
}
