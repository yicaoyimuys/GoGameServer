package consul

import (
	. "core/libs"
	"core/libs/array"
	"github.com/hashicorp/consul/api"
	"strings"
)

type ConsulClient struct {
	consulClient *api.Client
}

func InitClient() (*ConsulClient, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	consulClient := &ConsulClient{
		consulClient: client,
	}
	return consulClient, nil
}

//func (this *ConsulClient) GetServices(service string) []string {
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

func (this *ConsulClient) GetServices(service string) []string {
	services, _, _ := this.consulClient.Health().Service(service, "", true, &api.QueryOptions{})
	filterServices := getFilterServices()
	results := []string{}
	if services != nil {
		for _, entry := range services {
			if array.InArray(filterServices, entry.Service.Address) {
				continue
			}
			addr := entry.Service.Address + ":" + NumToString(entry.Service.Port)
			results = append(results, addr)
		}
	}
	return results
}

func (this *ConsulClient) DeregisterService(serviceID string) {
	this.consulClient.Agent().ServiceDeregister(serviceID)
}
