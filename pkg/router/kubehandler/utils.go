package kubehandler

import "github.com/gin-gonic/gin"

func KubeErrorResponse(ctx *gin.Context, statusCode int, err error) {
	ctx.JSON(statusCode, gin.H{"error": err.Error(), "status": statusCode})
	return
}
