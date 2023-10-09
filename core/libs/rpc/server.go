package rpc

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"

	"github.com/yicaoyimuys/GoGameServer/core/libs/logger"
	"github.com/yicaoyimuys/GoGameServer/core/libs/stack"
	"go.uber.org/zap"
)

func InitServer() (string, error) {
	listen, err := net.Listen("tcp", ":")
	if err != nil {
		return "", err
	}

	go func() {
		defer stack.TryError()
		defer listen.Close()

		for {
			conn, err := listen.Accept()
			if err != nil {
				logger.Error("Listen.Accept()", zap.Error(err))
			}
			go jsonrpc.ServeConn(conn)
		}
	}()

	serverPort := strconv.Itoa(listen.Addr().(*net.TCPAddr).Port)
	return serverPort, nil
}

func RegisterModule(name string, rcvr interface{}) error {
	err := rpc.RegisterName(name, rcvr)
	return err
}
