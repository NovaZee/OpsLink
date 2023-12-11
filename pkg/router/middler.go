package router

import (
	"github.com/denovo/permission/pkg/casbin"
	"github.com/denovo/permission/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/oppslink/protocol/logger"
	"go.uber.org/zap"
	"net/http"
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
func JWT(cb *casbin.Casbin) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Debugw("jwt check ", "error", r, "http errorCode", c.Writer.Status())
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
		enforce := cb.Enforcer.Enforce(claims.UserName, casbin.HttpV1, casbin.Read)
		if !enforce {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":  ErrorAuthCheckTokenFail,
				"status": http.StatusForbidden,
			})
		}
		c.Next()
	}
}

// ManagerMiddleware 管理Group路径中间件  /manager
func ManagerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Query("behavior") {
		case casbin.Read, casbin.Write, casbin.Admin:
			c.Next()
		default:
			// 当权限不足时，终止请求并返回JSON响应
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":  "Permission denied",
				"status": http.StatusForbidden,
			})
		}
	}
}
