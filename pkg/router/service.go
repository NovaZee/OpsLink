package router

import (
	"github.com/denovo/permission/pkg/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddPolicy(ctx *gin.Context, adapter *casbin.CasbinAdapter) {
	sub := ctx.Query("sub")
	obj := ctx.Query("obj")
	method := ctx.Query("method")
	if len(sub) == 0 || len(obj) == 0 || len(method) == 0 {
		ctx.JSONP(http.StatusBadRequest, gin.H{"message": "", "status": http.StatusBadRequest})
		return
	}
	policy, err := casbin.NewDefaultPolicy(adapter)
	if err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"message": "", "status": http.StatusBadRequest})
		return
	}
	casbinModel := casbin.NewCasbinModel(sub, obj, method)
	add := policy.Add(casbinModel)
	if add == false {
		ctx.JSONP(http.StatusOK, gin.H{"message": "添加重复", "status": http.StatusOK})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"message": "添加成功", "status": http.StatusOK})
	return
}

func DeletePolicy(ctx *gin.Context, adapter *casbin.CasbinAdapter) {
	sub := ctx.Query("sub")
	obj := ctx.Query("obj")
	method := ctx.Query("method")
	if len(sub) == 0 || len(obj) == 0 || len(method) == 0 {
		ctx.JSONP(http.StatusBadRequest, gin.H{"message": "", "status": http.StatusBadRequest})
		return
	}
	policy, err := casbin.NewDefaultPolicy(adapter)
	if err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"message": "", "status": http.StatusBadRequest})
		return
	}
	casbinModel := casbin.NewCasbinModel(sub, obj, method)
	add := policy.Delete(casbinModel)
	if add == false {
		ctx.JSONP(http.StatusOK, gin.H{"message": "删除失败", "status": http.StatusOK})
		return
	}
	ctx.JSONP(http.StatusOK, gin.H{"message": "删除成功", "status": http.StatusOK})
	return
}
