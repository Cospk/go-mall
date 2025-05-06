package rpc

import (
	"context"
	"fmt"
	pb "github.com/Cospk/go-mall/api/rpc/gen/go/user"
	"github.com/Cospk/go-mall/pkg/registry"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"sync"
	"time"
)

// RPC请求/响应结构

// UserServiceClient 定义用户服务客户端接口
type UserServiceClient interface {
	// Register 注册用户
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	// Login 登录
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	// Logout 登出
	Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error)
	// GetUserInfo 获取用户信息
	GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error)
	// RefreshUserToken 刷新用户token
	RefreshUserToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error)
	// PasswordResetApply 申请重置登录密码
	PasswordResetApply(ctx context.Context, req *pb.PasswordResetApplyRequest) (*pb.PasswordResetApplyResponse, error)
	// PasswordReset 重置登录密码
	PasswordReset(ctx context.Context, req *pb.PasswordResetRequest) (*pb.PasswordResetResponse, error)
	// GetUserInfoById 获取用户信息-ID
	GetUserInfoById(ctx context.Context, req *pb.GetUserInfoByIdRequest) (*pb.GetUserInfoResponse, error)
	// UpdateUserInfo 更新用户信息
	UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoResponse, error)
	// AddUserAddressInfo 添加用户收货地址
	AddUserAddressInfo(ctx context.Context, req *pb.AddUserAddressInfoRequest) (*pb.AddUserAddressInfoResponse, error)
	// GetUserAddressList 获取用户收货地址列表
	GetUserAddressList(ctx context.Context, req *pb.GetUserAddressListRequest) (*pb.GetUserAddressListResponse, error)
	// GetUserAddressInfo 获取单个收货地址信息
	GetUserAddressInfo(ctx context.Context, req *pb.GetUserAddressInfoRequest) (*pb.UserAddressInfo, error)
	// UpdateUserAddressInfo 更新用户收货地址信息
	UpdateUserAddressInfo(ctx context.Context, req *pb.UpdateUserAddressInfoRequest) (*pb.UpdateUserAddressInfoResponse, error)
	// DeleteUserAddressInfo 删除用户收货地址
	DeleteUserAddressInfo(ctx context.Context, req *pb.DeleteUserAddressInfoRequest) (*pb.DeleteUserAddressInfoResponse, error)
	// LoadBalanceTest 负载均衡测试
}

// 用户服务客户端实现
type userServiceClient struct {
	conn          *grpc.ClientConn
	client        pb.StreamGreeterClient
	registry      registry.Registry
	serviceName   string
	instances     []registry.ServiceInstance
	mu            sync.RWMutex
	watchCanceled context.CancelFunc
}

func NewUserServiceClient() UserServiceClient {
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
		return &userServiceClient{
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
		return &userServiceClient{
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

	client := &userServiceClient{
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
func (c *userServiceClient) watchService(ctx context.Context) {
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
func (c *userServiceClient) reconnect() {
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

// 实现接口方法

func (c *userServiceClient) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return c.client.Register(ctx, req)
}

func (c *userServiceClient) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return c.client.LoginUser(ctx, req)
}

func (c *userServiceClient) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	return c.client.LogoutUser(ctx, req)
}

func (c *userServiceClient) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	return c.client.GetUserInfo(ctx, req)
}

func (c *userServiceClient) RefreshUserToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	return c.client.RefreshUserToken(ctx, req)
}

func (c *userServiceClient) PasswordResetApply(ctx context.Context, req *pb.PasswordResetApplyRequest) (*pb.PasswordResetApplyResponse, error) {
	return c.client.PasswordResetApply(ctx, req)
}

func (c *userServiceClient) PasswordReset(ctx context.Context, req *pb.PasswordResetRequest) (*pb.PasswordResetResponse, error) {
	return c.client.PasswordReset(ctx, req)
}

func (c *userServiceClient) GetUserInfoById(ctx context.Context, req *pb.GetUserInfoByIdRequest) (*pb.GetUserInfoResponse, error) {
	return c.client.GetUserInfoById(ctx, req)
}

func (c *userServiceClient) UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoResponse, error) {
	return c.client.UpdateUserInfo(ctx, req)
}

func (c *userServiceClient) AddUserAddressInfo(ctx context.Context, req *pb.AddUserAddressInfoRequest) (*pb.AddUserAddressInfoResponse, error) {
	return c.client.AddUserAddressInfo(ctx, req)
}

func (c *userServiceClient) GetUserAddressList(ctx context.Context, req *pb.GetUserAddressListRequest) (*pb.GetUserAddressListResponse, error) {
	return c.client.GetUserAddressList(ctx, req)
}

func (c *userServiceClient) GetUserAddressInfo(ctx context.Context, req *pb.GetUserAddressInfoRequest) (*pb.UserAddressInfo, error) {
	return c.client.GetUserAddressInfo(ctx, req)
}

func (c *userServiceClient) UpdateUserAddressInfo(ctx context.Context, req *pb.UpdateUserAddressInfoRequest) (*pb.UpdateUserAddressInfoResponse, error) {
	return c.client.UpdateUserAddressInfo(ctx, req)
}

func (c *userServiceClient) DeleteUserAddressInfo(ctx context.Context, req *pb.DeleteUserAddressInfoRequest) (*pb.DeleteUserAddressInfoResponse, error) {
	return c.client.DeleteUserAddressInfo(ctx, req)
}
