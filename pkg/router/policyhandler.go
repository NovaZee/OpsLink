package router

import (
	"github.com/denovo/permission/pkg/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PolicyHandler struct {
}

func BuildPolicy() *PolicyHandler {
	return &PolicyHandler{}
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

func (ph *PolicyHandler) Register(r *Router) {
	admin := r.Router.Group("/manager")
	{
		admin.POST("addPolicy", func(ctx *gin.Context) { ph.AddPolicy(ctx, r.cb) })
		admin.GET("deletePolicy", func(ctx *gin.Context) { ph.DeletePolicy(ctx, r.cb) })
		admin.POST("update", func(ctx *gin.Context) {})
	}
	//admin.Use(ManagerMiddleware())
}
