package router

import "github.com/gin-gonic/gin"

type Management interface {
	Front
	WriteRegister(g gin.IRoutes, middle ...gin.HandlerFunc)
}

type Front interface {
	ReadRegister(g gin.IRoutes, middle ...gin.HandlerFunc)
	GetName() string
}
