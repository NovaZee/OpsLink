//go:build wireinject
// +build wireinject

package pkg

import (
	clientv3 "github.com/coreos/etcd/clientv3"
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/casbin"
	"github.com/denovo/permission/pkg/etcd"
	"github.com/google/wire"
)

func InitializeServer(cfg *config.Config) (*OpsLinkServer, error) {
	wire.Build(
		initEtcd,
		initCasbin,
		NewOpsLinkServer,
	)
	return &OpsLinkServer{}, nil
}

func initCasbin(conf *config.Config) (*casbin.Casbin, error) {
	return casbin.InitCasbin(conf)
}
func initEtcd(conf *config.Config) (*clientv3.Client, error) {
	return etcd.InitEtcd(conf)
}
