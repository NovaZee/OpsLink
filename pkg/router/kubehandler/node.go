package kubehandler

import (
	"github.com/denovo/permission/pkg/kubenates/kubeservice"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NodeController struct {
	NodeService *kubeservice.NodeService
	middlewares []gin.HandlerFunc
}

func BuildNode(ns *kubeservice.NodeService, middleware ...gin.HandlerFunc) *NodeController {
	return &NodeController{
		NodeService: ns,
		middlewares: middleware,
	}
}

func (nc *NodeController) List(ctx *gin.Context) {
	KubeSuccessResponse(ctx, http.StatusOK, nc.NodeService.List(ctx))
	return
}

// Register pod controller 路由 框架规范
func (nc *NodeController) Register(g *gin.Engine) {
	pods := g.Group("v1/nodes").Use(nc.middlewares...)
	{
		pods.GET("list", func(ctx *gin.Context) { nc.List(ctx) })
	}
}
