package router

import (
	"github.com/denovo/permission/pkg/casbin"
	"github.com/denovo/permission/pkg/router/kubehandler"
	"github.com/denovo/permission/pkg/service"
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

	handler []Handler
}

func InitRouter(opslinkServer *service.OpsLinkServer) (*Router, error) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	engine := gin.Default()
	engine.Use(Logger())
	engine.GET("/validate")

	router, err := NewRouter(engine, opslinkServer.Casbin)
	if err != nil {
		return nil, err
	}
	router.InitHandler(opslinkServer)
	registerHandlers(router, router.handler...)
	engine.Run(":" + opslinkServer.Config.Server.HttpPort).Error()
	return router, nil
}

func (r *Router) InitHandler(opslinkServer *service.OpsLinkServer) {
	handlers := []Handler{
		//BuildPolicy(opslinkServer.Casbin, ManagerMiddleware()),
		//BuildRole(opslinkServer.Casbin, opslinkServer.StoreService, ManagerMiddleware()),
		////todo:kube资源过多时，由于使用的路由中间件一致，可以继续抽离模块，尽量避免在路由模块操作
		//kubehandler.BuildRole(opslinkServer.K8sClient.RBACHandler, JWT(opslinkServer.Casbin)),
		//kubehandler.BuildNode(opslinkServer.K8sClient.NodeHandler, JWT(opslinkServer.Casbin)),
		kubehandler.BuildDeployments(opslinkServer.K8sClient.DepHandler),
		//kubehandler.BuildPod(opslinkServer.K8sClient.PodHandler, JWT(opslinkServer.Casbin)),
		//kubehandler.BuildConfigMap(opslinkServer.K8sClient.ConfigMapHandler, JWT(opslinkServer.Casbin)),
		//kubehandler.BuildService(opslinkServer.K8sClient.ServiceHandler, JWT(opslinkServer.Casbin)),
		//kubehandler.BuildNamespace(opslinkServer.K8sClient.NamespaceHandler, JWT(opslinkServer.Casbin)),
	}
	r.handler = handlers
}

// registerHandlers 将多个处理程序注册到 Gin 路由器上
func registerHandlers(router *Router, handlers ...Handler) {
	for _, h := range handlers {
		rGroup := router.Router.Group("/v1/r/" + h.GetName())
		h.ReadRegister(rGroup)
	}

	for _, h := range handlers {
		wGroup := router.Router.Group("/v1/w/" + h.GetName())
		h.WriteRegister(wGroup)
	}
}
func NewRouter(g *gin.Engine, cb *casbin.Casbin) (*Router, error) {
	return &Router{
		Router: g,
		cb:     cb,
	}, nil
}

//// InitAccessingRouting 用户访问路由
//func (r *Router) InitAccessingRouting() {
//	// 访问请求通过jwt校验->casbin校验
//	admin := r.Router.Group("/v1")
//	{
//		admin.POST("index", func(ctx *gin.Context) {
//		})
//
//	}
//	admin.Use(JWT(r.cb))
//}
