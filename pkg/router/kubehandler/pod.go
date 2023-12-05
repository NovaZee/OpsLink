package kubehandler

import (
	"github.com/denovo/permission/pkg/kubenates/kubeservice"
	"github.com/gin-gonic/gin"
)

type PodController struct {
	PodService  *kubeservice.PodService
	middlewares []gin.HandlerFunc
}

func BuildPod(ps *kubeservice.PodService, middleware ...gin.HandlerFunc) *PodController {
	return &PodController{
		PodService:  ps,
		middlewares: middleware,
	}
}

func (dc *PodController) List(ctx *gin.Context) {
	_ = ctx.DefaultQuery("namespace", "default")
}
func (dc *PodController) Delete(ctx *gin.Context) {
}

// Register pod controller 路由 框架规范
func (dc *PodController) Register(g *gin.Engine) {
	pods := g.Group("v1/pods").Use(dc.middlewares...)
	{
		pods.GET("list", func(ctx *gin.Context) { dc.List(ctx) }) ///deployments?namespace=

	}
}
