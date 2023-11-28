package router

import (
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
	Router *gin.Engine
	cb     *casbin.Casbin

	storeService store.StoreService
	handler      []Handler
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
	router.InitHandler(opslinkServer)
	registerHandlers(router, router.handler...)
	engine.Run(":" + opslinkServer.Config.Server.HttpPort).Error()
	return router, nil
}

func (r *Router) InitHandler(opslinkServer *service.OpsLinkServer) {
	clientSet := opslinkServer.K8sClient.Clientset
	handlers := []Handler{
		BuildDeployments(clientSet, opslinkServer.K8sClient.DepHandler),
	}
	r.handler = handlers
}

// registerHandlers 将多个处理程序注册到 Gin 路由器上
func registerHandlers(router *Router, handlers ...Handler) {
	for _, h := range handlers {
		h.Register(router)
	}
}

func NewRouter(g *gin.Engine, ca *casbin.Casbin, ss store.StoreService) (*Router, error) {
	return &Router{
		Router:       g,
		cb:           ca,
		storeService: ss,
	}, nil
}

// InitAccessingRouting 用户访问路由
func (r *Router) InitAccessingRouting() {
	// 访问请求通过jwt校验->casbin校验
	admin := r.Router.Group("/v1")
	{
		admin.POST("index", func(ctx *gin.Context) {
		})

	}
	admin.Use(JWT(r))
}
