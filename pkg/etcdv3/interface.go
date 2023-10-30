package etcdv3

type Interface interface {
	RolesClient
}

type RolesClient interface {
	// RolesCfg RolesCfg() RoleClientInterface
	RolesCfg() RoleClientInterface
}
