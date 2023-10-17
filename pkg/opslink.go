package pkg

import (
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/casbin"
	etcdv3 "github.com/denovo/permission/pkg/clientv3"
)

type OpsLinkServer struct {
	config    *config.OpsLinkConfig
	Casbin    *casbin.Casbin
	Interface *etcdv3.Interface
}

func NewOpsLinkServer(config *config.OpsLinkConfig, casbin *casbin.Casbin) (os *OpsLinkServer, err error) {
	os = &OpsLinkServer{
		config: config,
		Casbin: casbin,
	}
	return
}
