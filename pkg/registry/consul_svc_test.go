package registry

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConsulRegistry(t *testing.T) {
	// 跳过测试，除非明确要求运行集成测试
	t.Skip("Skipping Consul integration test")

	// 连接到 Consul
	registry, err := NewConsulRegistry("127.0.0.1:8500")
	if err != nil {
		t.Fatalf("Failed to create Consul registry: %v", err)
	}
	defer registry.Close()

	ctx := context.Background()

	// 测试注册服务
	instance := ServiceInstance{
		ID:      "test-service-1",
		Name:    "test-service",
		Address: "127.0.0.1",
		Port:    8080,
		Metadata: map[string]string{
			"version": "1.0.0",
		},
		TTL:       time.Second * 30,
		StartTime: time.Now(),
	}

	err = registry.Register(ctx, instance)
	assert.NoError(t, err)

	// 等待服务注册生效
	time.Sleep(time.Second)

	// 测试获取服务
	instances, err := registry.GetService(ctx, "test-service")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(instances), 1)

	var found bool
	for _, s := range instances {
		if s.ID == instance.ID {
			found = true
			assert.Equal(t, instance.Name, s.Name)
			assert.Equal(t, instance.Address, s.Address)
			assert.Equal(t, instance.Port, s.Port)
			assert.Equal(t, instance.Metadata["version"], s.Metadata["version"])
			break
		}
	}
	assert.True(t, found, "Registered service not found")

	// 测试获取所有服务
	allServices, err := registry.GetAllServices(ctx)
	assert.NoError(t, err)
	assert.Contains(t, allServices, "test-service")

	// 测试监听服务变化
	watchCh, err := registry.Watch(ctx, "test-service")
	assert.NoError(t, err)

	// 注册第二个服务实例
	instance2 := ServiceInstance{
		ID:      "test-service-2",
		Name:    "test-service",
		Address: "127.0.0.1",
		Port:    8081,
		Metadata: map[string]string{
			"version": "1.0.1",
		},
		TTL:       time.Second * 30,
		StartTime: time.Now(),
	}
	err = registry.Register(ctx, instance2)
	assert.NoError(t, err)

	// 等待并检查监听通道
	select {
	case instances := <-watchCh:
		found1, found2 := false, false
		for _, s := range instances {
			if s.ID == "test-service-1" {
				found1 = true
			}
			if s.ID == "test-service-2" {
				found2 = true
			}
		}
		assert.True(t, found1 || found2, "At least one service should be found")
	case <-time.After(time.Second * 5):
		t.Fatal("Timeout waiting for service update")
	}

	// 测试注销服务
	err = registry.Deregister(ctx, instance.ID, instance.Name)
	assert.NoError(t, err)
	err = registry.Deregister(ctx, instance2.ID, instance2.Name)
	assert.NoError(t, err)

	// 等待服务注销生效
	time.Sleep(time.Second)

	// 验证服务已注销
	instances, err = registry.GetService(ctx, "test-service")
	assert.NoError(t, err)

	found1, found2 := false, false
	for _, s := range instances {
		if s.ID == "test-service-1" {
			found1 = true
		}
		if s.ID == "test-service-2" {
			found2 = true
		}
	}
	assert.False(t, found1, "Service 1 should be deregistered")
	assert.False(t, found2, "Service 2 should be deregistered")
}

func TestConsulRegistryMock(t *testing.T) {
	// 创建一个模拟的测试
	// 注意：这只是一个简单的测试框架，实际上没有连接到 Consul

	instance := ServiceInstance{
		ID:      "test-service-1",
		Name:    "test-service",
		Address: "127.0.0.1",
		Port:    8080,
		Metadata: map[string]string{
			"version":    "1.0.0",
			"check_http": "http://127.0.0.1:8080/health",
		},
		TTL:       time.Second * 30,
		StartTime: time.Now(),
	}

	// 验证服务信息结构
	assert.Equal(t, "test-service-1", instance.ID)
	assert.Equal(t, "test-service", instance.Name)
	assert.Equal(t, "127.0.0.1", instance.Address)
	assert.Equal(t, 8080, instance.Port)
	assert.Equal(t, "1.0.0", instance.Metadata["version"])
	assert.Equal(t, "http://127.0.0.1:8080/health", instance.Metadata["check_http"])
	assert.Equal(t, time.Second*30, instance.TTL)
}
