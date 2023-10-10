//go:build wireinject
// +build wireinject

package pkg

import (
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/casbin"
	"github.com/denovo/permission/pkg/componment"
	"github.com/denovo/permission/pkg/router"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializeServer(cfg *config.Config) (*OpsLinkServer, error) {
	wire.Build(
		getDbEngine,
		getCasbinAdapter,
		initRouter,
		NewOpsLinkServer,
	)
	return &OpsLinkServer{}, nil
}

func getDbEngine(conf *config.Config) (*gorm.DB, error) {
	return componment.InitDBConnection(conf)
}

func getCasbinAdapter(engine *gorm.DB, conf *config.Config) *casbin.CasbinAdapter {
	return casbin.NewCasbinAdapter(engine, conf)
}

func initRouter(ca *casbin.CasbinAdapter) (*router.Router, error) {
	return router.InitRouter(ca)
}
