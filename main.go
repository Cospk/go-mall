package main

import (
	"github.com/Cospk/go-mall/global"
	"github.com/Cospk/go-mall/initialize"
	"go.uber.org/zap"
)

func main() {

	// 初始化配置
	initialize.InitConfig()

	// 初始化日志
	initialize.InitLogger()

	// 初始化路由
	Router := initialize.InitWebRouter()

	err := Router.Run("127.0.0.1:8080")

	if err != nil {
		global.Logger.Info("服务启动失败", zap.Error(err))
	}

}
