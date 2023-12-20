package router

import "github.com/gin-gonic/gin"

// Handler 接口定义了处理程序的方法
type Handler interface {
	ReadRegister(g *gin.RouterGroup, middle ...gin.HandlerFunc)
	WriteRegister(g *gin.RouterGroup, middle ...gin.HandlerFunc)
	GetName() string
}
