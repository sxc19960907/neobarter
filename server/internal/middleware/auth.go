package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	jwtPkg "github.com/neobarter/server/internal/pkg/jwt"
	"github.com/neobarter/server/internal/pkg/response"
)

func Auth(jwtManager *jwtPkg.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证信息")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		claims, err := jwtManager.ParseToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "认证已过期，请重新登录")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("phone", claims.Phone)
		c.Set("user_type", claims.UserType)
		c.Next()
	}
}

// GetUserID 从上下文获取当前用户ID
func GetUserID(c *gin.Context) int64 {
	userID, _ := c.Get("user_id")
	return userID.(int64)
}

// GetUserType 从上下文获取用户类型
func GetUserType(c *gin.Context) string {
	userType, _ := c.Get("user_type")
	return userType.(string)
}
