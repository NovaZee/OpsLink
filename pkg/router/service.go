package router

import (
	"github.com/denovo/permission/pkg/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AddPolicy 新增权限策略
func AddPolicy(ctx *gin.Context, c *casbin.Casbin) {
	sub := ctx.Query("sub")
	obj := ctx.Query("obj")
	method := ctx.Query("method")
	if len(sub) == 0 || len(obj) == 0 || len(method) == 0 {
		ctx.JSONP(http.StatusBadRequest, gin.H{"message": "", "status": http.StatusBadRequest})
		return
	}
	casbinModel := casbin.NewCasbinModel(sub, obj, method)
	add := c.DefaultPolicy.Add(casbinModel)
	if add == false {
		ctx.JSONP(http.StatusOK, gin.H{"message": "添加重复", "status": http.StatusOK})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"message": "添加成功", "status": http.StatusOK})
	return
}

// DeletePolicy  删除权限策略
func DeletePolicy(ctx *gin.Context, c *casbin.Casbin) {
	sub := ctx.Query("sub")
	obj := ctx.Query("obj")
	method := ctx.Query("method")
	if len(sub) == 0 || len(obj) == 0 || len(method) == 0 {
		ctx.JSONP(http.StatusBadRequest, gin.H{"message": "", "status": http.StatusBadRequest})
		return
	}
	casbinModel := casbin.NewCasbinModel(sub, obj, method)
	add := c.DefaultPolicy.Delete(casbinModel)
	if add == false {
		ctx.JSONP(http.StatusOK, gin.H{"message": "删除失败", "status": http.StatusOK})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"message": "删除成功", "status": http.StatusOK})
	return
}
