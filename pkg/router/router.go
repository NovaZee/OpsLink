package router

import (
	"errors"
	config "github.com/denovo/permission/configration"
	"github.com/denovo/permission/pkg/casbin"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type Router struct {
	router *gin.Engine
	cb     *casbin.Casbin
}

func InitRouter(ca *casbin.Casbin, conf *config.Config) error {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	engine := gin.Default()
	engine.Use(Logger())
	engine.Run(conf.Server.HttpPort)
	router, err := NewRouter(engine, ca)
	if err != nil {
		return err
	}
	router.InitRouting()
	s := engine.Run(":" + conf.Server.HttpPort).Error()
	if len(s) != 0 {
		return errors.New(s)
	}
	return nil
}

func NewRouter(g *gin.Engine, ca *casbin.Casbin) (*Router, error) {
	return &Router{
		router: g,
		cb:     ca,
	}, nil
}

func (r *Router) InitRouting() {
	v1 := r.router.Group("/manager", ManagerMiddleware())
	{
		v1.POST("addPolicy", func(ctx *gin.Context) {
			AddPolicy(ctx, r.cb)
		})
		v1.GET("deletePolicy", func(ctx *gin.Context) {
			DeletePolicy(ctx, r.cb)
		})
		v1.POST("update", func(ctx *gin.Context) {
		})
		v1.GET("listRoom")
	}

}
