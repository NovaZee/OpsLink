package router

import (
	"errors"
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/casbin"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

const (
	ErrorAuthCheckTokenFail    = " check fail "
	ErrorAuthCheckTokenExpired = " token expired "
	ErrorParamsError           = " bind params error  "
)

type Router struct {
	router *gin.Engine
	cb     *casbin.Casbin
}

func InitRouter(ca *casbin.Casbin, conf *config.Config) error {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	engine := gin.Default()
	engine.Use(Logger(), JWT())
	engine.Run(conf.Server.HttpPort)
	router, err := NewRouter(engine, ca)
	if err != nil {
		return err
	}
	router.InitAdminRouting()
	s := engine.Run(":" + conf.Server.HttpPort).Error()
	if len(s) != 0 {
		return errors.New(s)
	}
	return nil
}

func NewRouter(g *gin.Engine, ca *casbin.Casbin) (*Router, error) {
	return &Router{
		router: g,
		cb:     ca,
	}, nil
}

func (r *Router) InitAdminRouting() {
	admin := r.router.Group("/manager", ManagerMiddleware())
	{
		admin.POST("addPolicy", func(ctx *gin.Context) {
			AddPolicy(ctx, r.cb)
		})
		admin.GET("deletePolicy", func(ctx *gin.Context) {
			DeletePolicy(ctx, r.cb)
		})
		admin.POST("update", func(ctx *gin.Context) {
		})
	}
}
func (r *Router) InitLoginRouting() {
	admin := r.router.Group("/v1", ManagerMiddleware())
	{
		admin.POST("login", func(ctx *gin.Context) {
			LogIn(ctx)
		})
		admin.GET("signIn", func(ctx *gin.Context) {
			SignIn(ctx)
		})
	}
}
