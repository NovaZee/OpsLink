package router

import (
	"github.com/denovo/permission/pkg/casbin"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type Router struct {
	router *gin.Engine
	csba   *casbin.CasbinAdapter
}

func InitRouter(ca *casbin.CasbinAdapter) (*Router, error) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	engine := gin.Default()
	engine.Use(Logger())
	defer engine.Run(":8081")
	router, err := NewRouter(engine, ca)
	router.InitRouting()
	return router, err
}

func NewRouter(g *gin.Engine, ca *casbin.CasbinAdapter) (*Router, error) {
	return &Router{
		router: g,
		csba:   ca,
	}, nil
}

func (r *Router) InitRouting() {

	v1 := r.router.Group("/v1")
	{
		v1.POST("addPolicy", func(ctx *gin.Context) {
			AddPolicy(ctx, r.csba)
		})
		v1.GET("listParticipants", func(ctx *gin.Context) {
		})
		v1.POST("update", func(ctx *gin.Context) {
		})
		v1.GET("listRoom")
	}

	manager := r.router.Group("/manager")
	{

		manager.POST("addPolicy", func(ctx *gin.Context) {
			AddPolicy(ctx, r.csba)
		})
		manager.GET("listParticipants", func(ctx *gin.Context) {
		})
		manager.POST("update", func(ctx *gin.Context) {
		})
		manager.GET("listRoom")
	}
}
