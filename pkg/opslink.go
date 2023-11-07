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

	doneChan   chan struct{}
	closedChan chan struct{}
}

func NewOpsLinkServer(config *config.OpsLinkConfig, casbin *casbin.Casbin, store opsstore.StoreService, kcs *kubeclient.KubernetesClient) (os *OpsLinkServer, err error) {
	os = &OpsLinkServer{
		Config:        config,
		Casbin:        casbin,
		StoreService:  store,
		KubeClientSet: kcs,
		closedChan:    make(chan struct{}),
	}
	return
}

func (server *OpsLinkServer) Start() error {
	server.doneChan = make(chan struct{})

	<-server.doneChan

	close(server.closedChan)
	return nil
}

func (server *OpsLinkServer) Stop(force bool) {
	//todo:如果使用本地内存启动，关闭之前等待数据同步
	//Before closing, check if there is any unsynchronized data.
	server.StoreService.Stop()
	close(server.doneChan)

	// wait for fully closed
	<-server.closedChan
}
