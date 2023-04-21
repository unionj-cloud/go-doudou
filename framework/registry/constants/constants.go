package constants

const (
	SD_NACOS      = "nacos"
	SD_ETCD       = "etcd"
	SD_MEMBERLIST = "memberlist"
	SD_ZK         = "zk"
)

type ServiceType string

const (
	GRPC_TYPE ServiceType = "grpc"
	REST_TYPE ServiceType = "rest"
)
