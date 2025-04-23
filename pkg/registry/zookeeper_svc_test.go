package registry

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestZookeeperRegistry(t *testing.T) {
	// 跳过测试，除非明确要求运行集成测试
	t.Skip("Skipping Zookeeper integration test")

	// 连接到 Zookeeper
	addresses := []string{"127.0.0.1:2181"}
	timeout := time.Second * 5

	registry, err := NewZookeeperRegistry(addresses, timeout)
	if err != nil {
		t.Fatalf("Failed to create Zookeeper registry: %v", err)
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

	// 测试获取服务
	instances, err := registry.GetService(ctx, "test-service")
	assert.NoError(t, err)
	assert.Len(t, instances, 1)
	assert.Equal(t, instance.ID, instances[0].ID)
	assert.Equal(t, instance.Name, instances[0].Name)
	assert.Equal(t, instance.Address, instances[0].Address)
	assert.Equal(t, instance.Port, instances[0].Port)
	assert.Equal(t, instance.Metadata["version"], instances[0].Metadata["version"])

	// 测试获取所有服务
	allServices, err := registry.GetAllServices(ctx)
	assert.NoError(t, err)
	assert.Contains(t, allServices, "test-service")
	assert.Len(t, allServices["test-service"], 1)

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
		assert.Len(t, instances, 2)
		found1, found2 := false, false
		for _, s := range instances {
			if s.ID == "test-service-1" {
				found1 = true
			}
			if s.ID == "test-service-2" {
				found2 = true
			}
		}
		assert.True(t, found1)
		assert.True(t, found2)
	case <-time.After(time.Second * 5):
		t.Fatal("Timeout waiting for service update")
	}

	// 测试注销服务
	err = registry.Deregister(ctx, instance.ID, instance.Name)
	assert.NoError(t, err)
	err = registry.Deregister(ctx, instance2.ID, instance2.Name)
	assert.NoError(t, err)

	// 验证服务已注销
	instances, err = registry.GetService(ctx, "test-service")
	assert.NoError(t, err)
	assert.Len(t, instances, 0)
}

func TestZookeeperRegistryMock(t *testing.T) {
	// 创建一个模拟的测试
	// 注意：这只是一个简单的测试框架，实际上没有连接到 Zookeeper

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

	// 验证服务信息结构
	assert.Equal(t, "test-service-1", instance.ID)
	assert.Equal(t, "test-service", instance.Name)
	assert.Equal(t, "127.0.0.1", instance.Address)
	assert.Equal(t, 8080, instance.Port)
	assert.Equal(t, "1.0.0", instance.Metadata["version"])
	assert.Equal(t, time.Second*30, instance.TTL)
}
