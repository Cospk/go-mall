package router

import (
	"github.com/Cospk/go-mall/api/controller"
	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(router *gin.RouterGroup) {
	UserRouter := router.Group("/user/")
	{
		// 注册
		UserRouter.POST("register", controller.RegisterUser)
		// 登录
		UserRouter.POST("login", controller.LoginUser)
		// 登出
		UserRouter.POST("logout", controller.LogoutUser)
		// 刷新token
		UserRouter.POST("refreshToken", controller.RefreshUserToken)

		// 重置密码

		// 获取用户信息

		// 修改用户信息

		// 更新用户信息

	}
}
