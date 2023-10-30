package pkg

import (
	config "github.com/denovo/permission/config"
	"github.com/denovo/permission/pkg/casbin"
	"github.com/denovo/permission/pkg/kubeclient"
	opsstore "github.com/denovo/permission/pkg/store"
)

type OpsLinkServer struct {
	Config       *config.OpsLinkConfig
	Casbin       *casbin.Casbin
	StoreService opsstore.StoreService

	KubeClientSet *kubeclient.KubernetesClient
}

func NewOpsLinkServer(config *config.OpsLinkConfig, casbin *casbin.Casbin, store opsstore.StoreService, kcs *kubeclient.KubernetesClient) (os *OpsLinkServer, err error) {
	os = &OpsLinkServer{
		Config:        config,
		Casbin:        casbin,
		StoreService:  store,
		KubeClientSet: kcs,
	}
	return
}
