//go:build wireinject
// +build wireinject

package pkg

import (
	clientv3 "github.com/coreos/etcd/clientv3"
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/casbin"
	etcdv3 "github.com/denovo/permission/pkg/clientv3"
	"github.com/denovo/permission/pkg/etcd"
	"github.com/google/wire"
)

func InitializeServer(cfg *config.Config) (*OpsLinkServer, error) {
	wire.Build(
		initCasbin,
		initEtcd,
		initRouter,
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

func initEtcd(conf *config.Config) (*etcdv3.Interface, error) {
	return etcdv3.New(cfg)
}
func initRouter(casbin *casbin.Casbin, conf *config.Config, etcdV3 *etcdv3.Interface) (*clientv3.Client, error) {
	return router.InitRouter(casbin, conf, etcdV3)
}

//back, err := etcdv3.New(cfg)
//if err != nil {
//return err
//}
//
////rolesCfg := back.RolesCfg()
//logger.Infow("start server ", "port", cfg.Server.HttpPort)
//error = router.InitRouter(server.Casbin, cfg, back)
//if error != nil {
//return error
//}
