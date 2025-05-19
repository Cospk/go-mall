package rpc

import (
	"context"
	"fmt"
	"github.com/Cospk/go-mall/api/rpc/gen/go/user"
	"github.com/Cospk/go-mall/pkg/registry"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// UserServer 用户服务gRPC服务器
type UserServer struct {
	server     *grpc.Server
	userServer user.StreamGreeterServer
	registry   registry.Registry
	instance   registry.ServiceInstance
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewUserServer 创建新的用户服务gRPC服务器
func NewUserServer(userServer user.StreamGreeterServer) *UserServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &UserServer{
		server:     grpc.NewServer(),
		userServer: userServer,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start 启动gRPC服务器并注册到服务发现
func (s *UserServer) Start(address string) error {
	// 解析主机和端口
	host, port, err := parseAddress(address)
	if err != nil {
		log.Fatalf("无效的地址格式: %v", err)
		return err
	}

	// 初始化服务注册
	if err := s.initRegistry(); err != nil {
		log.Printf("警告: 服务注册初始化失败: %v", err)
		// 继续启动服务，但不进行注册
	}

	// 启动gRPC服务器
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("监听失败: %v", err)
		return err
	}

	// 注册用户服务
	user.RegisterStreamGreeterServer(s.server, s.userServer)

	// 如果注册中心已初始化，则注册服务
	if s.registry != nil {
		// 创建服务实例
		s.instance = registry.ServiceInstance{
			ID:        fmt.Sprintf("user-service-%s", uuid.New().String()),
			Name:      "user-service",
			Address:   host,
			Port:      port,
			Metadata:  map[string]string{"version": "1.0.0"},
			TTL:       time.Second * 30,
			StartTime: time.Now(),
		}

		// 注册服务
		if err := s.registry.Register(s.ctx, s.instance); err != nil {
			log.Printf("服务注册失败: %v", err)
		} else {
			log.Printf("服务已注册到注册中心: %s:%d", host, port)
		}
	}

	// 处理优雅关闭
	go s.handleSignals()

	log.Printf("用户服务启动在: %s", address)
	return s.server.Serve(lis)
}

// Stop 停止gRPC服务器并从服务发现中注销
func (s *UserServer) Stop() {
	// 从注册中心注销服务
	if s.registry != nil && s.instance.ID != "" {
		if err := s.registry.Deregister(s.ctx, s.instance.ID, s.instance.Name); err != nil {
			log.Printf("服务注销失败: %v", err)
		} else {
			log.Printf("服务已从注册中心注销")
		}
		s.registry.Close()
	}

	// 取消上下文
	s.cancel()

	// 优雅停止gRPC服务器
	if s.server != nil {
		s.server.GracefulStop()
		log.Println("用户服务已停止")
	}
}

// 初始化服务注册
func (s *UserServer) initRegistry() error {
	// 从配置中获取注册中心类型
	registryType := viper.GetString("registry.type")
	if registryType == "" {
		registryType = "etcd" // 默认使用etcd
	}

	// 获取注册中心地址
	endpoints := viper.GetStringSlice("registry.endpoints")
	if len(endpoints) == 0 {
		endpoints = []string{"localhost:2379"} // 默认etcd地址
	}

	// 获取超时设置
	timeout := viper.GetInt("registry.timeout")
	if timeout == 0 {
		timeout = 5 // 默认5秒
	}

	// 获取前缀
	prefix := viper.GetString("registry.prefix")
	if prefix == "" {
		prefix = "/services/" // 默认前缀
	}

	// 创建注册中心配置
	config := registry.RegistryConfig{
		Type:      registry.RegistryType(registryType),
		Endpoints: endpoints,
		Timeout:   timeout,
		Prefix:    prefix,
	}

	// 创建注册中心客户端
	reg, err := registry.NewRegistry(s.ctx, config)
	if err != nil {
		return err
	}

	s.registry = reg
	return nil
}

// 处理系统信号，实现优雅关闭
func (s *UserServer) handleSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-signals
	log.Println("接收到关闭信号，开始优雅关闭...")
	s.Stop()
}

// 解析地址为主机和端口
func parseAddress(address string) (string, int, error) {
	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return "", 0, err
	}

	// 如果主机为空，使用本机IP
	if host == "" {
		host = "localhost"
	}

	// 解析端口
	var port int
	_, err = fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		return "", 0, err
	}

	return host, port, nil
}
