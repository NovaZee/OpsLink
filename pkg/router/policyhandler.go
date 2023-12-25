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
func (ph *PolicyHandler) AddPolicy(ctx *gin.Context) {
	casbinModel, err := kubehandler.ProcessManagerRequestParams(ctx)
	if err != nil {
		return
	}
	add := ph.cb.Add(casbinModel)
	if add == false {
		ctx.JSONP(http.StatusOK, gin.H{"message": "添加重复", "status": http.StatusOK})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"message": "添加成功", "status": http.StatusOK})
	return
}

// DeletePolicy  删除权限策略 -manager
func (ph *PolicyHandler) DeletePolicy(ctx *gin.Context) {
	casbinModel, err := kubehandler.ProcessManagerRequestParams(ctx)
	if err != nil {
		return
	}
	add, err := ph.cb.Delete(casbinModel)
	if add == false || err != nil {
		ctx.JSONP(http.StatusOK, gin.H{"message": "删除失败", "status": http.StatusOK})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"message": "删除成功", "status": http.StatusOK})
	return
}

// UpdatePolicy  删除权限策略 -manager
func (ph *PolicyHandler) UpdatePolicy(ctx *gin.Context) {
	var req *casbin.UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	update, err := ph.cb.Update(req)
	if update == false || err != nil {
		ctx.JSONP(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
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
func (ph *PolicyHandler) ListUses(ctx *gin.Context) {

}
func (ph *PolicyHandler) ListRoles(ctx *gin.Context) {

}

func (ph *PolicyHandler) ListMyPolicy(ctx *gin.Context) {
	uname := ctx.Param("uname")
	//r2, err := ph.ss.Get(ctx, uname)
	//if err != nil || r2 == nil {
	//	ErrorResponse(ctx, http.StatusBadRequest, uname+" 不存在")
	//	return
	//}

	policy, err := ph.cb.ListMyPolicy(uname)
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, uname+" 不存在")
		return
	}
	kubehandler.KubeSuccessMsgResponse(ctx, http.StatusOK, policy)
	return
}

func (ph *PolicyHandler) GetName() string {
	return "policy"
}

// ReadRegister 实现deployment controller 路由 框架规范
func (ph *PolicyHandler) ReadRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	read := g.Use(middle...)
	{
		read.GET("policies/:uname", func(ctx *gin.Context) { ph.ListMyPolicy(ctx) })
	}
}

// WriteRegister 实现deployment controller 路由 框架规范
func (ph *PolicyHandler) WriteRegister(g gin.IRoutes, middle ...gin.HandlerFunc) {
	admin := g.Use(middle...)
	{
		admin.POST("add", func(ctx *gin.Context) { ph.AddPolicy(ctx) })
		admin.GET("delete", func(ctx *gin.Context) { ph.DeletePolicy(ctx) })
		admin.POST("update", func(ctx *gin.Context) { ph.UpdatePolicy(ctx) })

		admin.GET("listRole", func(ctx *gin.Context) { ph.ListALlRoles(ctx) })
		admin.GET("users", func(ctx *gin.Context) { ph.ListUses(ctx) })
		admin.POST("roles", func(ctx *gin.Context) { ph.ListRoles(ctx) })
	}
}
