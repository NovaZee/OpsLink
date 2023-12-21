package router

import (
	"fmt"
	"github.com/denovo/permission/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/oppslink/protocol/logger"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

// Logger Logger中间件 集成到自己的日志库
func Logger() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 执行时间
		nowTime := time.Now()
		logger.Debugw(" http request",
			zap.Any(" request", c.Request.URL),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", time.Since(nowTime)),
		)
		c.Next()
	}
}

// JWT token验证中间件
func (r *Router) JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Debugw("jwt check ", "error", r)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s", r), "status": http.StatusInternalServerError})
				return
			}
		}()
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":  ErrorAuthCheckTokenFail,
				"status": http.StatusForbidden,
			})
		}
		claims, err := util.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":  ErrorAuthCheckTokenFail,
				"status": http.StatusForbidden,
			})
		} else if time.Now().Unix() > claims.ExpiresAt {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":  ErrorAuthCheckTokenExpired,
				"status": http.StatusBadRequest,
			})
		}
		//TODO:
		//获取请求的资源
		domain, _ := c.Get("domain")
		resource, _ := c.Get("resource")
		action, _ := c.Get("action")
		_, _ = c.Get("version")
		println(claims.UserName, domain.(string), resource.(string), action.(string))
		enforce := r.cb.Enforcer.Enforce(claims.UserName, domain.(string), resource.(string), action.(string))
		if !enforce {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":  ErrorAuthCheckTokenFail,
				"status": http.StatusForbidden,
			})
		}
		c.Next()
	}
}

//list资源列表单独设置中间件

// ExtractParams 对于read或者write的中间件权限控制
func (r *Router) ExtractParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从路径中解析域、资源和操作
		path := c.FullPath()
		parts := strings.Split(path, "/")
		version := parts[1]
		action := parts[2]
		var resource, domain string
		if action == "w" {
			action = "write"
		}
		if action == "r" {
			action = "read"
		}
		// ns name 为空，证明当前路径请求的是列表
		if c.Param("name") != "" {
			resource = c.Param("name")
		} else {
			resource = "*"
		}
		if c.Param("ns") != "" {
			domain = c.Param("ns")
		} else {
			domain = "*"
		}
		c.Set("version", version)
		c.Set("action", action)
		c.Set("resource", resource)
		c.Set("domain", domain)

		c.Next()
	}
}
