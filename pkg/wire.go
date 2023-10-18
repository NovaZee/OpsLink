//go:build wireinject
// +build wireinject

package pkg

import (
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/casbin"
	etcdv3 "github.com/denovo/permission/pkg/etcdv3"
	"github.com/google/wire"
)

func InitializeServer(cfg *config.OpsLinkConfig) (*OpsLinkServer, error) {
	wire.Build(
		initCasbin,
		initEtcd,
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
