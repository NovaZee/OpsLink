package pkg

import (
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/casbin"
	etcdv3 "github.com/denovo/permission/pkg/etcdv3"
)

type OpsLinkServer struct {
	Config    *config.OpsLinkConfig
	Casbin    *casbin.Casbin
	Interface etcdv3.Interface
}

func NewOpsLinkServer(config *config.OpsLinkConfig, casbin *casbin.Casbin, interface1 etcdv3.Interface) (os *OpsLinkServer, err error) {
	os = &OpsLinkServer{
		Config:    config,
		Casbin:    casbin,
		Interface: interface1,
	}
	return
}
