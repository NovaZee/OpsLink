package router

import "github.com/gin-gonic/gin"

type ManagementHandler interface {
	FrontHandler
	WriteRegister(g gin.IRoutes, middle ...gin.HandlerFunc)
}

type FrontHandler interface {
	ReadRegister(g gin.IRoutes, middle ...gin.HandlerFunc)
	GetName() string
}
