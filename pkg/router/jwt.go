package router

import (
	"github.com/denovo/permission/pkg/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// JWT token验证中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":  ErrorAuthCheckTokenFail,
				"status": http.StatusBadRequest,
			})
		}
		claims, err := util.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":  ErrorAuthCheckTokenFail,
				"status": http.StatusBadRequest,
			})
		} else if time.Now().Unix() > claims.ExpiresAt {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":  ErrorAuthCheckTokenExpired,
				"status": http.StatusBadRequest,
			})
		}
		c.Next()
	}
}
