package router

import (
	"github.com/Cospk/go-mall/api/controller"
	"github.com/Cospk/go-mall/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(router *gin.RouterGroup) {
	UserRouter := router.Group("/user/")
	{
		// 注册
		UserRouter.POST("register", controller.RegisterUser)
		// 登录
		UserRouter.POST("login", controller.LoginUser)
	}
	UserRouter.Use(middleware.AuthMiddleware())
	{
		// 登出
		UserRouter.POST("logout", controller.LogoutUser)
		// 刷新token
		UserRouter.POST("refreshToken", controller.RefreshUserToken)

		// 重置密码
		UserRouter.POST("password/apply-reset", controller.PasswordResetApply)

		// 重置密码
		UserRouter.POST("password/reset", controller.PasswordReset)

		// 获取用户信息
		UserRouter.GET("info", controller.UserInfo)

		// 更新用户信息
		UserRouter.PATCH("info", controller.UpdateUserInfo)

	}
}
