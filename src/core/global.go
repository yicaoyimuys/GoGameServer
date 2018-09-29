package core

import (
	"core/libs/grpc/ipc"
	"core/libs/mysql"
	"core/libs/redis"
	"core/libs/rpc"
)

type IService interface {
	Env() string
	Name() string
	ID() int
	Identify() string
	GetIpcClient(serviceName string) *ipc.Client
	GetRpcClient(serviceName string) *rpc.Client
	GetRedisClient(redisAliasName string) *redis.Client
	GetMysqlClient(dbAliasName string) *mysql.Client
	GetIpcServer() *ipc.Server
	Ip() string
	Port(serviceType string) string
}

var (
	Service IService
)
