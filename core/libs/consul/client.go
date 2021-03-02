package consul

import (
	"GoGameServer/core/libs/array"
	"GoGameServer/core/libs/common"
	"GoGameServer/core/libs/stack"
	"sort"
	"strings"

	"github.com/hashicorp/consul/api"
)

type Client struct {
	consulClient *api.Client
}

type ServiceInfo struct {
	ID      string
	Name    string
	Address string
	Port    string
	SortKey string
}

func NewClient() (*Client, error) {
	//开启consulKV
	err := InitKV(true)
	stack.CheckError(err)

	//开启consul客户端
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	consulClient := &Client{
		consulClient: client,
	}
	return consulClient, nil
}

//func (this *Client) GetServices(service string) []string {
//	services, _ := this.consulClient.Agent().Services()
//	results := []string{}
//	if services != nil {
//		for _, value := range services {
//			if value.Service != service {
//				continue
//			}
//			addr := value.Address + ":" + NumToString(value.Port)
//			results = append(results, addr)
//		}
//	}
//	return results
//}

func getFilterServices() []string {
	filterServicesStr := KV_Get("FilterServices")
	arr := strings.Split(filterServicesStr, ";")
	var result = []string{}
	for _, service := range arr {
		if len(service) == 0 {
			continue
		}
		result = append(result, strings.TrimSpace(service))
	}
	return result
}

func (this *Client) GetServices(service string) []string {
	services, _, _ := this.consulClient.Health().Service(service, "", true, &api.QueryOptions{})
	filterServices := getFilterServices()
	serviceDatas := []ServiceInfo{}
	if services != nil {
		for _, entry := range services {
			if array.InArray(filterServices, entry.Service.Address) {
				continue
			}

			arr := strings.Split(entry.Service.ID, "-")
			serviveId := arr[2]
			data := ServiceInfo{
				ID:      entry.Service.ID,
				Name:    entry.Service.Service,
				Address: entry.Service.Address,
				Port:    common.NumToString(entry.Service.Port),
				SortKey: entry.Service.Address + "-" + serviveId,
			}
			serviceDatas = append(serviceDatas, data)
		}
	}

	//排序(从小到大)
	sort.Slice(serviceDatas, func(i, j int) bool {
		return serviceDatas[i].SortKey < serviceDatas[j].SortKey
	})

	//组装返回数据
	results := []string{}
	for i := 0; i < len(serviceDatas); i++ {
		data := serviceDatas[i]
		addr := data.Address + ":" + data.Port
		results = append(results, addr)
	}
	return results
}

func (this *Client) DeregisterService(serviceID string) {
	this.consulClient.Agent().ServiceDeregister(serviceID)
}
