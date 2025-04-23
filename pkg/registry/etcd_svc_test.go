package registry

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试 NewEtcdRegistry 函数
func TestNewEtcdRegistry(t *testing.T) {
	// 测试无效配置
	_, err := NewEtcdRegistry(context.Background(), EtcdConfig{})
	assert.Equal(t, ErrInvalidConfig, err)

	// 测试有效配置但无法连接的情况
	config := EtcdConfig{
		Endpoints:   []string{"localhost:23790"}, // 使用一个不可能存在的端口
		DialTimeout: 1 * time.Second,
		Prefix:      "/test/",
	}
	_, err = NewEtcdRegistry(context.Background(), config)
	assert.Error(t, err)
}

// 获取测试用的 etcd 配置
func getTestEtcdConfig() EtcdConfig {
	// 如果环境变量中设置了 etcd 地址，则使用环境变量中的地址
	endpoints := []string{"localhost:2379"}
	if env := os.Getenv("ETCD_ENDPOINTS"); env != "" {
		endpoints = []string{env}
	}

	return EtcdConfig{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
		Prefix:      "/test-services/",
	}
}

// 跳过测试如果无法连接到 etcd
func skipIfEtcdNotAvailable(t *testing.T) *EtcdRegistry {
	// 如果设置了环境变量，则跳过测试
	if os.Getenv("SKIP_ETCD_TESTS") == "true" {
		t.Skip("跳过 etcd 测试，设置 SKIP_ETCD_TESTS=false 以启用")
	}

	config := getTestEtcdConfig()
	registry, err := NewEtcdRegistry(context.Background(), config)
	if err != nil {
		t.Skipf("无法连接到 etcd: %v", err)
	}

	return registry
}

// 测试 Register 和 GetService 方法
func TestEtcdRegistry_RegisterAndGetService(t *testing.T) {
	registry := skipIfEtcdNotAvailable(t)
	defer registry.Close()

	// 创建测试实例
	instance := ServiceInstance{
		ID:      "test-instance-1",
		Name:    "test-service",
		Address: "localhost",
		Port:    8080,
		TTL:     time.Second * 10,
	}

	// 注册服务
	err := registry.Register(context.Background(), instance)
	require.NoError(t, err)

	// 获取服务
	instances, err := registry.GetService(context.Background(), instance.Name)
	require.NoError(t, err)
	require.Len(t, instances, 1)
	assert.Equal(t, instance.ID, instances[0].ID)
	assert.Equal(t, instance.Name, instances[0].Name)
	assert.Equal(t, instance.Address, instances[0].Address)
	assert.Equal(t, instance.Port, instances[0].Port)
}

// 测试 Deregister 方法
func TestEtcdRegistry_Deregister(t *testing.T) {
	registry := skipIfEtcdNotAvailable(t)
	defer registry.Close()

	// 创建测试实例
	instance := ServiceInstance{
		ID:      "test-instance-2",
		Name:    "test-service-deregister",
		Address: "localhost",
		Port:    8081,
		TTL:     time.Second * 10,
	}

	// 注册服务
	err := registry.Register(context.Background(), instance)
	require.NoError(t, err)

	// 确认服务已注册
	instances, err := registry.GetService(context.Background(), instance.Name)
	require.NoError(t, err)
	require.Len(t, instances, 1)

	// 注销服务
	err = registry.Deregister(context.Background(), instance.ID, instance.Name)
	require.NoError(t, err)

	// 确认服务已注销
	_, err = registry.GetService(context.Background(), instance.Name)
	assert.Equal(t, ErrNotFound, err)
}

// 测试 GetAllServices 方法
func TestEtcdRegistry_GetAllServices(t *testing.T) {
	registry := skipIfEtcdNotAvailable(t)
	defer registry.Close()

	// 创建测试实例
	instance1 := ServiceInstance{
		ID:      "test-instance-all-1",
		Name:    "test-service-all-1",
		Address: "localhost",
		Port:    8082,
		TTL:     time.Second * 10,
	}

	instance2 := ServiceInstance{
		ID:      "test-instance-all-2",
		Name:    "test-service-all-2",
		Address: "localhost",
		Port:    8083,
		TTL:     time.Second * 10,
	}

	// 注册服务
	err := registry.Register(context.Background(), instance1)
	require.NoError(t, err)

	err = registry.Register(context.Background(), instance2)
	require.NoError(t, err)

	// 获取所有服务
	services, err := registry.GetAllServices(context.Background())
	require.NoError(t, err)

	// 验证结果
	assert.Contains(t, services, instance1.Name)
	assert.Contains(t, services, instance2.Name)
	assert.Len(t, services[instance1.Name], 1)
	assert.Len(t, services[instance2.Name], 1)

	// 清理
	registry.Deregister(context.Background(), instance1.ID, instance1.Name)
	registry.Deregister(context.Background(), instance2.ID, instance2.Name)
}

// 测试 Watch 方法
func TestEtcdRegistry_Watch(t *testing.T) {
	registry := skipIfEtcdNotAvailable(t)
	defer registry.Close()

	serviceName := "test-service-watch"

	// 启动监听
	ch, err := registry.Watch(context.Background(), serviceName)
	require.NoError(t, err)

	// 创建测试实例
	instance := ServiceInstance{
		ID:      "test-instance-watch",
		Name:    serviceName,
		Address: "localhost",
		Port:    8084,
		TTL:     time.Second * 10,
	}

	// 在另一个 goroutine 中注册服务
	go func() {
		time.Sleep(100 * time.Millisecond) // 稍微延迟，确保监听已经启动
		err := registry.Register(context.Background(), instance)
		if err != nil {
			t.Errorf("注册服务失败: %v", err)
		}
	}()

	// 等待监听结果
	select {
	case instances := <-ch:
		assert.Len(t, instances, 1)
		assert.Equal(t, instance.ID, instances[0].ID)
	case <-time.After(3 * time.Second):
		t.Fatal("监听超时")
	}

	// 清理
	registry.Deregister(context.Background(), instance.ID, instance.Name)
}

// 测试 TTL 和自动续约
func TestEtcdRegistry_TTLAndKeepAlive(t *testing.T) {
	registry := skipIfEtcdNotAvailable(t)
	defer registry.Close()

	// 创建测试实例，使用较短的 TTL
	instance := ServiceInstance{
		ID:      "test-instance-ttl",
		Name:    "test-service-ttl",
		Address: "localhost",
		Port:    8085,
		TTL:     time.Second * 2, // 2秒 TTL
	}

	// 注册服务
	err := registry.Register(context.Background(), instance)
	require.NoError(t, err)

	// 等待超过 TTL 的时间，但由于自动续约，服务应该仍然存在
	time.Sleep(time.Second * 3)

	// 获取服务，确认仍然存在
	instances, err := registry.GetService(context.Background(), instance.Name)
	require.NoError(t, err)
	require.Len(t, instances, 1)
	assert.Equal(t, instance.ID, instances[0].ID)

	// 清理
	registry.Deregister(context.Background(), instance.ID, instance.Name)
}

// 测试关闭连接后的行为
func TestEtcdRegistry_CloseAndReconnect(t *testing.T) {
	registry := skipIfEtcdNotAvailable(t)

	// 关闭连接
	err := registry.Close()
	require.NoError(t, err)

	// 尝试在关闭后使用，应该返回错误
	instance := ServiceInstance{
		ID:      "test-instance-close",
		Name:    "test-service-close",
		Address: "localhost",
		Port:    8086,
		TTL:     time.Second * 10,
	}

	err = registry.Register(context.Background(), instance)
	assert.Error(t, err, "在关闭连接后应该返回错误")

	// 重新创建连接
	newRegistry := skipIfEtcdNotAvailable(t)
	defer newRegistry.Close()

	// 使用新连接应该正常
	err = newRegistry.Register(context.Background(), instance)
	assert.NoError(t, err)

	// 清理
	newRegistry.Deregister(context.Background(), instance.ID, instance.Name)
}
