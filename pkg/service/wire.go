//go:build wireinject
// +build wireinject

package service

import (
	config "github.com/denovo/permission/config"
	"github.com/denovo/permission/pkg/casbin"
	kubeclient "github.com/denovo/permission/pkg/kubeclient"
	"github.com/denovo/permission/pkg/store"
	"github.com/google/wire"
)

func InitializeServer(cfg *config.OpsLinkConfig) (*OpsLinkServer, error) {
	wire.Build(
		initCasbin,
		initStore,
		initClientSet,
		initSignalService,
		NewOpsLinkServer,
	)
	return &OpsLinkServer{}, nil
}

func initCasbin(conf *config.OpsLinkConfig) (*casbin.Casbin, error) {
	return casbin.InitCasbin(conf)
}

func initStore(conf *config.OpsLinkConfig) (store.StoreService, error) {
	return store.NewStoreService(conf)
}

func initSignalService(conf *config.OpsLinkConfig) (*SignalService, error) {
	return NewSignalService(), nil
}

func initClientSet(conf *config.OpsLinkConfig) (*kubeclient.KubernetesClient, error) {
	clinetInterface, err := kubeclient.NewClientInterface(conf, kubeclient.K8sClientTypeKubernetes)
	if err != nil {
		return nil, err
	}
	return kubeclient.GetClientSet(clinetInterface), nil
}
