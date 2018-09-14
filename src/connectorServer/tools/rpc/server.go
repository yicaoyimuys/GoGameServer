package rpc

import (
	. "connectorServer/tools"
	"connectorServer/tools/stack"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
)

func InitServer() (string, error) {
	listen, err := net.Listen("tcp", ":")
	if err != nil {
		return "", err
	}

	go func() {
		defer stack.PrintPanicStackError()
		defer listen.Close()

		for {
			conn, err := listen.Accept()
			if err != nil {
				ERR("listen.Accept(): ", err)
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
