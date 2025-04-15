package main

import (
	"github.com/Cospk/go-mall/api/router"
	"github.com/Cospk/go-mall/internal/dal/cache"
	"github.com/Cospk/go-mall/internal/dal/dao"
	"github.com/Cospk/go-mall/pkg/config"
	"github.com/Cospk/go-mall/pkg/logger"
	"go.uber.org/zap"
)

func main() {

	// 初始化配置
	config.InitConfig()

	// 初始化日志
	logger.InitLogger()

	//初始化数据库
	dao.InitGorm()

	// 初始化缓存
	cache.InitRedis()

	// 初始化路由
	Router := router.InitWebRouter()

	err := Router.Run("127.0.0.1:8080")

	if err != nil {
		logger.Logger.Info("服务启动失败", zap.Error(err))
	}

}
