package grpc

import (
	. "core/libs"
	"core/libs/consul"
	"core/libs/hash"
	"core/libs/timer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"reflect"
	"sort"
	"sync"
	"time"
)

type GrpcClient struct {
	consulClient *consul.ConsulClient
	serviceName  string

	services      []string
	servicesMutex sync.Mutex

	links     map[string]*grpc.ClientConn
	linkMutex sync.Mutex

	newPbClientFunc func(*grpc.ClientConn) interface{}
}

func InitClient(consulClient *consul.ConsulClient, serviceName string, newPbClientFunc func(*grpc.ClientConn) interface{}) *GrpcClient {
	client := &GrpcClient{
		consulClient:    consulClient,
		serviceName:     serviceName,
		links:           make(map[string]*grpc.ClientConn),
		newPbClientFunc: newPbClientFunc,
	}
	client.initServices()
	client.loop()
	return client
}

func (this *GrpcClient) loop() {
	timer.DoTimer(10*1000, this.initServices)
}

func (this *GrpcClient) initServices() {
	this.servicesMutex.Lock()
	this.services = this.consulClient.GetServices(this.serviceName)
	sort.Strings(this.services)
	this.servicesMutex.Unlock()

	this.traceServices()
}

func (this *GrpcClient) traceServices() {
	this.servicesMutex.Lock()
	for _, value := range this.services {
		DEBUG(this.serviceName, "Service", value)
	}
	DEBUG("--------------------------------------------")
	this.servicesMutex.Unlock()
}

func (this *GrpcClient) removeService(service string) {
	this.servicesMutex.Lock()
	for index, value := range this.services {
		if value == service {
			this.services = append(this.services[:index], this.services[index+1:]...)
		}
	}
	this.servicesMutex.Unlock()

	this.traceServices()
}

func (this *GrpcClient) GetServiceByFlag(flag string) string {
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

func (this *GrpcClient) GetServiceByRandom() string {
	return this.GetServiceByFlag(NumToString(time.Now().Unix()))
}

func (this *GrpcClient) getLink(service string) *grpc.ClientConn {
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
		ERR("grpcServer connect fail", service)
		return nil
	} else {
		INFO("grpcServer connect success", service)
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

func (this *GrpcClient) removeLink(service string) {
	this.linkMutex.Lock()
	if link, ok := this.links[service]; ok {
		link.Close()
		delete(this.links, service)
	}
	this.linkMutex.Unlock()

	ERR("grpcServer disconnected", service)
}

func (this *GrpcClient) Call(service string, serviceMethod string, arg interface{}) interface{} {
	////根据Flag分配Service
	//if flag == "" {
	//	flag = NumToString(time.Now().Unix())
	//}
	//service := this.getServiceByFlag(flag)
	//if service == "" {
	//	ERR("grpcServer no exists")
	//	return nil, ""
	//}

	//根据Service获取链接
	link := this.getLink(service)
	if link == nil {
		this.removeService(service)
		return this.Call(service, serviceMethod, arg)
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
		return this.Call(service, serviceMethod, arg)
	}

	return reply
}
