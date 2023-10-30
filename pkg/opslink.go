package pkg

import (
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/casbin"
	etcdv3 "github.com/denovo/permission/pkg/etcdv3"
	"github.com/denovo/permission/pkg/kubeclient"
)

type OpsLinkServer struct {
	Config    *config.OpsLinkConfig
	Casbin    *casbin.Casbin
	Interface etcdv3.Interface

	kubeClientSet *kubeclient.KubernetesClient
}

func NewOpsLinkServer(config *config.OpsLinkConfig, casbin *casbin.Casbin, interface1 etcdv3.Interface, kcs *kubeclient.KubernetesClient) (os *OpsLinkServer, err error) {
	os = &OpsLinkServer{
		Config:        config,
		Casbin:        casbin,
		Interface:     interface1,
		kubeClientSet: kcs,
	}
	return
}
