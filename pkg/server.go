package pkg

import (
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/casbin"
	"github.com/denovo/permission/pkg/router"
	"gorm.io/gorm"
)

type OpsLinkServer struct {
	config        *config.Config
	DBEngine      *gorm.DB
	CasbinAdapter *casbin.CasbinAdapter
	router        *router.Router
}

func NewOpsLinkServer(config *config.Config, DBEngine *gorm.DB, CasbinAdapter *casbin.CasbinAdapter, router *router.Router) (os *OpsLinkServer, err error) {
	os = &OpsLinkServer{
		config:        config,
		DBEngine:      DBEngine,
		CasbinAdapter: CasbinAdapter,
		router:        router,
	}
	return
}
