package router

import (
	"github.com/denovo/permission/pkg/router/kubehandler"
	"github.com/denovo/permission/pkg/service/casbin"
	"github.com/denovo/permission/pkg/service/store"
	"github.com/denovo/permission/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PolicyHandler struct {
	cb          *casbin.Casbin
	middlewares []gin.HandlerFunc

	ss store.StoreService
}

func BuildRolePolicy(cb *casbin.Casbin, ss store.StoreService, middleware ...gin.HandlerFunc) *PolicyHandler {
	return &PolicyHandler{
		cb:          cb,
		middlewares: middleware,
		ss:          ss,
	}
}

// AddPolicy 新增权限策略 -manager
func (ph *PolicyHandler) AddPolicy(ctx *gin.Context, c *casbin.Casbin) {
	casbinModel, err := kubehandler.ProcessManagerRequestParams(ctx)
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
	casbinModel, err := kubehandler.ProcessManagerRequestParams(ctx)
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
	casbinModel, err := kubehandler.ProcessManagerRequestParams(ctx)
	if err != nil {
		return
	}
	add := c.Update(casbinModel)
	if add == false {
		ctx.JSONP(http.StatusOK, gin.H{"message": "更新失败", "status": http.StatusOK})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"message": "更新成功", "status": http.StatusOK})
	return
}

func (ph *PolicyHandler) ListALlRoles(ctx *gin.Context) {
	pageNo := ctx.DefaultQuery("pageNo", "1")
	PageSize := ctx.DefaultQuery("PageSize", "10")
	num, err := strconv.Atoi(pageNo)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid parameter: id must be an integer"})
		return
	}
	if num < 1 {
		num = 1
	}
	size, err := strconv.Atoi(PageSize)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid parameter: id must be an integer"})
		return
	}
	role := ph.ss.GetRole()
	length := len(role.Roles)
	if length != 0 {
		start, end := util.PageSlice(length, num, size)
		kubehandler.KubeSuccessMsgResponse(ctx, http.StatusOK, role.Roles[start:end])
		return
	}
}

func (ph *PolicyHandler) GetName() string {
	return "policy"
}

// ReadRegister 实现deployment controller 路由 框架规范
func (ph *PolicyHandler) ReadRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	read := g.Use(middle...)
	{
		read.GET("listRole", func(ctx *gin.Context) { ph.ListALlRoles(ctx) })
	}
}

// WriteRegister 实现deployment controller 路由 框架规范
func (ph *PolicyHandler) WriteRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	admin := g.Use(middle...)
	{
		admin.POST("add", func(ctx *gin.Context) { ph.AddPolicy(ctx, ph.cb) })
		admin.GET("delete", func(ctx *gin.Context) { ph.DeletePolicy(ctx, ph.cb) })
		admin.POST("update", func(ctx *gin.Context) { ph.UpdatePolicy(ctx, ph.cb) })
	}
}
