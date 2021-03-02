package ipc

import (
	myGprc "GoGameServer/core/libs/grpc"
	"GoGameServer/core/libs/stack"
	"io"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
)

type ServerRecvHandle func(stream *Stream, msg *Req)
type StreamSession interface {
	Close()
}

type Stream struct {
	transferServer Ipc_TransferServer
	sessions       []StreamSession
	sessionsMutex  sync.Mutex
	closeFlag      int32
}

func (this *Stream) Send(userSessionIds []uint64, data []byte) error {
	msg := &Res{
		UserSessionIds: userSessionIds,
		Data:           data,
	}
	return this.transferServer.Send(msg)
}

func (this *Stream) IsClosed() bool {
	return atomic.LoadInt32(&this.closeFlag) == 1
}

func (this *Stream) AddSession(session StreamSession) {
	if this.IsClosed() {
		return
	}

	this.sessionsMutex.Lock()
	defer this.sessionsMutex.Unlock()

	this.sessions = append(this.sessions, session)
}

func (this *Stream) RemoveSession(session StreamSession) {
	if this.IsClosed() {
		return
	}

	this.sessionsMutex.Lock()
	defer this.sessionsMutex.Unlock()

	for index, s := range this.sessions {
		if s == session {
			this.sessions = append(this.sessions[:index], this.sessions[index+1:]...)
		}
	}
}

func (this *Stream) close() {
	if atomic.CompareAndSwapInt32(&this.closeFlag, 0, 1) {
		this.transferServer.Context().Done()

		this.sessionsMutex.Lock()
		defer this.sessionsMutex.Unlock()

		for _, session := range this.sessions {
			session.Close()
		}
		this.sessions = nil
	}
}

type Server struct {
	serverRecvHandle ServerRecvHandle
	streams          []*Stream
	streamMutex      sync.Mutex
}

func (this *Server) addStream(stream *Stream) {
	this.streamMutex.Lock()
	defer this.streamMutex.Unlock()

	this.streams = append(this.streams, stream)
}

func (this *Server) removeStream(stream *Stream) {
	this.streamMutex.Lock()
	defer this.streamMutex.Unlock()

	for index, s := range this.streams {
		if s == stream {
			this.streams = append(this.streams[:index], this.streams[index+1:]...)
		}
	}

	stream.close()
}

func (this *Server) SendToClient(stream *Stream, userSessionIds []uint64, data []byte) {
	stream.Send(userSessionIds, data)
}

func (this *Server) SendToAllClient(userSessionIds []uint64, data []byte) {
	this.streamMutex.Lock()
	defer this.streamMutex.Unlock()

	for _, stream := range this.streams {
		stream.Send(userSessionIds, data)
	}
}

func (this *Server) Transfer(stream Ipc_TransferServer) error {
	defer stack.TryError()

	s := &Stream{stream, []StreamSession{}, sync.Mutex{}, 0}
	this.addStream(s)

	defer this.removeStream(s)

	for {
		in, err := s.transferServer.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		go this.dealServerRecvHandle(s, in)
	}
}

func (this *Server) dealServerRecvHandle(stream *Stream, msg *Req) {
	defer stack.TryError()

	this.serverRecvHandle(stream, msg)
}

func InitServer(serverRecvHandle ServerRecvHandle) (*Server, string, error) {
	ipcServer := &Server{
		serverRecvHandle: serverRecvHandle,
		streams:          []*Stream{},
	}
	serverPort, err := myGprc.InitServer(func(grpcServer *grpc.Server) {
		//注册处理模块
		RegisterIpcServer(grpcServer, ipcServer)
	})
	return ipcServer, serverPort, err
}
