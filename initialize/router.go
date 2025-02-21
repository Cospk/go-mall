package initialize

import "github.com/gin-gonic/gin"

func InitWebRouter() *gin.Engine {
	Router := gin.Default()
	Router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 使用中间件
	Router.Use()

	Router.Group("api/v1")
	// 注册路由

	return Router
}
