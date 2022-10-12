package grpc_resolver_nacos

import "github.com/wubin1989/nacos-sdk-go/v2/vo"

func GetService(serviceName string, clusters []string, groupName string, nacosClient NacosClient) (NacosService, error) {
	service, err := nacosClient.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
		Clusters:    clusters,
		GroupName:   groupName, // 默认值DEFAULT_GROUP
	})
	return service, err
}

func GetOneHealthyInstance(serviceName string, clusters []string, groupName string, nacosClient NacosClient) (*NacosInstance, error) {
	instance, err := nacosClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
		Clusters:    clusters,
		GroupName:   groupName, // 默认值DEFAULT_GROUP
	})
	return instance, err
}

func GetHealthyInstances(serviceName string, clusters []string, groupName string, nacosClient NacosClient) ([]NacosInstance, error) {
	instances, err := nacosClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		Clusters:    clusters,
		GroupName:   groupName, // 默认值DEFAULT_GROUP
		HealthyOnly: true,
	})
	return instances, err
}

func RegisterInstance(ip string, port uint64, serviceName string, weight float64, enable bool, healthy bool,
	metadata map[string]string, clusterName string, groupName string, ephemeral bool, nacosClient NacosClient) (bool, error) {
	param := vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: serviceName,
		Weight:      weight,
		Enable:      enable,
		Healthy:     healthy,
		Metadata:    metadata,
		ClusterName: clusterName,
		GroupName:   groupName,
		Ephemeral:   ephemeral,
	}
	success, err := nacosClient.RegisterInstance(param)
	return success, err
}
