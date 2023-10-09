package grpc

import (
	"io"
	"reflect"
	"sync"
	"time"

	"github.com/yicaoyimuys/GoGameServer/core/libs/consul"
	"github.com/yicaoyimuys/GoGameServer/core/libs/hash"
	"github.com/yicaoyimuys/GoGameServer/core/libs/logger"
	"github.com/yicaoyimuys/GoGameServer/core/libs/timer"
	"go.uber.org/zap"

	"github.com/spf13/cast"
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
		logger.Debug("Service", zap.String("ServiceName", this.serviceName), zap.String("ServiceAddress", value))
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
	return this.GetServiceByFlag(cast.ToString(time.Now().Unix()))
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
		logger.Error("GrpcServer Connect Fail", zap.String("Service", service))
		return nil
	} else {
		logger.Info("GrpcServer Connect Success", zap.String("Service", service))
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

	logger.Error("grpcServer disconnected", zap.String("Service", service))
}

func (this *Client) Call(service string, serviceMethod string, arg interface{}) interface{} {
	//根据Service获取链接
	link := this.getLink(service)
	if link == nil {
		this.removeService(service)
		logger.Error("service not exists", zap.String("Service", service))
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
		logger.Error("service is error", zap.String("Service", service))
		return nil
	}

	return reply
}
