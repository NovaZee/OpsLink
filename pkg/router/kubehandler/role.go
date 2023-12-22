package kubehandler

import (
	"github.com/denovo/permission/pkg/service/kubenates/kubeservice"
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

// GetName 实现deployment controller 路由 框架规范
func (rc *RoleController) GetName() string {
	return "role"
}

// ReadRegister 实现deployment controller 路由 框架规范
func (rc *RoleController) ReadRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	rbac := g.Use(middle...)
	{
		rbac.GET("", func(ctx *gin.Context) { rc.ListTarget(ctx) })
		rbac.GET("/sa", func(ctx *gin.Context) { rc.ListSaTarget(ctx) })
	}
}

// WriteRegister 实现deployment controller 路由 框架规范
func (rc *RoleController) WriteRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
}
