package router

import (
	"github.com/Cospk/go-mall/api/demo/controller"
	"github.com/gin-gonic/gin"
)

func RegisterDemoRouter(router *gin.RouterGroup) {
	DemoRouter := router.Group("/demo")
	{
		// 测试ping
		DemoRouter.GET("/ping", controller.TestPing)
		// 测试配置文件读取
		DemoRouter.GET("/config-read", controller.TestConfigRead)
		// 测试日志的使用
		DemoRouter.GET("/log", controller.TestLog)
		// 测试服务的访问日志
		DemoRouter.GET("/access-log-test", controller.TestAccessLog)
		// 测试服务的崩溃日志
		DemoRouter.GET("/panic-log-test", controller.TestPanicLog)
		// 测试项目自定义的AppError 打印错误链和错误发生的位置
		DemoRouter.GET("/customized-error-test", controller.TestAppError)
		// 测试统一响应 -- 返回对象数据
		DemoRouter.GET("/response-obj", controller.TestResponseObj)
		// 测试统一响应 -- 返回数组数据
		DemoRouter.GET("/response-arr", controller.TestResponseList)
		// 测试统一响应 -- 返回错误
		DemoRouter.GET("/response-error", controller.TestResponseError)
		// 测试GORM Logger
		DemoRouter.GET("/gorm-logger", controller.TestGormLogger)

		// 测试Token生成
		DemoRouter.GET("/get-token", controller.TestGetToken)
		// 测试Token验证
		DemoRouter.GET("/verify-token", controller.TestVerifyToken)
		// 测试Token刷新
		DemoRouter.GET("/refresh-token", controller.TestRefreshToken)
	}
}
