package rpc

import (
	"context"
	"fmt"
	pb "github.com/Cospk/go-mall/api/rpc/gen/go/order"
	"github.com/Cospk/go-mall/pkg/registry"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"sync"
	"time"
)

// OrderServiceClient 订单服务客户端接口
type OrderServiceClient interface {
	// 创建订单
	CreateOrder(ctx context.Context, req *pb.OrderCreateRequest) (*pb.OrderCreateReply, error)
	// 获取用户订单列表
	GetUserOrders(ctx context.Context, req *pb.UserOrdersRequest) (*pb.UserOrdersReply, error)
	// 获取订单详情
	GetOrderInfo(ctx context.Context, req *pb.OrderInfoRequest) (*pb.OrderInfoReply, error)
	// 取消订单
	CancelOrder(ctx context.Context, req *pb.OrderCancelRequest) (*pb.CommonResponse, error)
	// 创建订单支付
	OrderCreatePay(ctx context.Context, req *pb.OrderPayCreateRequest) (*pb.OrderPayCreateReply, error)
}

type orderServiceClient struct {
	conn          *grpc.ClientConn
	client        pb.StreamGreeterClient
	registry      registry.Registry
	serviceName   string
	instances     []registry.ServiceInstance
	mu            sync.RWMutex
	watchCanceled context.CancelFunc
}

func NewOrderServiceClient() OrderServiceClient {

	// 从配置文件中获取服务名称和注册中心配置
	serviceName := viper.GetString("service.order.name")
	if serviceName == "" {
		serviceName = "order-service"
	}
	// 创建注册中心客户端
	registryType := viper.GetString("registry.type")
	if registryType == "" {
		registryType = "etcd"
	}

	endpoints := viper.GetStringSlice("registry.endpoints")
	if len(endpoints) == 0 {
		endpoints = []string{"localhost:2379"}
	}
	timeout := viper.GetInt("registry.timeout")
	if timeout == 0 {
		timeout = 5
	}
	prefix := viper.GetString("registry.prefix")
	if prefix == "" {
		prefix = "/services/"
	}
	config := registry.RegistryConfig{
		Type:      registry.RegistryType(registryType),
		Endpoints: endpoints,
		Timeout:   timeout,
		Prefix:    prefix,
	}
	ctx, cancel := context.WithCancel(context.Background())
	reg, err := registry.NewRegistry(ctx, config)
	if err != nil {
		// 如果服务发现失败，回退到直连模式
		fmt.Printf("服务发现初始化失败: %v, 回退到直连模式\n", err)
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		return &orderServiceClient{
			conn:          conn,
			client:        pb.NewStreamGreeterClient(conn),
			watchCanceled: cancel,
		}
	}

	// 获取服务实例
	instances, err := reg.GetService(ctx, serviceName)
	if err != nil || len(instances) == 0 {
		// 如果获取服务实例失败，回退到直连模式
		fmt.Printf("获取服务实例失败: %v, 回退到直连模式\n", err)
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		return &orderServiceClient{
			conn:          conn,
			client:        pb.NewStreamGreeterClient(conn),
			registry:      reg,
			serviceName:   serviceName,
			watchCanceled: cancel,
		}
	}

	// 随机选择一个实例
	rand.Seed(time.Now().UnixNano())
	instance := instances[rand.Intn(len(instances))]
	addr := fmt.Sprintf("%s:%d", instance.Address, instance.Port)

	// 建立gRPC连接
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := &orderServiceClient{
		conn:          conn,
		client:        pb.NewStreamGreeterClient(conn),
		registry:      reg,
		serviceName:   serviceName,
		instances:     instances,
		watchCanceled: cancel,
	}

	//
	go client.watchService(ctx)

	return client
}

// 监听服务变化
func (c *orderServiceClient) watchService(ctx context.Context) {
	if c.registry == nil {
		return
	}

	watchCh, err := c.registry.Watch(ctx, c.serviceName)
	if err != nil {
		fmt.Printf("监听服务变化失败: %v\n", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case instances, ok := <-watchCh:
			if !ok {
				return
			}

			if len(instances) == 0 {
				continue
			}

			c.mu.Lock()
			oldInstances := c.instances
			c.instances = instances
			c.mu.Unlock()

			// 如果实例列表变化，可以考虑重新连接
			if !instancesEqual(oldInstances, instances) {
				c.reconnect()
			}
		}
	}
}

// 重新连接到服务
func (c *orderServiceClient) reconnect() {
	c.mu.RLock()
	instances := c.instances
	c.mu.RUnlock()

	if len(instances) == 0 {
		return
	}

	// 随机选择一个实例
	rand.Seed(time.Now().UnixNano())
	instance := instances[rand.Intn(len(instances))]
	addr := fmt.Sprintf("%s:%d", instance.Address, instance.Port)

	// 创建新连接
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("重新连接服务失败: %v\n", err)
		return
	}

	// 关闭旧连接
	if c.conn != nil {
		c.conn.Close()
	}

	c.mu.Lock()
	c.conn = conn
	c.client = pb.NewStreamGreeterClient(conn)
	c.mu.Unlock()
}

// 比较两个实例列表是否相同
func instancesEqual(a, b []registry.ServiceInstance) bool {
	if len(a) != len(b) {
		return false
	}

	// 简单比较，只检查长度变化
	// 实际应用中可能需要更复杂的比较逻辑
	return true
}

//func (c *orderServiceClient) Close() error {
//	if c.watchCanceled != nil {
//		c.watchCanceled()
//	}
//
//	if c.registry != nil {
//		c.registry.Close()
//	}
//
//	if c.conn != nil {
//		return c.conn.Close()
//	}
//
//	return nil
//}

// 创建订单
func (c *orderServiceClient) CreateOrder(ctx context.Context, req *pb.OrderCreateRequest) (*pb.OrderCreateReply, error) {

	return c.client.CreateOrder(ctx, req)
}

// 获取用户订单列表
func (c *orderServiceClient) GetUserOrders(ctx context.Context, req *pb.UserOrdersRequest) (*pb.UserOrdersReply, error) {
	return c.client.GetUserOrders(ctx, req)
}

// 获取订单详情
func (c *orderServiceClient) GetOrderInfo(ctx context.Context, req *pb.OrderInfoRequest) (*pb.OrderInfoReply, error) {
	return c.client.GetOrderInfo(ctx, req)
}

// 取消订单
func (c *orderServiceClient) CancelOrder(ctx context.Context, req *pb.OrderCancelRequest) (*pb.CommonResponse, error) {
	return c.client.CancelOrder(ctx, req)
}

// 创建订单支付
func (c *orderServiceClient) OrderCreatePay(ctx context.Context, req *pb.OrderPayCreateRequest) (*pb.OrderPayCreateReply, error) {
	return c.client.OrderCreatePay(ctx, req)
}
