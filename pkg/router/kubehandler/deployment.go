package kubehandler

import (
	"github.com/denovo/permission/pkg/kubenates/kubeservice"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DeploymentController struct {
	DeploymentService *kubeservice.DeploymentService
	middlewares       []gin.HandlerFunc
}

func BuildDeployments(ds *kubeservice.DeploymentService, middleware ...gin.HandlerFunc) *DeploymentController {
	return &DeploymentController{
		DeploymentService: ds,
		middlewares:       middleware,
	}
}

func (dc *DeploymentController) List(ctx *gin.Context) {
	namespace := ctx.DefaultQuery("namespace", "default")
	res, err := dc.DeploymentService.List(ctx, namespace)
	if err != nil {
		return
	}
	// 配合前端
	ctx.JSON(http.StatusOK, gin.H{"data": res, "status": http.StatusOK})
	return
}
func (dc *DeploymentController) Delete(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")
	dc.DeploymentService.Delete(ctx, ns, name)
}

// Register 实现deployment controller 路由 框架规范
func (dc *DeploymentController) Register(g *gin.Engine) {

	deployments := g.Group("v1/deployments").Use(dc.middlewares...)
	{
		deployments.GET("list", func(ctx *gin.Context) { dc.List(ctx) }) ///deployments?namespace=
		deployments.POST("delete/:ns/:name", func(ctx *gin.Context) { dc.Delete(ctx) })
	}
}
