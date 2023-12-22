package kubehandler

import (
	"github.com/denovo/permission/pkg/service/kubenates/kubeservice"
	"github.com/gin-gonic/gin"
	v3yaml "gopkg.in/yaml.v3"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"strconv"
)

type DeploymentController struct {
	DeploymentService *kubeservice.DeploymentService
}

func BuildDeployments(ds *kubeservice.DeploymentService) *DeploymentController {
	return &DeploymentController{
		DeploymentService: ds,
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
		KubeSuccessResponse(ctx, http.StatusOK)
		return
	}
	deployment.Spec.Replicas = &valInt32

	_, err = dc.DeploymentService.Update(ctx, ns, name, deployment)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	KubeSuccessResponse(ctx, http.StatusOK)
	return
}

func (dc *DeploymentController) upgrade(ctx *gin.Context) {
	ns := ctx.Param("ns")
	_ = ctx.Param("name")
	deployment, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	err = dc.DeploymentService.ApplyByYaml(ctx, ns, deployment, true)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	KubeSuccessResponse(ctx, http.StatusOK)
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

func (dc *DeploymentController) Rollout(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")
	err := dc.DeploymentService.Rollout(ctx, ns, name)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": "success", "status": http.StatusOK})
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
	ctx.Header("Content-Disposition", "attachment; filename="+name+".yaml")
	ctx.Header("Content-Type", "application/x-yaml")

	// Send the Deployment YAML as a response
	ctx.Data(http.StatusOK, "application/x-yaml", yaml)
	return
}

func (dc *DeploymentController) applyByYaml(ctx *gin.Context) {
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
	err = dc.DeploymentService.ApplyByYaml(ctx, ns, data, false)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": true, "status": http.StatusOK})
	return
}

// MiddleHandler 实现deployment controller 路由 框架规范
func (dc *DeploymentController) MiddleHandler() string {
	return "deployments"
}

// GetName 实现deployment controller 路由 框架规范
func (dc *DeploymentController) GetName() string {
	return "deployments"
}

// ReadRegister 实现deployment controller 路由 框架规范
func (dc *DeploymentController) ReadRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	rd := g.Use(middle...)
	{
		rd.GET("list", func(ctx *gin.Context) { dc.list(ctx) })
		rd.GET("yaml/:ns/:name", func(ctx *gin.Context) { dc.downYaml(ctx) })
		rd.GET("checkout/:ns/:name", func(ctx *gin.Context) { dc.checkout(ctx) })
	}
}

// WriteRegister 实现deployment controller 路由 框架规范
func (dc *DeploymentController) WriteRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	wd := g.Use(middle...)
	{
		wd.POST("delete/:ns/:name", func(ctx *gin.Context) { dc.delete(ctx) })
		wd.POST("apply/:ns", func(ctx *gin.Context) { dc.applyByYaml(ctx) })
		wd.PUT("patch/:ns/:name", func(ctx *gin.Context) { dc.patch(ctx) })
		// deployment的所有更新操作
		wd.POST("upgrade/:ns/:name", func(ctx *gin.Context) { dc.upgrade(ctx) })
		wd.GET("rollout/:ns/:name", func(ctx *gin.Context) { dc.Rollout(ctx) })

		wd.PUT("scale/:ns/:name", func(ctx *gin.Context) { dc.scale(ctx) })
	}
}
