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
				logger.Debugw("gin error ", "error", r)
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
			return
		}
		claims, err := util.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":  ErrorAuthCheckTokenFail,
				"status": http.StatusForbidden,
			})
			return
		} else if time.Now().Unix() > claims.ExpiresAt {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":  ErrorAuthCheckTokenExpired,
				"status": http.StatusBadRequest,
			})
			return
		}
		//TODO:
		//获取请求的资源
		namespace, _ := c.Get("namespace")
		source, _ := c.Get("source")
		action, _ := c.Get("action")
		enforce, err := r.cb.Enforcer.Enforce(claims.UserName, namespace.(string), source.(string), action.(string))
		if !enforce {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":  ErrorBadPermission,
				"status": http.StatusForbidden,
			})
			return
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
		var source, namespace string
		if action == "w" {
			action = "write"
		}
		if action == "r" {
			action = "read"
		}
		// ns name 为空，证明当前路径请求的是列表
		if c.Param("name") != "" {
			source = c.Param("name")
		} else {
			source = "*"
		}
		if c.Param("ns") != "" {
			namespace = c.Param("ns")
		} else if parts[3] == "policy" {
			namespace = "policy"
		} else {
			namespace = "*"
		}
		c.Set("version", version)
		c.Set("action", action)
		c.Set("source", source)
		c.Set("namespace", namespace)

		c.Next()
	}
}
