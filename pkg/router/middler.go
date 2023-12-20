package router

import (
	"fmt"
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
		enforce := cb.Enforcer.Enforce("A", "kube-system", "B", casbin.Write)
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
		//TODO:
		case casbin.Write, casbin.Read:
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

//list资源列表单独设置中间件

// ExtractParams 对于read或者write的中间件权限控制
func ExtractParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求中获取参数，这里以获取 query 参数为例
		user := c.Query("user")
		ns := c.Query("namespace")
		res := c.Query("resource")
		action := c.Query("action")

		// 将获取的参数存入 Gin 的上下文中
		c.Set("user", user)
		c.Set("ns", ns)
		c.Set("res", res)
		c.Set("action", action)

		// 继续后续处理
		c.Next()
	}
}
