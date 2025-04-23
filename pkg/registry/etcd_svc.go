package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdRegistry 实现基于etcd的服务注册与发现
type EtcdRegistry struct {
	client     *clientv3.Client
	prefix     string
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// EtcdConfig etcd配置
type EtcdConfig struct {
	Endpoints   []string
	DialTimeout time.Duration
	Prefix      string
}

// NewEtcdRegistry 创建etcd注册中心
func NewEtcdRegistry(ctx context.Context, config EtcdConfig) (*EtcdRegistry, error) {
	if len(config.Endpoints) == 0 {
		return nil, ErrInvalidConfig
	}

	if config.DialTimeout == 0 {
		config.DialTimeout = 5 * time.Second
	}

	if config.Prefix == "" {
		config.Prefix = "/services/"
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Endpoints,
		DialTimeout: config.DialTimeout,
	})

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &EtcdRegistry{
		client:     client,
		prefix:     config.Prefix,
		ctx:        ctx,
		cancelFunc: cancel,
	}, nil
}

// Register 注册服务
func (e *EtcdRegistry) Register(ctx context.Context, instance ServiceInstance) error {
	key := fmt.Sprintf("%s%s/%s", e.prefix, instance.Name, instance.ID)
	value, err := json.Marshal(instance)
	if err != nil {
		return err
	}

	// 设置TTL
	lease, err := e.client.Grant(e.ctx, int64(instance.TTL.Seconds()))
	if err != nil {
		return err
	}

	_, err = e.client.Put(e.ctx, key, string(value), clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}

	// 自动续约
	keepAliveCh, err := e.client.KeepAlive(e.ctx, lease.ID)
	if err != nil {
		return err
	}

	// 处理续约响应
	go func() {
		for {
			select {
			case <-e.ctx.Done():
				return
			case <-keepAliveCh:
				// 续约成功，可以记录日志
			}
		}
	}()

	return nil
}

// Deregister 注销服务
func (e *EtcdRegistry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	// 在实际应用中，可能需要先查找服务名称
	// 这里简化处理，假设已知服务名称和实例ID
	_, err := e.client.Delete(e.ctx, instanceID)
	return err
}

// GetService 获取服务实例列表
func (e *EtcdRegistry) GetService(ctx context.Context, serviceName string) ([]ServiceInstance, error) {
	key := fmt.Sprintf("%s%s/", e.prefix, serviceName)
	resp, err := e.client.Get(e.ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, ErrNotFound
	}

	instances := make([]ServiceInstance, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var instance ServiceInstance
		if err := json.Unmarshal(kv.Value, &instance); err != nil {
			continue
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// GetAllServices 获取所有服务
func (e *EtcdRegistry) GetAllServices(ctx context.Context) (map[string][]ServiceInstance, error) {
	resp, err := e.client.Get(e.ctx, e.prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	services := make(map[string][]ServiceInstance)
	for _, kv := range resp.Kvs {
		var instance ServiceInstance
		if err := json.Unmarshal(kv.Value, &instance); err != nil {
			continue
		}

		if _, ok := services[instance.Name]; !ok {
			services[instance.Name] = make([]ServiceInstance, 0)
		}
		services[instance.Name] = append(services[instance.Name], instance)
	}

	return services, nil
}

// Watch 监听服务变化
func (e *EtcdRegistry) Watch(ctx context.Context, serviceName string) (<-chan []ServiceInstance, error) {
	key := fmt.Sprintf("%s%s/", e.prefix, serviceName)

	// 创建通道
	ch := make(chan []ServiceInstance, 10)

	// 启动监听
	go func() {
		defer close(ch)

		watchCh := e.client.Watch(e.ctx, key, clientv3.WithPrefix())
		for {
			select {
			case <-e.ctx.Done():
				return
			case <-watchCh:
				// 服务发生变化，获取最新列表
				instances, err := e.GetService(context.Background(), serviceName)
				if err != nil {
					continue
				}
				ch <- instances
			}
		}
	}()

	return ch, nil
}

// Close 关闭连接
func (e *EtcdRegistry) Close() error {
	e.cancelFunc()
	return e.client.Close()
}
