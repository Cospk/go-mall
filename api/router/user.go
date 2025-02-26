package router

import (
	"github.com/Cospk/go-mall/api/controller"
	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(router *gin.RouterGroup) {
	UserRouter := router.Group("/user")
	{
		// 登录
		UserRouter.POST("/login", controller.LoginUser)
	}
}
