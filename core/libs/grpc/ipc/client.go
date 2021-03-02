package ipc

import (
	"GoGameServer/core/libs/consul"
	myGprc "GoGameServer/core/libs/grpc"
	"GoGameServer/core/libs/stack"
	"errors"
	"sync"

	"google.golang.org/grpc"
)

type ClientRecvHandle func(stream Ipc_TransferClient, msg *Res)

type Client struct {
	grpcClient        *myGprc.Client
	recvHandle        ClientRecvHandle
	serverStreams     map[string]Ipc_TransferClient
	serverStreamMutex sync.Mutex
}

func NewClient(consulClient *consul.Client, serviceName string, handle ClientRecvHandle) *Client {
	grpcClient := myGprc.NewClient(consulClient, serviceName, func(conn *grpc.ClientConn) interface{} {
		return NewIpcClient(conn)
	})

	client := &Client{
		grpcClient:    grpcClient,
		recvHandle:    handle,
		serverStreams: make(map[string]Ipc_TransferClient),
	}
	return client
}

func (this *Client) dealRecvHandle(stream Ipc_TransferClient, msg *Res) {
	defer stack.TryError()

	this.recvHandle(stream, msg)
}

func (this *Client) loop(service string, stream Ipc_TransferClient) {
	defer stack.TryError()
	defer this.removeStream(service)

	for {
		in, err := stream.Recv()
		if err != nil {
			return
		}
		go this.dealRecvHandle(stream, in)
	}
}

func (this *Client) getStream(service string) Ipc_TransferClient {
	this.serverStreamMutex.Lock()
	defer this.serverStreamMutex.Unlock()

	//检测是否已经存在
	stream, ok := this.serverStreams[service]
	if ok {
		return stream
	}

	//创建新的stream
	transferClient := this.grpcClient.Call(service, "Transfer", nil)
	if transferClient == nil {
		return nil
	}
	stream = transferClient.(Ipc_TransferClient)
	this.serverStreams[service] = stream
	go this.loop(service, stream)

	return stream
}

func (this *Client) removeStream(service string) {
	this.serverStreamMutex.Lock()
	if stream, ok := this.serverStreams[service]; ok {
		stream.Context().Done()
		stream.CloseSend()
		delete(this.serverStreams, service)
	}
	this.serverStreamMutex.Unlock()
}

func (this *Client) GetServiceByRandom() string {
	return this.grpcClient.GetServiceByRandom()
}

func (this *Client) GetServiceByFlag(flag string) string {
	return this.grpcClient.GetServiceByFlag(flag)
}

func (this *Client) Send(senderServiceIdentify string, userSessionId uint64, data []byte, receiverService string) error {
	if receiverService == "" {
		return errors.New("service is null")
	}

	stream := this.getStream(receiverService)
	if stream == nil {
		return errors.New("stream is null")
	}

	return stream.Send(&Req{
		ServiceIdentify: senderServiceIdentify,
		UserSessionId:   userSessionId,
		Data:            data,
	})
}
