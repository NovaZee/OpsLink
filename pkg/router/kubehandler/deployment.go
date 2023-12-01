package kubehandler

import (
	"github.com/denovo/permission/pkg/kubenates/kubeservice"
	"github.com/gin-gonic/gin"
	"io"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
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
	res, err := dc.DeploymentService.List(namespace)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "status": http.StatusOK})
	return
}
func (dc *DeploymentController) Delete(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")
	err := dc.DeploymentService.Delete(ctx, ns, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"data": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": true, "status": http.StatusOK})
	return
}

func (dc *DeploymentController) DownYaml(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")

	yaml, err := dc.DeploymentService.DownToYaml(ns, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"data": err.Error(), "status": http.StatusInternalServerError})
	}
	// Set response headers for downloading the file
	ctx.Header("Content-Disposition", "attachment; filename=deployment.yaml")
	ctx.Header("Content-Type", "application/x-yaml")

	// Send the Deployment YAML as a response
	ctx.Data(http.StatusOK, "application/x-yaml", yaml)
	return
}

func (dc *DeploymentController) ApplyByYaml(ctx *gin.Context) {
	ns := ctx.Param("ns")
	// 从请求中获取上传的文件
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"data": err, "status": http.StatusInternalServerError})
		return
	}
	defer file.Close()
	// 读取上传的文件内容为二进制字节流
	data, err := io.ReadAll(file)

	// 创建一个 Unstructured 对象来装载 YAML 内容
	decode := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	deployment := &v1.Deployment{}
	_, _, err = decode.Decode(data, nil, deployment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"data": err, "status": http.StatusInternalServerError})
		return
	}

	err = dc.DeploymentService.Apply(ns, deployment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"data": err, "status": http.StatusInternalServerError})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": true, "status": http.StatusOK})
	return
}

// Register 实现deployment controller 路由 框架规范
func (dc *DeploymentController) Register(g *gin.Engine) {

	deployments := g.Group("v1/deployments").Use(dc.middlewares...)
	{
		deployments.GET("list", func(ctx *gin.Context) { dc.List(ctx) }) ///deployments?namespace=
		deployments.POST("delete/:ns/:name", func(ctx *gin.Context) { dc.Delete(ctx) })
		deployments.GET("download/:ns/:name", func(ctx *gin.Context) { dc.DownYaml(ctx) })
		deployments.POST("apply/:ns", func(ctx *gin.Context) { dc.ApplyByYaml(ctx) })
	}
}
