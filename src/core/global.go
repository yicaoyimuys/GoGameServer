package core

import (
	"core/libs/grpc/ipc"
)

type IService interface {
	Env() string
	Name() string
	ID() int
	Port() string
	GetIpcClient(serviceName string) *ipc.Client
}

var (
	Service IService
)
