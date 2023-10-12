package router

import (
	"errors"
	"github.com/denovo/permission/pkg/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AddPolicy 新增权限策略
func AddPolicy(ctx *gin.Context, c *casbin.Casbin) {
	casbinModel, err := processManagerRequestParams(ctx)
	if err != nil {
		return
	}
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
	casbinModel, err := processManagerRequestParams(ctx)
	if err != nil {
		return
	}
	add := c.DefaultPolicy.Delete(casbinModel)
	if add == false {
		ctx.JSONP(http.StatusOK, gin.H{"message": "删除失败", "status": http.StatusOK})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"message": "删除成功", "status": http.StatusOK})
	return
}

// 共享的请求参数处理逻辑
func processManagerRequestParams(ctx *gin.Context) (*casbin.CasbinModel, error) {
	role := ctx.Query("role")
	source := ctx.Query("source")
	behavior := ctx.Query("behavior")
	if len(role) == 0 || len(source) == 0 || len(behavior) == 0 {
		ctx.JSONP(http.StatusBadRequest, gin.H{"message": "params errors", "status": http.StatusBadRequest})
		return nil, errors.New("params errors")
	}
	casbinModel := casbin.NewCasbinModel(role, source, behavior)
	return casbinModel, nil
}
