package rpc

import (
	"GoGameServer/core/libs/common"
	"GoGameServer/core/libs/consul"
	"GoGameServer/core/libs/hash"
	"GoGameServer/core/libs/logger"
	"GoGameServer/core/libs/timer"
	"errors"
	"io"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"
)

type Client struct {
	consulClient *consul.Client
	serviceName  string

	services      []string
	servicesMutex sync.Mutex

	links     map[string]*rpc.Client
	linkMutex sync.Mutex
}

func NewClient(consulClient *consul.Client, serviceName string) *Client {
	client := &Client{
		consulClient: consulClient,
		serviceName:  serviceName,
		links:        make(map[string]*rpc.Client),
	}
	client.initServices()
	client.loop()
	return client
}

func (this *Client) loop() {
	timer.SetTimeOut(5*1000, this.initServices)
}

func (this *Client) initServices() {
	this.servicesMutex.Lock()
	this.services = this.consulClient.GetServices(this.serviceName)
	this.servicesMutex.Unlock()

	this.initLinks()
	this.traceServices()
}

func (this *Client) initLinks() {
	if len(this.services) == 0 {
		timer.SetTimeOut(300, this.initServices)
		return
	}

	if len(this.links) == 0 {
		for _, service := range this.services {
			this.getLink(service)
		}
	}
}

func (this *Client) traceServices() {
	return
	this.servicesMutex.Lock()
	logger.Debug("----------rpc start " + this.serviceName + "----------")
	for _, value := range this.services {
		logger.Debug(this.serviceName, "Service", value)
	}
	logger.Debug("-----------rpc end " + this.serviceName + "-----------")
	this.servicesMutex.Unlock()
}

func (this *Client) removeService(service string) {
	this.servicesMutex.Lock()
	for index, value := range this.services {
		if value == service {
			this.services = append(this.services[:index], this.services[index+1:]...)
		}
	}
	this.servicesMutex.Unlock()

	this.traceServices()
}

func (this *Client) getServiceByFlag(flag string) string {
	this.servicesMutex.Lock()
	service := ""
	servicesLen := len(this.services)
	if servicesLen > 0 {
		num := hash.GetHash([]byte(flag))
		index := int(num % uint32(servicesLen))
		service = this.services[index]
	}
	this.servicesMutex.Unlock()

	return service
}

func (this *Client) getLink(service string) *rpc.Client {
	//监测是否已经存在
	this.linkMutex.Lock()
	link, ok := this.links[service]
	this.linkMutex.Unlock()

	if ok {
		return link
	}

	//连接Rpc服务器
	conn, err := net.DialTimeout("tcp", service, time.Second*3)
	if err != nil {
		logger.Error("rpcServer connect fail", service)
		return nil
	} else {
		logger.Info("rpcServer connect success", service)
	}
	link = jsonrpc.NewClient(conn)

	//防止重复链接
	this.linkMutex.Lock()
	if link2, ok := this.links[service]; ok {
		link.Close()
		link = link2
	} else {
		this.links[service] = link
	}
	this.linkMutex.Unlock()

	return link
}

func (this *Client) removeLink(service string) {
	this.linkMutex.Lock()
	if link, ok := this.links[service]; ok {
		link.Close()
		delete(this.links, service)
	}
	this.linkMutex.Unlock()

	logger.Error("rpcServer disconnected", service)
}

func (this *Client) Call(serviceMethod string, args interface{}, reply interface{}, flag string) error {
	if flag == "" {
		flag = common.NumToString(time.Now().Unix())
	}
	service := this.getServiceByFlag(flag)
	if service == "" {
		return errors.New("rpcServer no exists")
	}

	link := this.getLink(service)
	if link == nil {
		this.removeService(service)
		return this.Call(serviceMethod, args, reply, flag)
	}

	err := link.Call(serviceMethod, args, reply)
	if err == io.ErrUnexpectedEOF || err == rpc.ErrShutdown {
		this.removeLink(service)
		return this.Call(serviceMethod, args, reply, flag)
	}
	return err
}

func (this *Client) CallAll(serviceMethod string, args interface{}, reply interface{}) {
	for _, value := range this.services {
		link := this.getLink(value)
		if link == nil {
			continue
		}

		err := link.Call(serviceMethod, args, reply)
		if err == io.ErrUnexpectedEOF || err == rpc.ErrShutdown {
			this.removeLink(value)
		}
	}
}
