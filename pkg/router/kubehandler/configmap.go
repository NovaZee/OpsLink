package kubehandler

import (
	"github.com/denovo/permission/pkg/service/kubenates/kubeservice"
	"github.com/denovo/permission/protoc/kube"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type ConfigMapController struct {
	ConfigMapService *kubeservice.ConfigMapService
	middlewares      []gin.HandlerFunc
}

func BuildConfigMap(cms *kubeservice.ConfigMapService, middleware ...gin.HandlerFunc) *ConfigMapController {
	return &ConfigMapController{
		ConfigMapService: cms,
		middlewares:      middleware,
	}
}

func (cmc *ConfigMapController) ListAll(ctx *gin.Context) {
	ns := ctx.DefaultQuery("namespace", "default")
	res, err := cmc.ConfigMapService.ListByNamespace(ns)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
	}
	KubeSuccessMsgResponse(ctx, http.StatusOK, res)
	return
}
func (cmc *ConfigMapController) Get(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")
	res, err := cmc.ConfigMapService.GetConfigMap(ns, name)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
	}
	KubeSuccessMsgResponse(ctx, http.StatusOK, res)
	return
}
func (cmc *ConfigMapController) Apply(ctx *gin.Context) {
	cf := &kube.ConfigMap{}
	err := ctx.ShouldBindJSON(cf)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = cmc.ConfigMapService.Apply(ctx, cf)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	KubeSuccessResponse(ctx, http.StatusOK)
	return
}
func (cmc *ConfigMapController) delete(ctx *gin.Context) {
	ns := ctx.DefaultQuery("namespace", "default")
	name := ctx.DefaultQuery("name", "")
	err := cmc.ConfigMapService.Delete(ctx, ns, name)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	KubeSuccessResponse(ctx, http.StatusOK)
	return
}

func (cmc *ConfigMapController) applyByYaml(ctx *gin.Context) {
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
	err = cmc.ConfigMapService.ApplyByYaml(ctx, ns, data)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": true, "status": http.StatusOK})
	return
}

// GetName 实现deployment controller 路由 框架规范
func (cmc *ConfigMapController) GetName() string {
	return "configmap"
}

// ReadRegister 实现deployment controller 路由 框架规范
func (cmc *ConfigMapController) ReadRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	cm := g.Use(middle...)
	{
		cm.GET("", func(ctx *gin.Context) { cmc.ListAll(ctx) })
		cm.GET("/:ns/:name", func(ctx *gin.Context) { cmc.Get(ctx) })
	}
}

// WriteRegister 实现deployment controller 路由 框架规范
func (cmc *ConfigMapController) WriteRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	cm := g.Use(middle...)
	{
		cm.POST("", func(ctx *gin.Context) { cmc.Apply(ctx) })
		cm.DELETE("", func(ctx *gin.Context) { cmc.delete(ctx) })

		cm.POST("apply/:ns", func(ctx *gin.Context) { cmc.applyByYaml(ctx) })
	}
}
