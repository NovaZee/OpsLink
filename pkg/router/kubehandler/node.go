package kubehandler

import (
	"github.com/denovo/permission/pkg/kubenates/kubeservice"
	"github.com/denovo/permission/protoc/kube"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"net/http"
)

type NodeController struct {
	NodeService *kubeservice.NodeService
	middlewares []gin.HandlerFunc
}

func BuildNode(ns *kubeservice.NodeService, middleware ...gin.HandlerFunc) *NodeController {
	return &NodeController{
		NodeService: ns,
		middlewares: middleware,
	}
}

func (nc *NodeController) List(ctx *gin.Context) {
	KubeSuccessMsgResponse(ctx, http.StatusOK, nc.NodeService.List(ctx))
	return
}

func (nc *NodeController) Modify(ctx *gin.Context) {
	node := &kube.FrontNode{}
	_ = ctx.ShouldBindJSON(node)
	get := nc.NodeService.Get(node.Name)
	if get == nil {
		KubeNotFoundResponse(ctx, http.StatusOK)
		return
	}
	for _, label := range node.Labels {
		get.Labels[label.Key] = label.Value
	}
	var ts []corev1.Taint
	for _, taint := range node.Taints {
		t := corev1.Taint{
			Key:    taint.Key,
			Value:  taint.Value,
			Effect: corev1.TaintEffect(taint.Effect)}
		ts = append(ts, t)
	}
	get.Spec.Taints = ts
	_, err := nc.NodeService.Update(ctx, get)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return

	}
	KubeSuccessResponse(ctx, http.StatusOK)
	return
}

// Register pod controller 路由 框架规范
func (nc *NodeController) Register(g *gin.Engine) {
	pods := g.Group("v1/nodes").Use(nc.middlewares...)
	{
		pods.GET("list", func(ctx *gin.Context) { nc.List(ctx) })
		pods.POST("modify", func(ctx *gin.Context) { nc.Modify(ctx) })
	}
}
