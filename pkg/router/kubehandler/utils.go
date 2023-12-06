package kubehandler

import (
	"github.com/gin-gonic/gin"
)

func KubeErrorResponse(ctx *gin.Context, statusCode int, err error) {
	ctx.JSON(statusCode, gin.H{"error": err.Error(), "status": statusCode})
	return
}

func KubeSuccessResponse(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, gin.H{"data": data, "status": statusCode})
	return
}

func KubeSuccessYamlResponse(ctx *gin.Context, statusCode int, out []byte) {
	ctx.Data(statusCode, "application/x-yaml", out)
	return
}
