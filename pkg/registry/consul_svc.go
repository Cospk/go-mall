package registry

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
)

// ConsulRegistry 是基于 Consul 的服务注册与发现实现
type ConsulRegistry struct {
	client   *api.Client
	mutex    sync.RWMutex
	watchMap map[string]chan []ServiceInstance
}

// NewConsulRegistry 创建一个新的 Consul 注册中心
func NewConsulRegistry(addr string) (*ConsulRegistry, error) {
	config := api.DefaultConfig()
	if addr != "" {
		config.Address = addr
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, ErrInvalidConfig
	}

	return &ConsulRegistry{
		client:   client,
		watchMap: make(map[string]chan []ServiceInstance),
	}, nil
}

// Register 注册服务实例
func (r *ConsulRegistry) Register(ctx context.Context, instance ServiceInstance) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 创建 Consul 服务注册信息
	registration := &api.AgentServiceRegistration{
		ID:      instance.ID,
		Name:    instance.Name,
		Address: instance.Address,
		Port:    instance.Port,
		Tags:    []string{},
		Meta:    instance.Metadata,
	}

	// 添加健康检查
	if instance.TTL > 0 {
		ttl := fmt.Sprintf("%ds", int(instance.TTL.Seconds()))
		registration.Check = &api.AgentServiceCheck{
			TTL:                            ttl,
			DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", int(instance.TTL.Seconds()*2)),
		}
	} else if checkURL, ok := instance.Metadata["check_http"]; ok {
		registration.Check = &api.AgentServiceCheck{
			HTTP:                           checkURL,
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30s",
		}
	}

	// 注册服务
	if err := r.client.Agent().ServiceRegister(registration); err != nil {
		return ErrRegisterFailed
	}

	// 如果使用 TTL 检查，需要定期更新
	if instance.TTL > 0 {
		go r.keepAlive(instance.ID, instance.TTL)
	}

	return nil
}

// keepAlive 定期更新服务健康状态
func (r *ConsulRegistry) keepAlive(serviceID string, ttl time.Duration) {
	ticker := time.NewTicker(ttl / 2)
	defer ticker.Stop()

	for range ticker.C {
		err := r.client.Agent().UpdateTTL("service:"+serviceID, "healthy", api.HealthPassing)
		if err != nil {
			// 如果更新失败，可能服务已被注销，停止更新
			return
		}
	}
}

// Deregister 注销服务实例
func (r *ConsulRegistry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := r.client.Agent().ServiceDeregister(instanceID); err != nil {
		return ErrDeregisterFailed
	}

	return nil
}

// GetService 获取指定服务的所有实例
func (r *ConsulRegistry) GetService(ctx context.Context, serviceName string) ([]ServiceInstance, error) {
	// 查询服务
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return []ServiceInstance{}, nil
	}

	instances := make([]ServiceInstance, 0, len(entries))
	for _, entry := range entries {
		instance := ServiceInstance{
			ID:        entry.Service.ID,
			Name:      entry.Service.Service,
			Address:   entry.Service.Address,
			Port:      entry.Service.Port,
			Metadata:  entry.Service.Meta,
			StartTime: time.Now(), // Consul 不提供启动时间，使用当前时间
		}

		// 尝试从元数据中获取 TTL
		if ttlStr, ok := entry.Service.Meta["ttl"]; ok {
			if ttl, err := strconv.Atoi(ttlStr); err == nil {
				instance.TTL = time.Duration(ttl) * time.Second
			}
		}

		instances = append(instances, instance)
	}

	return instances, nil
}

// GetAllServices 获取所有服务
func (r *ConsulRegistry) GetAllServices(ctx context.Context) (map[string][]ServiceInstance, error) {
	// 获取所有服务名称
	services, _, err := r.client.Catalog().Services(nil)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]ServiceInstance)
	for serviceName := range services {
		instances, err := r.GetService(ctx, serviceName)
		if err != nil {
			continue
		}
		result[serviceName] = instances
	}

	return result, nil
}

// Watch 监听服务变化
func (r *ConsulRegistry) Watch(ctx context.Context, serviceName string) (<-chan []ServiceInstance, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 检查是否已经在监听该服务
	if ch, ok := r.watchMap[serviceName]; ok {
		return ch, nil
	}

	// 创建通道
	watchCh := make(chan []ServiceInstance, 10)
	r.watchMap[serviceName] = watchCh

	// 启动监听协程
	go r.watchService(ctx, serviceName, watchCh)

	return watchCh, nil
}

// watchService 监听服务变化并发送更新
func (r *ConsulRegistry) watchService(ctx context.Context, serviceName string, watchCh chan []ServiceInstance) {
	defer func() {
		r.mutex.Lock()
		delete(r.watchMap, serviceName)
		r.mutex.Unlock()
		close(watchCh)
	}()

	var lastIndex uint64 = 0

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// 使用阻塞查询监听服务变化
			entries, meta, err := r.client.Health().Service(serviceName, "", true, &api.QueryOptions{
				WaitIndex: lastIndex,
				WaitTime:  time.Second * 30,
			})

			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			// 如果索引没变，说明服务列表没变化
			if meta.LastIndex == lastIndex {
				continue
			}
			lastIndex = meta.LastIndex

			// 转换服务信息
			instances := make([]ServiceInstance, 0, len(entries))
			for _, entry := range entries {
				instance := ServiceInstance{
					ID:        entry.Service.ID,
					Name:      entry.Service.Service,
					Address:   entry.Service.Address,
					Port:      entry.Service.Port,
					Metadata:  entry.Service.Meta,
					StartTime: time.Now(), // Consul 不提供启动时间，使用当前时间
				}

				// 尝试从元数据中获取 TTL
				if ttlStr, ok := entry.Service.Meta["ttl"]; ok {
					if ttl, err := strconv.Atoi(ttlStr); err == nil {
						instance.TTL = time.Duration(ttl) * time.Second
					}
				}

				instances = append(instances, instance)
			}

			// 发送服务列表更新
			select {
			case watchCh <- instances:
			case <-ctx.Done():
				return
			}
		}
	}
}

// Close 关闭连接
func (r *ConsulRegistry) Close() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 关闭所有监听通道
	for name, ch := range r.watchMap {
		close(ch)
		delete(r.watchMap, name)
	}

	return nil
}
