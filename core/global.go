package core

import (
	"GoGameServer/core/libs/grpc/ipc"
	"GoGameServer/core/libs/mongo"
	"GoGameServer/core/libs/mysql"
	"GoGameServer/core/libs/redis"
	"GoGameServer/core/libs/rpc"
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
	GetMongoClient(dbAliasName string) *mongo.Client
	GetIpcServer() *ipc.Server
	Ip() string
	Port(serviceType string) string
}

var (
	Service IService
)
