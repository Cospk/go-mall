package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"sync"
	"time"

	"github.com/go-zookeeper/zk"
)

// ZookeeperRegistry 是基于 Zookeeper 的服务注册与发现实现
type ZookeeperRegistry struct {
	conn      *zk.Conn
	basePath  string
	addresses []string
	timeout   time.Duration
	mutex     sync.RWMutex
	watchMap  map[string]chan []ServiceInstance
}

// NewZookeeperRegistry 创建一个新的 Zookeeper 注册中心
func NewZookeeperRegistry(addresses []string, timeout time.Duration) (*ZookeeperRegistry, error) {
	if len(addresses) == 0 {
		return nil, ErrInvalidConfig
	}

	conn, _, err := zk.Connect(addresses, timeout)
	if err != nil {
		return nil, fmt.Errorf("connect to zookeeper failed: %w", err)
	}

	r := &ZookeeperRegistry{
		conn:      conn,
		basePath:  "/services",
		addresses: addresses,
		timeout:   timeout,
		watchMap:  make(map[string]chan []ServiceInstance),
	}

	// 确保基础路径存在
	exists, _, err := conn.Exists(r.basePath)
	if err != nil {
		return nil, err
	}
	if !exists {
		_, err := conn.Create(r.basePath, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return nil, err
		}
	}

	return r, nil
}

// Register 注册服务实例
func (r *ZookeeperRegistry) Register(ctx context.Context, instance ServiceInstance) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	servicePath := path.Join(r.basePath, instance.Name)
	instancePath := path.Join(servicePath, instance.ID)

	// 确保服务路径存在
	exists, _, err := r.conn.Exists(servicePath)
	if err != nil {
		return ErrRegisterFailed
	}
	if !exists {
		_, err := r.conn.Create(servicePath, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return ErrRegisterFailed
		}
	}

	// 序列化服务信息
	data, err := json.Marshal(instance)
	if err != nil {
		return ErrRegisterFailed
	}

	// 创建服务实例节点
	exists, _, err = r.conn.Exists(instancePath)
	if err != nil {
		return ErrRegisterFailed
	}

	if exists {
		_, err = r.conn.Set(instancePath, data, -1)
	} else {
		// 使用临时节点，确保服务实例在连接断开时自动注销
		_, err = r.conn.Create(instancePath, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	}
	if err != nil {
		return ErrRegisterFailed
	}

	return nil
}

// Deregister 注销服务实例
func (r *ZookeeperRegistry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	instancePath := path.Join(r.basePath, serviceName, instanceID)

	exists, _, err := r.conn.Exists(instancePath)
	if err != nil {
		return ErrDeregisterFailed
	}

	if exists {
		err = r.conn.Delete(instancePath, -1)
		if err != nil {
			return ErrDeregisterFailed
		}
	}

	return nil
}

// GetService 获取指定服务的所有实例
func (r *ZookeeperRegistry) GetService(ctx context.Context, serviceName string) ([]ServiceInstance, error) {
	servicePath := path.Join(r.basePath, serviceName)

	exists, _, err := r.conn.Exists(servicePath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return []ServiceInstance{}, nil
	}

	// 获取所有服务实例
	children, _, err := r.conn.Children(servicePath)
	if err != nil {
		return nil, err
	}

	instances := make([]ServiceInstance, 0, len(children))
	for _, child := range children {
		instancePath := path.Join(servicePath, child)
		data, _, err := r.conn.Get(instancePath)
		if err != nil {
			continue
		}

		var instance ServiceInstance
		if err := json.Unmarshal(data, &instance); err != nil {
			continue
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// GetAllServices 获取所有服务
func (r *ZookeeperRegistry) GetAllServices(ctx context.Context) (map[string][]ServiceInstance, error) {
	// 获取所有服务名称
	services, _, err := r.conn.Children(r.basePath)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]ServiceInstance)
	for _, serviceName := range services {
		instances, err := r.GetService(ctx, serviceName)
		if err != nil {
			continue
		}
		result[serviceName] = instances
	}

	return result, nil
}

// Watch 监听服务变化
func (r *ZookeeperRegistry) Watch(ctx context.Context, serviceName string) (<-chan []ServiceInstance, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 检查是否已经在监听该服务
	if ch, ok := r.watchMap[serviceName]; ok {
		return ch, nil
	}

	servicePath := path.Join(r.basePath, serviceName)

	// 确保服务路径存在
	exists, _, err := r.conn.Exists(servicePath)
	if err != nil {
		return nil, err
	}
	if !exists {
		_, err := r.conn.Create(servicePath, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return nil, err
		}
	}

	// 创建通道
	watchCh := make(chan []ServiceInstance, 10)
	r.watchMap[serviceName] = watchCh

	// 启动监听协程
	go r.watchService(ctx, serviceName, watchCh)

	return watchCh, nil
}

// watchService 监听服务变化并发送更新
func (r *ZookeeperRegistry) watchService(ctx context.Context, serviceName string, watchCh chan []ServiceInstance) {
	defer func() {
		r.mutex.Lock()
		delete(r.watchMap, serviceName)
		r.mutex.Unlock()
		close(watchCh)
	}()

	servicePath := path.Join(r.basePath, serviceName)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// 获取当前服务列表并设置监听
			children, _, childrenWatcher, err := r.conn.ChildrenW(servicePath)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			instances := make([]ServiceInstance, 0, len(children))
			for _, child := range children {
				instancePath := path.Join(servicePath, child)
				data, _, err := r.conn.Get(instancePath)
				if err != nil {
					continue
				}

				var instance ServiceInstance
				if err := json.Unmarshal(data, &instance); err != nil {
					continue
				}
				instances = append(instances, instance)
			}

			// 发送服务列表更新
			select {
			case watchCh <- instances:
			case <-ctx.Done():
				return
			}

			// 等待变更事件
			select {
			case <-childrenWatcher:
			case <-ctx.Done():
				return
			}
		}
	}
}

// Close 关闭连接
func (r *ZookeeperRegistry) Close() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 关闭所有监听通道
	for name, ch := range r.watchMap {
		close(ch)
		delete(r.watchMap, name)
	}

	r.conn.Close()
	return nil
}
