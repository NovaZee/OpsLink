package pkg

import (
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/casbin"
	etcdv3 "github.com/denovo/permission/pkg/clientv3"
	"github.com/denovo/permission/pkg/router"
)

type OpsLinkServer struct {
	config    *config.OpsLinkConfig
	Casbin    *casbin.Casbin
	Interface *etcdv3.Interface
	Router    *router.Router
}

func NewOpsLinkServer(config *config.OpsLinkConfig, casbin *casbin.Casbin, interface1 *etcdv3.Interface, router *router.Router) (os *OpsLinkServer, err error) {
	os = &OpsLinkServer{
		config:    config,
		Casbin:    casbin,
		Interface: interface1,
		Router:    router,
	}
	return
}
