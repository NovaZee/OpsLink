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

	PublicHandler  []FrontHandler
	FrontHandler   []FrontHandler
	ManagerHandler []ManagementHandler
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
	registerPublic(router, router.FrontHandler...)
	registerFront(router, router.FrontHandler...)
	registerManager(router, router.ManagerHandler...)
	engine.Run(":" + opslinkServer.Config.Server.HttpPort).Error()
	return router, nil
}

func (r *Router) InitHandler(opslinkServer *service.OpsLinkServer) {
	handler := opslinkServer.K8sClient.K8sHandler
	front := []FrontHandler{
		BuildRole(opslinkServer.Casbin, opslinkServer.StoreService),
	}
	in := []ManagementHandler{
		//todo:考虑更加优雅的做法
		kubehandler.BuildRole(handler.RBACHandler),
		kubehandler.BuildNode(handler.NodeHandler),
		kubehandler.BuildDeployments(handler.DepHandler),
		kubehandler.BuildPod(handler.PodHandler),
		kubehandler.BuildConfigMap(handler.ConfigMapHandler),
		kubehandler.BuildService(handler.ServiceHandler),
		kubehandler.BuildNamespace(handler.NamespaceHandler),
		BuildPolicy(opslinkServer.Casbin),
	}
	r.FrontHandler = front
	r.ManagerHandler = in
}

// registerHandlers 将多个处理程序注册到 Gin 路由器上
func registerFront(router *Router, handlers ...FrontHandler) {
	for _, h := range handlers {
		h.ReadRegister(router.Router.Group("v1/f"))
	}
}

// registerHandlers 将多个处理程序注册到 Gin 路由器上
func registerPublic(router *Router, handlers ...FrontHandler) {
	for _, h := range handlers {
		h.ReadRegister(router.Router.Group("/v1/p"))
	}
}

// registerHandlers 将多个处理程序注册到 Gin 路由器上
func registerManager(router *Router, handlers ...ManagementHandler) {
	for _, h := range handlers {
		use := router.Router.Group("/v1/r/" + h.GetName()).Use(router.ExtractParams()).Use(router.JWT())
		h.ReadRegister(use)
	}

	for _, h := range handlers {
		h.WriteRegister(router.Router.Group("/v1/w/" + h.GetName()).Use(router.ExtractParams()).Use(router.JWT()))
	}
}

func NewRouter(g *gin.Engine, cb *casbin.Casbin) (*Router, error) {
	return &Router{
		Router: g,
		cb:     cb,
	}, nil
}
