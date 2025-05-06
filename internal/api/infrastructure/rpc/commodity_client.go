package rpc

import (
	"context"
	"fmt"
	pb "github.com/Cospk/go-mall/api/rpc/gen/go/commodity"
	"github.com/Cospk/go-mall/pkg/registry"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"sync"
	"time"
)

// CommodityServiceClient 商品服务客户端接口
type CommodityServiceClient interface {
	// 获取分类层级结构
	GetCategoryHierarchy(ctx context.Context, req *pb.EmptyRequest) (*pb.CategoryHierarchyReply, error)
	// 根据父ID获取分类列表
	GetCategoriesWithParentId(ctx context.Context, req *pb.ParentIdRequest) (*pb.CategoriesReply, error)
	// 获取分类下的商品列表
	CommoditiesInCategory(ctx context.Context, req *pb.CategoryCommoditiesRequest) (*pb.CommoditiesReply, error)
	// 商品搜索
	CommoditySearch(ctx context.Context, req *pb.SearchRequest) (*pb.CommoditiesReply, error)
	// 获取商品详情
	CommodityInfo(ctx context.Context, req *pb.CommodityIdRequest) (*pb.CommodityDetailReply, error)
}
type commodityServiceClient struct {
	conn          *grpc.ClientConn
	client        pb.StreamGreeterClient
	registry      registry.Registry
	serviceName   string
	instances     []registry.ServiceInstance
	mu            sync.RWMutex
	watchCanceled context.CancelFunc
}

func NewCommodityServiceClient() CommodityServiceClient {
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
		return &commodityServiceClient{
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
		return &commodityServiceClient{
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

	client := &commodityServiceClient{
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
func (c *commodityServiceClient) watchService(ctx context.Context) {
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
func (c *commodityServiceClient) reconnect() {
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

// 获取分类层级结构
func (c *commodityServiceClient) GetCategoryHierarchy(ctx context.Context, req *pb.EmptyRequest) (*pb.CategoryHierarchyReply, error) {
	return c.client.GetCategoryHierarchy(ctx, req)
}

// 根据父ID获取分类列表
func (c *commodityServiceClient) GetCategoriesWithParentId(ctx context.Context, req *pb.ParentIdRequest) (*pb.CategoriesReply, error) {
	return c.client.GetCategoriesWithParentId(ctx, req)
}

// 获取分类下的商品列表
func (c *commodityServiceClient) CommoditiesInCategory(ctx context.Context, req *pb.CategoryCommoditiesRequest) (*pb.CommoditiesReply, error) {
	return c.client.CommoditiesInCategory(ctx, req)
}

// 商品搜索
func (c *commodityServiceClient) CommoditySearch(ctx context.Context, req *pb.SearchRequest) (*pb.CommoditiesReply, error) {
	return c.client.CommoditySearch(ctx, req)
}

// 获取商品详情
func (c *commodityServiceClient) CommodityInfo(ctx context.Context, req *pb.CommodityIdRequest) (*pb.CommodityDetailReply, error) {
	return c.client.CommodityInfo(ctx, req)
}
