package initialize

import (
	"github.com/Cospk/go-mall/middleware"
	"github.com/gin-gonic/gin"
)

func InitWebRouter() *gin.Engine {
	Router := gin.Default()
	Router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 使用中间件
	Router.Use(middleware.TraceMiddleware(), middleware.LoggerMiddleware(), middleware.RecoveryMiddleware())

	Router.Group("api/v1")
	// 注册路由

	return Router
}
