package registry

import (
	"context"
	"testing"
	"time"
)

// 测试工厂方法
func TestNewRegistry(t *testing.T) {
	tests := []struct {
		name    string
		config  RegistryConfig
		wantErr bool
	}{
		{
			name: "etcd registry",
			config: RegistryConfig{
				Type:      RegistryTypeEtcd,
				Endpoints: []string{"localhost:2379"},
				Timeout:   5,
				Prefix:    "/services/",
			},
			wantErr: false,
		},
		{
			name: "consul registry",
			config: RegistryConfig{
				Type:      RegistryTypeConsul,
				Endpoints: []string{"localhost:8500"},
				Timeout:   5,
				Prefix:    "/services/",
			},
			wantErr: true, // 因为我们还没有实现consul
		},
		{
			name: "zookeeper registry",
			config: RegistryConfig{
				Type:      RegistryTypeZookeeper,
				Endpoints: []string{"localhost:2181"},
				Timeout:   5,
				Prefix:    "/services/",
			},
			wantErr: true, // 因为我们还没有实现zookeeper
		},
		{
			name: "unsupported registry",
			config: RegistryConfig{
				Type:      "unsupported",
				Endpoints: []string{"localhost:2379"},
				Timeout:   5,
				Prefix:    "/services/",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRegistry(context.Background(), tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRegistry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// 测试服务实例结构体
func TestServiceInstance(t *testing.T) {
	instance := ServiceInstance{
		ID:        "test-service-1",
		Name:      "test-service",
		Address:   "localhost",
		Port:      8080,
		Metadata:  map[string]string{"version": "1.0.0"},
		TTL:       time.Second * 30,
		StartTime: time.Now(),
	}

	if instance.ID != "test-service-1" {
		t.Errorf("Expected ID to be 'test-service-1', got %s", instance.ID)
	}

	if instance.Name != "test-service" {
		t.Errorf("Expected Name to be 'test-service', got %s", instance.Name)
	}

	if instance.Address != "localhost" {
		t.Errorf("Expected Address to be 'localhost', got %s", instance.Address)
	}

	if instance.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", instance.Port)
	}

	if version, ok := instance.Metadata["version"]; !ok || version != "1.0.0" {
		t.Errorf("Expected Metadata['version'] to be '1.0.0', got %s", version)
	}

	if instance.TTL != time.Second*30 {
		t.Errorf("Expected TTL to be 30s, got %v", instance.TTL)
	}
}
