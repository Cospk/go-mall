package registry

import (
	"context"
	"github.com/Cospk/go-mall/pkg/errcode"
	"time"
)

// ServiceInstance 表示一个服务实例
type ServiceInstance struct {
	ID        string            // 实例唯一ID
	Name      string            // 服务名称
	Address   string            // 服务地址
	Port      int               // 服务端口
	Metadata  map[string]string // 元数据
	TTL       time.Duration     // 生存时间
	StartTime time.Time         // 启动时间
}

// Registry 定义服务注册与发现接口
type Registry interface {
	// Register 注册服务实例
	Register(ctx context.Context, instance ServiceInstance) error

	// Deregister 注销服务实例
	Deregister(ctx context.Context, instanceID string, serviceName string) error

	// GetService 获取指定服务的所有实例
	GetService(ctx context.Context, serviceName string) ([]ServiceInstance, error)

	// GetAllServices 获取所有服务
	GetAllServices(ctx context.Context) (map[string][]ServiceInstance, error)

	// Watch 监听服务变化
	Watch(ctx context.Context, serviceName string) (<-chan []ServiceInstance, error)

	// Close 关闭注册中心连接
	Close() error
}

// 常见错误定义
var (
	ErrNotFound         = errcode.NewError(5000, "service not found")
	ErrRegisterFailed   = errcode.NewError(5001, "failed to register service")
	ErrDeregisterFailed = errcode.NewError(5002, "failed to deregister service")
	ErrInvalidConfig    = errcode.NewError(5003, "invalid registry configuration")
)
