package grpc

import (
	"GoGameServer/core/libs/stack"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

func InitServer(registerPbServiceFunc func(*grpc.Server)) (string, error) {
	//创建监听
	listen, err := net.Listen("tcp", ":")
	if err != nil {
		return "", err
	}

	go func() {
		defer stack.TryError()
		defer listen.Close()

		//创建grpcServer
		grpcServer := grpc.NewServer()

		//注册服务
		registerPbServiceFunc(grpcServer)

		//服务开启
		grpcServer.Serve(listen)
	}()

	//返回端口
	serverPort := strconv.Itoa(listen.Addr().(*net.TCPAddr).Port)
	return serverPort, nil
}
