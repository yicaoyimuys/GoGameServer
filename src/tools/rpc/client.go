package rpc

import (
	"errors"
	"io"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sort"
	"sync"
	"time"
	. "tools"
	"tools/consul"
	"tools/hash"
	"tools/timer"
)

type RpcClient struct {
	consulClient *consul.ConsulClient
	serviceName  string

	services      []string
	servicesMutex sync.Mutex

	links     map[string]*rpc.Client
	linkMutex sync.Mutex
}

func InitClient(consulClient *consul.ConsulClient, serviceName string) *RpcClient {
	client := &RpcClient{
		consulClient: consulClient,
		serviceName:  serviceName,
		links:        make(map[string]*rpc.Client),
	}
	client.initServices()
	client.loop()
	return client
}

func (this *RpcClient) loop() {
	timer.DoTimer(10*1000, this.initServices)
}

func (this *RpcClient) initServices() {
	this.servicesMutex.Lock()
	this.services = this.consulClient.GetServices(this.serviceName)
	sort.Strings(this.services)
	this.servicesMutex.Unlock()

	this.traceServices()
}

func (this *RpcClient) traceServices() {
	this.servicesMutex.Lock()
	for _, value := range this.services {
		DEBUG(this.serviceName, "Service", value)
	}
	DEBUG("--------------------------------------------")
	this.servicesMutex.Unlock()
}

func (this *RpcClient) removeService(service string) {
	this.servicesMutex.Lock()
	for index, value := range this.services {
		if value == service {
			this.services = append(this.services[:index], this.services[index+1:]...)
		}
	}
	this.servicesMutex.Unlock()

	this.traceServices()
}

func (this *RpcClient) getServiceByFlag(flag string) string {
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

func (this *RpcClient) getLink(service string) *rpc.Client {
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
		ERR("rpcServer connect fail", service)
		return nil
	} else {
		INFO("rpcServer connect success", service)
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

func (this *RpcClient) removeLink(service string) {
	this.linkMutex.Lock()
	if link, ok := this.links[service]; ok {
		link.Close()
		delete(this.links, service)
	}
	this.linkMutex.Unlock()

	ERR("rpcServer disconnected", service)
}

func (this *RpcClient) Call(serviceMethod string, args interface{}, reply interface{}, flag string) error {
	if flag == "" {
		flag = NumToString(time.Now().Unix())
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

	err := link.Call(this.serviceName+"."+serviceMethod, args, reply)
	if err == io.ErrUnexpectedEOF || err == rpc.ErrShutdown {
		this.removeLink(service)
		return this.Call(serviceMethod, args, reply, flag)
	}
	return err
}
