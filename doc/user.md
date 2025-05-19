## 项目结构

```bash
cmd/
  └── user/
      └── main.go            # 主程序入口 
config/
  └── config.yaml            # 配置文件
internal/
  └── user/
      ├── domain/               # 领域模型
      │   ├── entity/           # 实体
      │   │   ├── user.go       # 用户实体
      │   │   └── address.go    # 地址实体
      │   └── repository/       # 仓储接口
      │       ├── user_repo.go  # 用户仓储接口
      │       └── address_repo.go # 地址仓储接口
      ├── infrastructure/       # 基础设施
      │   ├── persistence/      # 持久化
      │   │   ├── mysql/        # MySQL实现
      │   │   │   ├── user_repo.go  # 用户仓储MySQL实现
      │   │   │   └── address_repo.go # 地址仓储MySQL实现
      │   │   └── entity/        # 数据库模型
      │   │       ├── user.go   # 用户数据库模型
      │   │       └── address.go # 地址数据库模型
      │   └── rpc/              # RPC服务
      │       └── server.go     # gRPC服务器
      ├── application/          # 应用层
      │   ├── service/          # 服务
      │   │   ├── user_service.go  # 用户服务
      │   │   └── address_service.go # 地址服务
      │   └── dto/              # 数据传输对象
      │       ├── user_dto.go   # 用户DTO
      │       └── address_dto.go # 地址DTO
      └── interfaces/           # 接口层
          └── rpc/              # RPC接口
              └── user_server.go # 用户RPC服务实现
```