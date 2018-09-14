package global

import (
	"core/libs/grpc/ipc"
	"core/libs/guid"
)

type services struct {
	Matching string
	Game     string
}

type ipcClients struct {
	Matching *ipc.Client
	Game     *ipc.Client
}

var (
	Env        string
	Guid       *guid.Guid
	ServerPort string
	ServerId   int
	ServerName string

	GameServerComputer int

	Services = services{
		Matching: "MatchingServer",
		Game:     "GameServer",
	}

	IpcClients = ipcClients{}
)

func InitGuid(serverId uint16) {
	Guid = guid.NewGuid(serverId)
}
