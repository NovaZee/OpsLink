package router

import (
	"context"
	"github.com/denovo/permission/pkg/casbin"
	"github.com/denovo/permission/pkg/service"
	store "github.com/denovo/permission/pkg/store"
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

	storeService store.StoreService
}

func InitRouter(opslinkServer *service.OpsLinkServer) (*Router, error) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	engine := gin.Default()
	engine.Use(Logger())
	engine.GET("/validate")

	router, err := NewRouter(engine, opslinkServer.Casbin, opslinkServer.StoreService)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	router.InitPolicyRouting()
	router.InitUserRouting(ctx)
	engine.Run(":" + opslinkServer.Config.Server.HttpPort).Error()
	return router, nil
}

func NewRouter(g *gin.Engine, ca *casbin.Casbin, ss store.StoreService) (*Router, error) {
	return &Router{
		router:       g,
		cb:           ca,
		storeService: ss,
	}, nil
}

// InitPolicyRouting InitAdminRouting 管理员路由
func (r *Router) InitPolicyRouting() {
	admin := r.router.Group("/manager")
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
	admin.Use(ManagerMiddleware())
}

// InitUserRouting 用户注册路由
func (r *Router) InitUserRouting(ctxEtcd context.Context) {
	admin := r.router.Group("/")
	{
		admin.POST("logIn", func(ctx *gin.Context) {
			LogIn(ctx, r, ctxEtcd)
		})
		admin.POST("signIn", func(ctx *gin.Context) {
			SignIn(ctx, r, ctxEtcd)
		})
	}
}

// InitAccessingRouting 用户访问路由
func (r *Router) InitAccessingRouting() {
	// 访问请求通过jwt校验->casbin校验
	admin := r.router.Group("/v1")
	{
		admin.POST("index", func(ctx *gin.Context) {
		})

	}
	admin.Use(JWT(r))
}
