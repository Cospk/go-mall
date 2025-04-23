/*
Package registry 提供了一个统一的服务注册与发现接口，支持多种注册中心实现。

该包的主要目标是提供一个抽象层，使应用程序可以轻松地在不同的服务注册中心之间切换，
而无需修改业务代码。目前支持的注册中心包括：
  - etcd
  - consul (计划中)
  - zookeeper (计划中)

基本使用示例:

	// 创建注册中心客户端
	reg, err := registry.NewRegistry(registry.RegistryConfig{
		Type:      registry.RegistryTypeEtcd,
		Endpoints: []string{"localhost:2379"},
		Timeout:   5,
		Prefix:    "/services/",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer reg.Close()

	// 注册服务
	err = reg.Register(registry.ServiceInstance{
		ID:       "service-1",
		Name:     "user-service",
		Address:  "localhost",
		Port:     8080,
		Metadata: map[string]string{"version": "1.0.0"},
		TTL:      time.Second * 30,
	})

	// 发现服务
	instances, err := reg.GetService("user-service")

设计原则:

1. 接口分离原则：核心功能通过 Registry 接口定义，与具体实现分离
2. 单一职责原则：每个注册中心实现只负责与特定服务的交互
3. 开闭原则：可以轻松添加新的注册中心实现，而无需修改现有代码
4. 依赖注入：通过工厂方法创建具体实现，降低耦合度

主要组件:

- Registry: 定义服务注册与发现的核心接口
- ServiceInstance: 表示一个服务实例的结构体
- EtcdRegistry/ConsulRegistry/ZookeeperRegistry: 特定注册中心的实现
- NewRegistry: 工厂方法，根据配置创建具体的注册中心实现

错误处理:

包中定义了一系列标准错误，如 ErrNotFound、ErrRegisterFailed 等，
用户可以通过这些错误进行错误处理和日志记录。

线程安全:

所有实现都保证线程安全，可以在多个 goroutine 中并发使用。
*/

package registry
