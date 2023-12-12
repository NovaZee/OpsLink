package kubehandler

import (
	"github.com/denovo/permission/pkg/kubenates/kubeservice"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NamespaceController struct {
	NamespaceService *kubeservice.NamespaceService
	middlewares      []gin.HandlerFunc
}

func BuildNamespace(nss *kubeservice.NamespaceService, middleware ...gin.HandlerFunc) *NamespaceController {
	return &NamespaceController{
		NamespaceService: nss,
		middlewares:      middleware,
	}
}

func (nss *NamespaceController) List(ctx *gin.Context) {
	KubeSuccessResponse(ctx, http.StatusOK, nss.NamespaceService.List())
	return
}

func (nss *NamespaceController) Add(ctx *gin.Context) {
	name := ctx.Param("name")
	create, err := nss.NamespaceService.Create(ctx, name)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	KubeSuccessResponse(ctx, http.StatusOK, create)
	return
}

func (nss *NamespaceController) Remove(ctx *gin.Context) {
	name := ctx.Param("name")
	err := nss.NamespaceService.Remove(ctx, name)
	if err != nil {
		KubeErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	KubeSuccessResponse(ctx, http.StatusOK, name)
	return
}

// Register NamespaceController 路由 框架规范
func (nss *NamespaceController) Register(g *gin.Engine) {

	ns := g.Group("v1/namespace").Use(nss.middlewares...)
	{
		ns.GET("list", func(ctx *gin.Context) { nss.List(ctx) })
		ns.PUT("add/:name", func(ctx *gin.Context) { nss.Add(ctx) })
		ns.PUT("remove/:name", func(ctx *gin.Context) { nss.Remove(ctx) })
	}
}
