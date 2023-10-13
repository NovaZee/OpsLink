package etcd

type EtcdClient interface {
	PermissionClient
	DiscoveryClient
}

type PermissionClient interface {
	SetPermissionPolicy(v any) error
	GetPermissionPolicy(v any) (string, error)
	DeletePermissionPolicy(key string) error
}

type DiscoveryClient interface {
	RegisterService(key, value string) error
	UnregisterService(key string) error
	DiscoverService(key string) (string, error)
}
