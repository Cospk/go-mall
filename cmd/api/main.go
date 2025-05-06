package main

import (
	"context"
	"github.com/Cospk/go-mall/internal/api/infrastructure/config"
	"github.com/Cospk/go-mall/internal/api/interfaces/http/router"
	"github.com/Cospk/go-mall/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// api网关，用户端请求，api网关转发请求到对应的服务
// 服务注册，服务发现，服务调用

func main() {

	// 初始化配置
	config.InitApiConfig()

	// 初始化日志
	logger.InitLogger()

	// 初始化路由
	Router := router.InitApiRouter()

	//err := Router.Run(config.App.Host + ":" + strconv.Itoa(config.App.Port))

	//if err != nil {
	//	logger.Logger.Info("服务启动失败", zap.Error(err))
	//}

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    config.App.Host + ":" + strconv.Itoa(config.App.Port),
		Handler: Router,
	}

	// 在单独的goroutine中启动服务器
	go func() {
		logger.Logger.Info("服务启动成功", zap.String("地址", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal("服务启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	// kill (无参数) 默认发送 syscall.SIGTERM
	// kill -2 发送 syscall.SIGINT
	// kill -9 发送 syscall.SIGKILL，但无法被捕获，所以不需要添加
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Logger.Info("正在关闭服务...")

	// 设置关闭超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		logger.Logger.Fatal("服务关闭异常", zap.Error(err))
	}

	// 等待关闭完成（超时或正常关闭）
	logger.Logger.Info("服务已关闭")
}
