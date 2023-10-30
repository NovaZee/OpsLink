//go:build wireinject
// +build wireinject

package pkg

import (
	config "github.com/denovo/permission/config"
	"github.com/denovo/permission/pkg/casbin"
	kubeclient "github.com/denovo/permission/pkg/kubeclient"
	etcdv3 "github.com/denovo/permission/pkg/store"
	"github.com/google/wire"
)

func InitializeServer(cfg *config.OpsLinkConfig) (*OpsLinkServer, error) {
	wire.Build(
		initCasbin,
		initEtcd,
		initClientSet,
		NewOpsLinkServer,
	)
	return &OpsLinkServer{}, nil
}

func initCasbin(conf *config.OpsLinkConfig) (*casbin.Casbin, error) {
	return casbin.InitCasbin(conf)
}

func initEtcd(conf *config.OpsLinkConfig) (etcdv3.Interface, error) {
	return etcdv3.New(conf)
}

func initClientSet(conf *config.OpsLinkConfig) (*kubeclient.KubernetesClient, error) {
	clinetInterface, err := kubeclient.NewClientInterface(conf, kubeclient.K8sClientTypeKubernetes)
	if err != nil {
		return nil, err
	}
	return kubeclient.GetClientSet(clinetInterface), nil
}
