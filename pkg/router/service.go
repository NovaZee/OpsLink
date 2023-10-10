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
	policy.Add(casbinModel)
}
