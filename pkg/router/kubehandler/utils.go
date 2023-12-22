package kubehandler

import (
	"errors"
	"github.com/denovo/permission/pkg/service/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
)

func KubeErrorResponse(ctx *gin.Context, statusCode int, err error) {
	ctx.JSON(statusCode, gin.H{"error": err.Error(), "status": statusCode})
	return
}

func KubeSuccessMsgResponse(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, gin.H{"data": data, "status": statusCode})
	return
}
func KubeSuccessResponse(ctx *gin.Context, statusCode int) {
	ctx.JSON(statusCode, gin.H{"data": "success", "status": statusCode})
	return
}

func KubeNotFoundResponse(ctx *gin.Context, statusCode int) {
	ctx.JSON(statusCode, gin.H{"data": "not exist！", "status": statusCode})
	return
}

func KubeSuccessYamlResponse(ctx *gin.Context, statusCode int, out []byte) {
	ctx.Data(statusCode, "application/x-yaml", out)
	return
}

// ProcessManagerRequestParams  共享的请求参数处理逻辑 -manager
func ProcessManagerRequestParams(ctx *gin.Context) (*casbin.CasbinModel, error) {
	role := ctx.Query("role")
	domain := ctx.Query("domain")
	pType := ctx.Query("type")
	if len(role) == 0 || len(domain) == 0 || len(pType) == 0 {
		ctx.JSONP(http.StatusBadRequest, gin.H{"message": "params errors", "status": http.StatusBadRequest})
		return nil, errors.New("params errors")
	}
	if pType == "p" {
		source := ctx.Query("source")
		behavior := ctx.Query("behavior")
		if len(source) == 0 || len(behavior) == 0 {
			ctx.JSONP(http.StatusBadRequest, gin.H{"message": "params errors", "status": http.StatusBadRequest})
			return nil, errors.New("params errors")
		}
	}
	source := ctx.Query("source")
	behavior := ctx.Query("behavior")
	if len(role) == 0 || len(source) == 0 || len(behavior) == 0 || len(domain) == 0 || len(pType) == 0 {
		ctx.JSONP(http.StatusBadRequest, gin.H{"message": "params errors", "status": http.StatusBadRequest})
		return nil, errors.New("params errors")
	}
	casbinModel := casbin.NewCasbinModel(pType, role, domain, source, behavior)
	return casbinModel, nil
}

func ErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{"message": message, "status": statusCode})
}
