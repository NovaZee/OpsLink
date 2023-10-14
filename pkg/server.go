package pkg

import (
	"github.com/coreos/etcd/clientv3"
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/casbin"
)

type OpsLinkServer struct {
	config *config.Config
	Etcd   *clientv3.Client
	Casbin *casbin.Casbin
}

func NewOpsLinkServer(config *config.Config, casbin *casbin.Casbin) (os *OpsLinkServer, err error) {
	os = &OpsLinkServer{
		config: config,
		Casbin: casbin,
	}
	return
}
