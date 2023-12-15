package kubehandler

import (
	"github.com/denovo/permission/pkg/kubenates/kubeservice"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type ServiceController struct {
	ServiceService *kubeservice.ServiceService
	middlewares    []gin.HandlerFunc
}

func BuildService(ss *kubeservice.ServiceService, middleware ...gin.HandlerFunc) *ServiceController {
	return &ServiceController{
		ServiceService: ss,
		middlewares:    middleware,
	}
}

func (sc *ServiceController) ListAll(ctx *gin.Context) {
	ns := ctx.DefaultQuery("namespace", "default")
	res, err := sc.ServiceService.ListServiceByNamespace(ns)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	KubeSuccessMsgResponse(ctx, 200, res)
	return
}

func (sc *ServiceController) Get(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")
	res, err := sc.ServiceService.Get(ns, name)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
	}
	KubeSuccessMsgResponse(ctx, http.StatusOK, res)
	return
}

func (sc *ServiceController) downYaml(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")

	yaml, err := sc.ServiceService.DownToYaml(ns, name)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	// Set response headers for downloading the file
	ctx.Header("Content-Disposition", "attachment; filename=deployment.yaml")
	ctx.Header("Content-Type", "application/x-yaml")

	// Send the Deployment YAML as a response
	ctx.Data(http.StatusOK, "application/x-yaml", yaml)
	return
}

func (sc *ServiceController) applyByYaml(ctx *gin.Context) {
	ns := ctx.Param("ns")
	// upload file
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	defer file.Close()
	// read to binary
	data, err := io.ReadAll(file)
	err = sc.ServiceService.ApplyByYaml(ctx, ns, data)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": true, "status": http.StatusOK})
	return
}

func (sc *ServiceController) Register(g *gin.Engine) {
	services := g.Group("v1/services").Use(sc.middlewares...)
	{
		services.GET("", func(ctx *gin.Context) { sc.ListAll(ctx) })
		services.GET("/:ns/:name", func(ctx *gin.Context) { sc.Get(ctx) })
		services.GET("downYaml/:ns/:name", func(ctx *gin.Context) { sc.downYaml(ctx) })
		services.POST("apply/:ns", func(ctx *gin.Context) { sc.applyByYaml(ctx) })
	}
}
