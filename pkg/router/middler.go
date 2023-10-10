package router

import (
	"github.com/gin-gonic/gin"
	"github.com/oppslink/protocol/logger"
	"go.uber.org/zap"
	"time"
)

// Logger 打印自己的日志库
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
