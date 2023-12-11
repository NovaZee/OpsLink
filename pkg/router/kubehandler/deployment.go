package kubehandler

import (
	"github.com/denovo/permission/pkg/kubenates/kubeservice"
	"github.com/gin-gonic/gin"
	v3yaml "gopkg.in/yaml.v3"
	"io"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"strconv"
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

// checkout 默认从
func (dc *DeploymentController) checkout(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")
	yaml, err := dc.DeploymentService.DownToYaml(ns, name)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Data(http.StatusOK, "application/x-yaml", yaml)
	return
}

func (dc *DeploymentController) checkoutFromApiServer(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")
	get, err := dc.DeploymentService.Client.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
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
	unstructuredObj["apiVersion"] = "apps/v1"
	unstructuredObj["kind"] = "Deployment"
	// 转换为 YAML 格式
	deploymentByte, err := v3yaml.Marshal(unstructuredObj)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Data(http.StatusOK, "application/x-yaml", deploymentByte)
	return
}

func (dc *DeploymentController) list(ctx *gin.Context) {
	namespace := ctx.DefaultQuery("namespace", "default")
	res, err := dc.DeploymentService.List(namespace)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res, "status": http.StatusOK})
	return
}

func (dc *DeploymentController) patch(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")
	all, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//var result map[string]interface{}
	//_ = json.Unmarshal(all, result)
	_, err = dc.DeploymentService.Patch(ctx, ns, name, all)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	return
}

func (dc *DeploymentController) scale(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")
	scale := ctx.Query("scale")
	deployment, err := dc.DeploymentService.GetDeployment(ns, name)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	atoi, err := strconv.ParseInt(scale, 10, 32)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	valInt32 := int32(atoi)
	replicas := *deployment.Spec.Replicas
	if valInt32 == replicas {
		KubeSuccessResponse(ctx, http.StatusOK, true)
		return
	}
	deployment.Spec.Replicas = &valInt32

	_, err = dc.DeploymentService.Update(ctx, ns, name, deployment)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	KubeSuccessResponse(ctx, http.StatusOK, true)
	return
}

func (dc *DeploymentController) update(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")
	all, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	deployment := &v1.Deployment{}
	err = json.Unmarshal(all, deployment)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	_, err = dc.DeploymentService.Update(ctx, ns, name, deployment)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	return
}

func (dc *DeploymentController) delete(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")
	err := dc.DeploymentService.Delete(ctx, ns, name)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": true, "status": http.StatusOK})
	return
}

func (dc *DeploymentController) downYaml(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")

	yaml, err := dc.DeploymentService.DownToYaml(ns, name)
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

func (dc *DeploymentController) applyByYaml(ctx *gin.Context) {
	ns := ctx.Param("ns")
	// 从请求中获取上传的文件
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
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
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = dc.DeploymentService.Apply(ns, deployment)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": true, "status": http.StatusOK})
	return
}

// Register 实现deployment controller 路由 框架规范
func (dc *DeploymentController) Register(g *gin.Engine) {

	deployments := g.Group("v1/deployments").Use(dc.middlewares...)
	{
		deployments.GET("list", func(ctx *gin.Context) { dc.list(ctx) })
		deployments.POST("delete/:ns/:name", func(ctx *gin.Context) { dc.delete(ctx) })
		deployments.GET("download/:ns/:name", func(ctx *gin.Context) { dc.downYaml(ctx) })
		deployments.POST("apply/:ns", func(ctx *gin.Context) { dc.applyByYaml(ctx) })
		deployments.PUT("patch/:ns/:name", func(ctx *gin.Context) { dc.patch(ctx) })
		deployments.POST("update/:ns/:name", func(ctx *gin.Context) { dc.update(ctx) })
		deployments.GET("checkout/:ns/:name", func(ctx *gin.Context) { dc.checkout(ctx) })

		deployments.PUT("scale/:ns/:name", func(ctx *gin.Context) { dc.scale(ctx) })
	}
}
