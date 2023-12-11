//go:build wireinject
// +build wireinject

package service

import (
	config "github.com/denovo/permission/config"
	"github.com/denovo/permission/pkg/casbin"
	"github.com/denovo/permission/pkg/kubenates"
	"github.com/denovo/permission/pkg/signal"
	"github.com/denovo/permission/pkg/store"
	"github.com/google/wire"
)

func InitializeServer(cfg *config.OpsLinkConfig) (*OpsLinkServer, error) {
	wire.Build(
		initCasbin,
		initStore,
		initOpsKube,
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

func initSignalService(conf *config.OpsLinkConfig) (*signal.SignalService, error) {
	return signal.NewSignalService(), nil
}
func initOpsKube(conf *config.OpsLinkConfig) (*kubenates.K8sClient, error) {
	return kubenates.NewK8sConfig(conf)
}

//func initClientSet(conf *config.OpsLinkConfig) (*kubenates.KubernetesClient, error) {
//	clinetInterface, err := kubenates.NewClientInterface(conf, kubenates.K8sClientTypeKubernetes)
//	if err != nil {
//		return nil, err
//	}
//	return kubenates.GetClientSet(clinetInterface), nil
//}
