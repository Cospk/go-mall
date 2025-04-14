package middleware

import (
	"github.com/Cospk/go-mall/internal/dal/cache"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
	"net/http"
	"strings"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的Token
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			// 尝试从Cookie中获取
			accessToken, _ = c.Cookie("token")
		}

		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": errcode.CodeUnauthorized, "msg": "未授权访问"})
			return
		}

		// 去掉可能的Bearer前缀
		accessToken = strings.TrimPrefix(accessToken, "Bearer ")

		// 首先尝试通过JWT直接验证Token
		userId, jwtErr := util.ParseUserIdFromToken(accessToken)

		// 如果JWT验证成功，再检查Redis中是否存在该Token（用于支持主动失效）
		if jwtErr == nil && userId > 0 {
			session, err := cache.GetAccessToken(c, accessToken)
			if err != nil {
				// Redis错误但JWT验证通过，仍然允许访问
				logger.New(c).Warn("Redis error but JWT valid", "err", err)
				// 设置用户ID到上下文
				c.Set("userId", userId)
				c.Next()
				return
			}

			if session.UserId > 0 {
				// Token有效，设置用户ID到上下文
				c.Set("userId", session.UserId)
				c.Set("platform", session.Platform)
				c.Set("sessionId", session.SessionId)
				c.Next()
				return
			}
		}

		// Token无效
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": errcode.CodeUnauthorized, "msg": "未授权访问"})
	}
}
