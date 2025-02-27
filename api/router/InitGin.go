package router

import (
	"errors"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/logger"
	"github.com/Cospk/go-mall/pkg/middleware"
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

	Router.GET("/customized-error-test", func(c *gin.Context) {

		// 使用 Wrap 包装原因error 生成 项目error
		err := errors.New("a dao error")
		appErr := errcode.Wrap("包装错误", err)
		bAppErr := errcode.Wrap("再包装错误", appErr)
		logger.NewLogger(c).Error("记录错误", "err", bAppErr)

		// 预定义的ErrServer, 给其追加错误原因的error
		err = errors.New("a domain error")
		apiErr := errcode.ErrServer.WithCause(err)
		logger.NewLogger(c).Error("API执行中出现错误", "err", apiErr)

		c.JSON(apiErr.HttpStatusCode(), gin.H{
			"code": apiErr.Code(),
			"msg":  apiErr.Msg(),
		})

	})

	router := Router.Group("api/v1")
	// 注册路由
	RegisterUserRouter(router)
	RegisterDemoRouter(router)

	return Router
}
