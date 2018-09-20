package core

import (
	"core/libs/grpc/ipc"
	"core/libs/rpc"
)

type IService interface {
	Env() string
	Name() string
	ID() int
	Port() string
	GetIpcClient(serviceName string) *ipc.Client
	GetRpcClient(serviceName string) *rpc.Client
}

var (
	Service IService
)
