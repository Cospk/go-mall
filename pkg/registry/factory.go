package registry

import (
	"context"
	"fmt"
	"time"
)

// RegistryType 注册中心类型
type RegistryType string

const (
	RegistryTypeEtcd      RegistryType = "etcd"
	RegistryTypeConsul    RegistryType = "consul"
	RegistryTypeZookeeper RegistryType = "zookeeper"
)

// RegistryConfig 注册中心配置
type RegistryConfig struct {
	Type      RegistryType
	Endpoints []string
	Timeout   int
	Prefix    string
	// 其他通用配置
}

// NewRegistry 创建注册中心实例
func NewRegistry(ctx context.Context, config RegistryConfig) (Registry, error) {
	switch config.Type {
	case RegistryTypeEtcd:
		return NewEtcdRegistry(ctx, EtcdConfig{
			Endpoints:   config.Endpoints,
			DialTimeout: time.Duration(config.Timeout) * time.Second,
			Prefix:      config.Prefix,
		})
	case RegistryTypeConsul:
		// 返回Consul实现
		return nil, fmt.Errorf("consul registry not implemented yet")
	case RegistryTypeZookeeper:
		// 返回Zookeeper实现
		return nil, fmt.Errorf("zookeeper registry not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported registry type: %s", config.Type)
	}
}
