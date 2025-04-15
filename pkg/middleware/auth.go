package middleware

import (
	"github.com/Cospk/go-mall/internal/logic/service"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的Token
		token := c.Request.Header.Get("Authorization")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 403, "msg": "未授权访问"})
			return
		}

		tokenVerify, err := service.NewUserService(c).VerifyAccessToken(token)
		if err != nil { // 验证Token时服务出错
			resp.NewResponse(c).Error(errcode.ErrServer)
			c.Abort()
			return
		}
		if !tokenVerify.Approved { // Token未通过验证
			resp.NewResponse(c).Error(errcode.ErrToken)
			c.Abort()
			return
		}
		c.Set("userId", tokenVerify.UserId)
		c.Set("sessionId", tokenVerify.SessionId)
		c.Set("platform", tokenVerify.Platform)
		c.Next()
	}
}
