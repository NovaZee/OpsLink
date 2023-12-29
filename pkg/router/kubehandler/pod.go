package kubehandler

import (
	"context"
	"github.com/denovo/permission/pkg/service/kubenates/kubeservice"
	"github.com/gin-gonic/gin"
	v3yaml "gopkg.in/yaml.v3"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"time"
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

func (pc *PodController) GetLogs(ctx *gin.Context) {
	ns := ctx.DefaultQuery("ns", "default")
	podName := ctx.DefaultQuery("pname", "")
	cname := ctx.DefaultQuery("cname", "")
	var tailLine int64 = 100
	opt := &v1.PodLogOptions{
		Follow:    true,
		Container: cname,
		TailLines: &tailLine,
	}

	cc, _ := context.WithTimeout(ctx, time.Minute*30) //设置半小时超时时间。否则会造成内存泄露
	req := pc.PodService.Client.CoreV1().Pods(ns).GetLogs(podName, opt)
	reader, _ := req.Stream(cc)
	defer reader.Close()

	// 分块发送的方式
	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf) // 如果 当前日志 读完了。 会阻塞

		if err != nil && err != io.EOF { //一旦超时 会进入 这个程序 ,,此时一定要break 掉
			break
		}

		w, err := ctx.Writer.Write([]byte(string(buf[0:n])))
		if w == 0 || err != nil {
			break
		}
		ctx.Writer.(http.Flusher).Flush()
	}

	return
}

func (pc *PodController) Pods(ctx *gin.Context) {
	ns := ctx.DefaultQuery("namespace", "default")
	KubeSuccessMsgResponse(ctx, http.StatusOK, pc.PodService.ListByNamespace(ns))
	return
}

// GetName 实现deployment controller 路由 框架规范
func (pc *PodController) GetName() string {
	return "pod"
}

// ReadRegister 实现deployment controller 路由 框架规范
func (pc *PodController) ReadRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	pods := g.Use(middle...)
	{
		pods.GET("list", func(ctx *gin.Context) { pc.List(ctx) })
		pods.GET("getDetail/:ns/:name", func(ctx *gin.Context) { pc.GetFromApiServer(ctx) })
		pods.GET("getPods", func(ctx *gin.Context) { pc.getPodsByLabel(ctx) })
		pods.GET("pods", func(ctx *gin.Context) { pc.Pods(ctx) })
		pods.GET("yaml/:ns/:name", func(ctx *gin.Context) { pc.downYaml(ctx) })
		pods.GET("logs", func(ctx *gin.Context) { pc.GetLogs(ctx) })
	}
}

// WriteRegister 实现deployment controller 路由 框架规范
func (pc *PodController) WriteRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {

}
