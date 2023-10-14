package etcd

import (
	"github.com/denovo/permission/pkg/service/role"
)

const CasbinRuleKey = "/casbin_policy/"
const RoleKey = "/role_key/"

type EtcdClient interface {
	PermissionClient
	DiscoveryClient
}

type PermissionClient interface {
	SetPermissionPolicy(v any) error
	GetPermissionPolicy(v any) ([]role.Role, error)
	DeletePermissionPolicy(v string) (int64, error)
	UpdatePermissionPolicy(o any, n any) error
}

type DiscoveryClient interface {
	RegisterService(key, value string) error
	UnregisterService(key string) error
	DiscoverService(key string) (string, error)
}
