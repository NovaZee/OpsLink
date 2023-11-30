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
}
func (dc *PodController) Delete(ctx *gin.Context) {
}

// Register pod controller 路由 框架规范
func (dc *PodController) Register(g *gin.Engine) {
}
