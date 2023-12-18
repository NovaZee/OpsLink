package kubehandler

import (
	"github.com/denovo/permission/pkg/kubenates/kubeservice"
	"github.com/gin-gonic/gin"
	v3yaml "gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
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

func (pc *PodController) GetFromApiServer(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")
	get, err := pc.PodService.GetDetail(ctx, ns, name)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	// 将 Deployment 对象转换为 Unstructured 对象
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(get)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	unstructuredObj["apiVersion"] = "v1"
	unstructuredObj["kind"] = "Pod"
	// 转换为 YAML 格式
	pod, err := v3yaml.Marshal(unstructuredObj)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	KubeSuccessYamlResponse(ctx, http.StatusOK, pod)
	return

}

// kubectl get pods -l app=nginx -n default
func (pc *PodController) getPodsByLabel(ctx *gin.Context) {
	ns := ctx.DefaultQuery("namespace", "default")
	label := ctx.DefaultQuery("label", "nginx")
	res, err := pc.PodService.GetByLabelInCache(ns, label)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	KubeSuccessMsgResponse(ctx, http.StatusOK, res)
	return
}

func (pc *PodController) downYaml(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")

	yaml, err := pc.PodService.DownToYaml(ns, name)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	// Set response headers for downloading the file
	ctx.Header("Content-Disposition", "attachment; filename="+name+".yaml")
	ctx.Header("Content-Type", "application/x-yaml")

	// Send the Deployment YAML as a response
	KubeSuccessYamlResponse(ctx, http.StatusOK, yaml)
	return
}

func (pc *PodController) GetFromCache(ctx *gin.Context) {
	_ = ctx.DefaultQuery("namespace", "default")
}

func (pc *PodController) List(ctx *gin.Context) {
	_ = ctx.DefaultQuery("namespace", "default")
}

func (pc *PodController) Delete(ctx *gin.Context) {
}

// Register pod controller 路由 框架规范
func (pc *PodController) Register(g *gin.Engine) {
	pods := g.Group("v1/pods").Use(pc.middlewares...)
	{
		pods.GET("list", func(ctx *gin.Context) { pc.List(ctx) })
		pods.GET("getDetail/:ns/:name", func(ctx *gin.Context) { pc.GetFromApiServer(ctx) })
		pods.GET("getPods", func(ctx *gin.Context) { pc.getPodsByLabel(ctx) })

		pods.GET("yaml/:ns/:name", func(ctx *gin.Context) { pc.downYaml(ctx) })
	}
}
