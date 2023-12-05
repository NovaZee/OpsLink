package router

import "github.com/gin-gonic/gin"

// Handler 接口定义了处理程序的方法
type Handler interface {
	Register(g *gin.Engine)
}
