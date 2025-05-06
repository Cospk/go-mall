package main

import (
	"flag"
	"fmt"
	pb "github.com/Cospk/go-mall/api/rpc/gen/go/user"
	"github.com/Cospk/go-mall/internal/user/application/service"
	"github.com/Cospk/go-mall/internal/user/infrastructure/mysql"
	grpcServer "github.com/Cospk/go-mall/internal/user/interfaces/grpc"

	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// User 用户服务
// 1. 用户注册
// 2. 用户登录
// 3. 用户信息查询
// 4. 用户信息修改
// 5. 用户信息删除
// 6. 用户信息导出
// 7. 用户信息导入

var (
	port       = flag.Int("port", 50051, "服务端口")
	configPath = flag.String("config", "configs/user-service.yaml", "配置文件路径")
)

func main() {
	// 加载配置

	// 连接数据库
	userRepo := mysql.NewUserRepository(db)

	// 创建服务
	userService := service.NewUserService(userRepo)

	// 创建gRPC服务
	grpcSrv := grpc.NewServer()
	userServer := grpcServer.NewUserServiceServer(userService)
	pb.RegisterStreamGreeterServer(grpcSrv, userServer)

	// 启动服务
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("监听端口失败: %v", err)
	}

	// 优雅关闭服务
	go func() {
		if err := grpcSrv.Serve(listener); err != nil {
			log.Fatalf("启动服务失败: %v", err)
		}
	}()
	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")
	grpcSrv.GracefulStop()
	log.Println("服务已经关闭")

}
