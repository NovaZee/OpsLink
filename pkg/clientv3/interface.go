package etcdv3

type Interface interface {
	RolesClient
}

type RolesClient interface {
	//RolesCfg() RoleClientInterface
	RolesCfg() RoleClientInterface
}
