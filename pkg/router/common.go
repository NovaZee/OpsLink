package router

import (
	"errors"
	"github.com/denovo/permission/pkg/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 共享的请求参数处理逻辑 -manager
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

func ErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{"message": message, "status": statusCode})
}
