package router

import (
	"github.com/denovo/permission/pkg/kubenates"
	"github.com/gin-gonic/gin"
	"github.com/oppslink/protocol/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentController struct {
	client            kubernetes.Interface
	DeploymentService *kubenates.DeploymentHandler
}

func BuildDeployments(client kubernetes.Interface, dh *kubenates.DeploymentHandler) *DeploymentController {
	return &DeploymentController{
		client:            client,
		DeploymentService: dh,
	}
}

func (dc *DeploymentController) List(ctx *gin.Context) {
	get, err := dc.client.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		println(err)
	}
	for _, item := range get.Items {
		logger.Infow("kubenates-system", "Namespace", item.Namespace, "Name", item.GetName())
	}
}
func (dc *DeploymentController) Delete(ctx *gin.Context) {
	ns := ctx.Param("ns")
	name := ctx.Param("name")

	err := dc.client.AppsV1().Deployments(ns).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		logger.Errorw("DeploymentController Delete ", err)
	}
	logger.Infow("DeploymentController Delete ", "Name", name, "Namespace", ns)
}

// Register 实现deployment controller 路由 框架规范
func (dc *DeploymentController) Register(r *Router) {
	deployments := r.Router.Group("v1/deployments").Use(JWT(r))
	{
		deployments.POST("list", func(ctx *gin.Context) { dc.List(ctx) })
		deployments.POST("delete/:ns/:name", func(ctx *gin.Context) { dc.Delete(ctx) })
	}
}
