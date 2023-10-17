package router

import (
	"context"
	"errors"
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/casbin"
	etcdv3 "github.com/denovo/permission/pkg/clientv3"
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

	roleClientv3 etcdv3.RoleClientInterface
}

func InitRouter(ca *casbin.Casbin, conf *config.OpsLinkConfig, back etcdv3.Interface) error {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	engine := gin.Default()
	engine.Use(Logger())
	engine.Run(conf.Server.HttpPort)
	router, err := NewRouter(engine, ca, back)
	if err != nil {
		return err
	}
	ctx := context.Background()
	router.InitAdminRouting()
	router.InitUserRouting(ctx)
	s := engine.Run(":" + conf.Server.HttpPort).Error()
	if len(s) != 0 {
		return errors.New(s)
	}
	return nil
}

func NewRouter(g *gin.Engine, ca *casbin.Casbin, back etcdv3.Interface) (*Router, error) {
	return &Router{
		router:       g,
		cb:           ca,
		roleClientv3: back.RolesCfg(),
	}, nil
}

// InitAdminRouting 管理员路由
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

// InitUserRouting 用户注册路由
func (r *Router) InitUserRouting(ctxEtcd context.Context) {
	admin := r.router.Group("/")
	{
		admin.POST("logIn", func(ctx *gin.Context) {
			LogIn(ctx)
		})
		admin.POST("signIn", func(ctx *gin.Context) {
			SignIn(ctx, r, ctxEtcd)
		})
	}
}

// InitAccessingRouting 用户访问路由
func (r *Router) InitAccessingRouting() {
	// 访问请求通过jwt校验->casbin校验
	admin := r.router.Group("/v1", JWT(r))
	{
		admin.POST("index", func(ctx *gin.Context) {
		})

	}
}
