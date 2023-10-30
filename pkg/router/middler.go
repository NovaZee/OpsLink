package router

import (
	"github.com/denovo/permission/pkg/casbin"
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
		logger.Infow(" http request",
			zap.Any(" request", c.Request.URL),
			zap.Any("response", c.Writer.Status()),
			zap.String("ip", c.ClientIP()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(nowTime)),
		)
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
