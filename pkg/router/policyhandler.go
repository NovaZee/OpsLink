package router

import (
	"github.com/denovo/permission/pkg/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PolicyHandler struct {
	cb          *casbin.Casbin
	middlewares []gin.HandlerFunc
}

func BuildPolicy(cb *casbin.Casbin, middleware ...gin.HandlerFunc) *PolicyHandler {
	return &PolicyHandler{
		cb:          cb,
		middlewares: middleware,
	}
}

// AddPolicy 新增权限策略 -manager
func (ph *PolicyHandler) AddPolicy(ctx *gin.Context, c *casbin.Casbin) {
	casbinModel, err := processManagerRequestParams(ctx)
	if err != nil {
		return
	}
	add := c.Add(casbinModel)
	if add == false {
		ctx.JSONP(http.StatusOK, gin.H{"message": "添加重复", "status": http.StatusOK})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"message": "添加成功", "status": http.StatusOK})
	return
}

// DeletePolicy  删除权限策略 -manager
func (ph *PolicyHandler) DeletePolicy(ctx *gin.Context, c *casbin.Casbin) {
	casbinModel, err := processManagerRequestParams(ctx)
	if err != nil {
		return
	}
	add := c.Delete(casbinModel)
	if add == false {
		ctx.JSONP(http.StatusOK, gin.H{"message": "删除失败", "status": http.StatusOK})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"message": "删除成功", "status": http.StatusOK})
	return
}

// UpdatePolicy  删除权限策略 -manager
func (ph *PolicyHandler) UpdatePolicy(ctx *gin.Context, c *casbin.Casbin) {
	casbinModel, err := processManagerRequestParams(ctx)
	if err != nil {
		return
	}
	add := c.Delete(casbinModel)
	if add == false {
		ctx.JSONP(http.StatusOK, gin.H{"message": "删除失败", "status": http.StatusOK})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"message": "删除成功", "status": http.StatusOK})
	return
}

func (ph *PolicyHandler) GetName() string {
	return "policy"
}

// ReadRegister 实现deployment controller 路由 框架规范
func (ph *PolicyHandler) ReadRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {

}

// WriteRegister 实现deployment controller 路由 框架规范
func (ph *PolicyHandler) WriteRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	admin := g.Use(middle...)
	{
		admin.POST("addPolicy", func(ctx *gin.Context) { ph.AddPolicy(ctx, ph.cb) })
		admin.GET("deletePolicy", func(ctx *gin.Context) { ph.DeletePolicy(ctx, ph.cb) })
		admin.POST("update", func(ctx *gin.Context) {})
	}
}
