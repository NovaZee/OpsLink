package router

import (
	"github.com/denovo/permission/pkg/router/kubehandler"
	"github.com/denovo/permission/pkg/service"
	"github.com/denovo/permission/pkg/service/casbin"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

const (
	ErrorAuthCheckTokenFail    = " check fail "
	ErrorAuthCheckTokenExpired = " token expired "
	ErrorParamsError           = " bind params error  "
	ErrorBadPermission         = " Bad Permission！  "
)

type Router struct {
	Router *gin.Engine
	cb     *casbin.Casbin

	LoginHandler   []Front
	ManagerHandler []Management
}

func InitRouter(opslinkServer *service.OpsLinkServer) (*Router, error) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	engine := gin.Default()
	engine.Use(Logger())
	engine.GET("/validate")

	router, err := NewRouter(engine, opslinkServer)
	if err != nil {
		return nil, err
	}
	router.InitHandler(opslinkServer)
	registerFront(router, router.LoginHandler...)
	registerManager(router, router.ManagerHandler...)
	engine.Run(":" + opslinkServer.Config.Server.HttpPort).Error()
	return router, nil
}

func (r *Router) InitHandler(opslinkServer *service.OpsLinkServer) {
	handler := opslinkServer.K8sClient.K8sHandler
	front := []Front{
		BuildPublic(opslinkServer.Casbin, opslinkServer.StoreService),
	}
	in := []Management{
		//todo:考虑更加优雅的做法
		kubehandler.BuildRole(handler.RBACHandler),
		kubehandler.BuildNode(handler.NodeHandler),
		kubehandler.BuildDeployments(handler.DepHandler),
		kubehandler.BuildPod(handler.PodHandler),
		kubehandler.BuildConfigMap(handler.ConfigMapHandler),
		kubehandler.BuildService(handler.ServiceHandler),
		kubehandler.BuildNamespace(handler.NamespaceHandler),
		BuildRolePolicy(opslinkServer.Casbin, opslinkServer.StoreService),
	}
	r.LoginHandler = front
	r.ManagerHandler = in
}

// registerFront 注册前台路由（查看）
func registerFront(router *Router, handlers ...Front) {
	for _, h := range handlers {
		h.ReadRegister(router.Router.Group("v1/p"))
	}
}

// registerManager 注册管理路由
func registerManager(router *Router, handlers ...Management) {
	for _, h := range handlers {
		h.ReadRegister(router.Router.Group("/v1/r/" + h.GetName()).Use(router.ExtractParams()).Use(router.JWT()))
		h.WriteRegister(router.Router.Group("/v1/w/" + h.GetName()).Use(router.ExtractParams()).Use(router.JWT()))
	}
}

func NewRouter(g *gin.Engine, ops *service.OpsLinkServer) (*Router, error) {
	return &Router{
		Router: g,
		cb:     ops.Casbin,
	}, nil
}
