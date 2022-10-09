package socket

import (
	"GoGameServer/core/libs/guid"
	"GoGameServer/core/libs/logger"
	"GoGameServer/core/libs/sessions"
	"GoGameServer/core/libs/stack"
	"net"
)

const (
	ServerNetworkType = "tcp4"
)

type Server struct {
	port string
	guid *guid.Guid

	sessionCreateHandle     sessions.FrontSessionCreateHandle
	sessionReceiveMsgHandle sessions.FrontSessionReceiveMsgHandle
}

func NewServer(port string, serviceId int) *Server {
	server := &Server{
		port: port,
		guid: guid.NewGuid(uint16(serviceId)),
	}
	return server
}

func (this *Server) SetSessionCreateHandle(handle sessions.FrontSessionCreateHandle) {
	this.sessionCreateHandle = handle
}

func (this *Server) SetSessionReceiveMsgHandle(handle sessions.FrontSessionReceiveMsgHandle) {
	this.sessionReceiveMsgHandle = handle
}

func (this *Server) Start() {
	logger.Info("front start socket...", this.port)

	go func() {
		defer stack.TryError()

		var err error
		addr, err := net.ResolveTCPAddr(ServerNetworkType, "0.0.0.0:"+this.port)
		stack.CheckError(err)

		listener, err := net.ListenTCP(ServerNetworkType, addr)
		stack.CheckError(err)

		defer listener.Close()
		logger.Info("socket waiting client connect...")
		for {
			conn, err := listener.Accept()
			stack.CheckError(err)

			go this.handleConnect(conn)
		}
	}()
}

func (this *Server) StartPing() {
	overTime := 15
	sessions.FrontSessionOpenPing(int64(overTime))
	logger.Info("session超时时间设置", overTime)
}

func (this *Server) handleConnect(conn net.Conn) {
	//捕获异常
	defer stack.TryError()

	//Session创建
	sessionId := this.guid.NewID()
	sessionCodec := NewFrontCodec(conn)
	session := sessions.NewFontSession(sessionId, sessionCodec)
	this.addFontSession(session)
}

func (this *Server) addFontSession(session *sessions.FrontSession) {
	sessions.AddFrontSession(session)
	if this.sessionCreateHandle != nil {
		this.sessionCreateHandle(session)
	}
	if this.sessionReceiveMsgHandle != nil {
		session.SetMsgHandle(this.sessionReceiveMsgHandle)
	}

	defer session.Close()
	for {
		msg, err := session.Receive()
		if err != nil || msg == nil {
			break
		}
	}
}
