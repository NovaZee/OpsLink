package kubehandler

import (
	"github.com/denovo/permission/pkg/service/kubenates/kubeservice"
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
	ctx.Header("Content-Disposition", "attachment; filename="+name+".yaml")
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
	err = sc.ServiceService.ApplyByYaml(ctx, ns, data, false)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": true, "status": http.StatusOK})
	return
}

func (sc *ServiceController) update(ctx *gin.Context) {
	ns := ctx.Param("ns")
	_ = ctx.Param("name")
	all, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	if err != nil {
		KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	err = sc.ServiceService.ApplyByYaml(ctx, ns, all, true)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	KubeSuccessResponse(ctx, http.StatusOK)
	return
}

// GetName 实现deployment controller 路由 框架规范
func (sc *ServiceController) GetName() string {
	return "service"
}

// ReadRegister 实现deployment controller 路由 框架规范
func (sc *ServiceController) ReadRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	services := g.Use(middle...)
	{
		services.GET("", func(ctx *gin.Context) { sc.ListAll(ctx) })
		services.GET("/:ns/:name", func(ctx *gin.Context) { sc.Get(ctx) })
		services.GET("yaml/:ns/:name", func(ctx *gin.Context) { sc.downYaml(ctx) })
	}
}

// WriteRegister 实现deployment controller 路由 框架规范
func (sc *ServiceController) WriteRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	services := g.Use(middle...)
	{

		services.POST("apply/:ns", func(ctx *gin.Context) { sc.applyByYaml(ctx) })
		services.POST("upgrade/:ns", func(ctx *gin.Context) { sc.update(ctx) })
	}
}
