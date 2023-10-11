// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package pkg

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/casbin"
	"github.com/denovo/permission/pkg/componment"
)

// Injectors from wire.go:

func InitializeServer(cfg *config.Config) (*OpsLinkServer, error) {
	client, err := initCasbin(cfg)
	if err != nil {
		return nil, err
	}
	casbin, err := initEtcd(cfg)
	if err != nil {
		return nil, err
	}
	opsLinkServer, err := NewOpsLinkServer(cfg, client, casbin)
	if err != nil {
		return nil, err
	}
	return opsLinkServer, nil
}

// wire.go:

func initEtcd(conf *config.Config) (*casbin.Casbin, error) {
	return casbin.InitCasbin(conf)
}

func initCasbin(conf *config.Config) (*clientv3.Client, error) {
	return componment.InitEtcd(conf)
}
