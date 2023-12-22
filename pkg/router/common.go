package router

import "github.com/gin-gonic/gin"

func ErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{"message": message, "status": statusCode})
}
