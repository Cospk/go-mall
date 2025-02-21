# go-mall

本项目是一个基于 Go 语言和 Gin 框架的实战项目，涵盖了从基础环境搭建到高级功能实现的完整开发流程。项目通过模块化设计和分层架构，展示了如何构建一个健壮、可维护的 Web 应用程序。项目还提供了 Docker-Compose 和 Kubernetes 的部署方案，方便开发者快速搭建开发环境和生产环境。

## 项目结构
项目分为多个部分，涵盖了从基础配置到业务模块实现的完整流程：

### 第一部分：基础环境与项目初始化

> Go 基础环境搭建和 Gin 项目初始化

- 定制化项目配置：使用 Viper 管理项目配置，支持热加载

- 项目日志管理：集成 Zap 日志库并配置日志自动切割，然后封装日志自动为日志添加traceID和程序位置

- 全局中间件：实现请求日志记录、跨域处理、错误恢复等，保证项目健壮性和可观测性

- 错误处理：自定义错误类型，错误链条串联和发生位置记录，确保错误处理的一致性和可维护性

- 接口规范化：规范接口响应格式，统一错误码和错误响应，分页响应的标准化，方便前端调用

### 第二部分：项目架构与模块化设计

- 项目的软件分层设计和约定：定义好项目的分层架构和模块划分

- 路由的分模块管理：模块化管理

- GORM与日志整合：慢查询和数据库错误监控

- Redis 的封装和统一管理

- 业务模块划分与解耦

- 用 Option 模式和对接层规范化外部 API 的对接

阶段性总结：为了让项目好 Debug 我们做了这些事情
> 用Docker-Compose、K8s 两种方式快速给项目搭建一套开发环境-- MySQL 和 Redis

### 第三部分：用户认证体系

- 业务需求分析与模块划分：根据需求划分业务模块，确保功能清晰、职责单一

- 多平台用户认证体系：实现支持多平台登录、同平台登录互踢、Token 泄漏检测的认证系统

- Token管理：Token 的派发、存储和认证管理、刷新机制方式过期、防偷窃踢人下线

- 用户密码安全：使用加密书算法存储用户密码，用户注册、登录、登出功能的实现

- 自定义Error增强：扩展自定义 Error，支持错误解包和 errors.Is 判定

### 第四部分：用户与商品管理

- 用户个人信息管理：实现密码的安全修改和重置，对用户信息脱敏处理

- 商品模块：实现商品分类管理，支持分类的增删改查，提供商品列表分页查询、商品搜索和商品详情功能

- 购物车模块：添加、修改购物车、购物项列表和结算信息功能实现，使用 职责链模式 解耦商品满减和优惠逻辑，提升代码可扩展性

- 订单模块：

   创建订单、订单查询和取消的功能实现

  对接微信支付接口，演示支付流程。

  使用 **模板 + 策略模式** 实现多场景支付，支持灵活扩展。



### 第五部分：测试与部署

- 项目的单元测试：测试的基础搭建和数据库的 Mock 测试，对接口、方法、Package 的 Mock 测试

- 容器化与k8s部署：将应用打包为 Docker 镜像，支持快速部署和管理，在 Kubernetes 上部署应用实现平滑重启和安全调度


总结：怎么把项目扩展成微服务


使用到的组件和框架
Gin：轻量级 Web 框架

Viper：配置管理

Zap：高性能日志库

GORM：ORM 框架

go-redis：Redis 客户端

lo：Go 语言实用工具库