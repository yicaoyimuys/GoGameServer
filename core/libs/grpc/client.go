package grpc

import (
	"GoGameServer/core/libs/common"
	"GoGameServer/core/libs/consul"
	"GoGameServer/core/libs/hash"
	"GoGameServer/core/libs/logger"
	"GoGameServer/core/libs/timer"
	"io"
	"reflect"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Client struct {
	consulClient *consul.Client
	serviceName  string

	services      []string
	servicesMutex sync.Mutex

	links     map[string]*grpc.ClientConn
	linkMutex sync.Mutex

	newPbClientFunc func(*grpc.ClientConn) interface{}
}

func NewClient(consulClient *consul.Client, serviceName string, newPbClientFunc func(*grpc.ClientConn) interface{}) *Client {
	client := &Client{
		consulClient:    consulClient,
		serviceName:     serviceName,
		links:           make(map[string]*grpc.ClientConn),
		newPbClientFunc: newPbClientFunc,
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
	logger.Debug("----------grpc start " + this.serviceName + "----------")
	for _, value := range this.services {
		logger.Debug(this.serviceName, "Service", value)
	}
	logger.Debug("-----------grpc end " + this.serviceName + "-----------")
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

func (this *Client) GetServiceByFlag(flag string) string {
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

func (this *Client) GetServiceByRandom() string {
	return this.GetServiceByFlag(common.NumToString(time.Now().Unix()))
}

func (this *Client) getLink(service string) *grpc.ClientConn {
	//监测是否已经存在
	this.linkMutex.Lock()
	link, ok := this.links[service]
	this.linkMutex.Unlock()

	if ok {
		return link
	}

	//连接Rpc服务器
	link, err := grpc.Dial(service, grpc.WithInsecure())
	if err != nil {
		logger.Error("grpcServer connect fail", service)
		return nil
	} else {
		logger.Info("grpcServer connect success", service)
	}

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

	logger.Error("grpcServer disconnected", service)
}

func (this *Client) Call(service string, serviceMethod string, arg interface{}) interface{} {
	//根据Service获取链接
	link := this.getLink(service)
	if link == nil {
		this.removeService(service)
		logger.Error("service not exists", service)
		return nil
	}

	//创建PbClient
	client := this.newPbClientFunc(link)

	//创建调用serviceMethod所需要的参数
	mArgs := []reflect.Value{
		reflect.ValueOf(context.Background()),
	}
	if arg != nil {
		mArgs = append(mArgs, reflect.ValueOf(arg))
	}

	//调用serviceMethod
	clientReflect := reflect.ValueOf(client)
	serviceResult := clientReflect.MethodByName(serviceMethod).Call(mArgs)

	//结果
	reply := serviceResult[0].Interface()
	err := serviceResult[1].Interface()
	if err != nil && err.(error) == io.ErrUnexpectedEOF {
		this.removeLink(service)
		logger.Error("service is error", service)
		return nil
	}

	return reply
}
