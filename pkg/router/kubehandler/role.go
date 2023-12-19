package kubehandler

import (
	"github.com/denovo/permission/pkg/kubenates/kubeservice"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RoleController struct {
	rs          *kubeservice.RBACService
	middlewares []gin.HandlerFunc
}

func BuildRole(rs *kubeservice.RBACService, middleware ...gin.HandlerFunc) *RoleController {
	return &RoleController{
		rs:          rs,
		middlewares: middleware,
	}
}

func (rc *RoleController) ListTarget(ctx *gin.Context) {
	ns := ctx.DefaultQuery("ns", "default")
	KubeSuccessMsgResponse(ctx, http.StatusOK, rc.rs.ListRoles(ns))
	return
}
func (rc *RoleController) ListSaTarget(ctx *gin.Context) {
	ns := ctx.DefaultQuery("ns", "default")
	KubeSuccessMsgResponse(ctx, http.StatusOK, rc.rs.ListSa(ns))
	return
}

func (rc *RoleController) Register(g *gin.Engine) {
	rbac := g.Group("v1/roles").Use(rc.middlewares...)
	{
		rbac.GET("", func(ctx *gin.Context) { rc.ListTarget(ctx) })
		rbac.GET("/sa", func(ctx *gin.Context) { rc.ListSaTarget(ctx) })
	}
}
