# 后端架构设计

## 技术栈选择

  Web框架: Gin
  ORM: GORM
  配置: Viper
  日志: logrus
  认证: JWT
  缓存: Redis (可选)
  数据库: PostgreSQL/MySQL

## 目录结构

基于Gin + GORM技术栈，采用**Clean Architecture**风格的目录结构：

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # 程序入口
├── internal/
│   ├── api/
│   │   ├── handlers/            # HTTP处理器
│   │   │   ├── auth.go
│   │   │   ├── customer.go
│   │   │   ├── license.go
│   │   │   └── system.go
│   │   ├── middleware/          # 中间件
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   └── logging.go
│   │   └── routes/              # 路由定义
│   │       └── router.go
│   ├── service/                 # 业务逻辑层
│   │   ├── auth_service.go
│   │   ├── customer_service.go
│   │   ├── license_service.go
│   │   └── interfaces.go        # 接口定义
│   ├── repository/              # 数据访问层
│   │   ├── customer_repo.go
│   │   ├── license_repo.go
│   │   └── interfaces.go
│   ├── models/                  # 数据模型
│   │   ├── customer.go
│   │   ├── license.go
│   │   └── common.go
│   ├── config/                  # 配置管理
│   │   └── config.go
│   └── cache/                   # 缓存层
│       ├── memory.go
│       └── redis.go
├── pkg/
│   ├── utils/                   # 工具函数
│   │   ├── crypto.go
│   │   ├── hardware.go
│   │   └── validator.go
│   └── logger/                  # 日志封装
│       └── logger.go
├── configs/
│   ├── config.yaml
│   └── config.example.yaml
├── migrations/                  # 数据库迁移
├── tools/                       # 硬件信息工具
└── go.mod
```

**核心设计原则：**
- **分层架构**：Handler → Service → Repository
- **依赖注入**：通过接口解耦各层
- **配置驱动**：Viper管理所有配置
- **缓存抽象**：支持内存/Redis切换

