# 后端架构设计

## 技术栈选择

  Web框架: Gin
  ORM: GORM,开发环境使用Auto Migration
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
│   ├── main.go                  # 程序入口
│   └── license_manager.sh       # 启动脚本
├── internal/
│   ├── api/
│   │   ├── handlers/            # HTTP处理器
│   │   │   ├── auth_handler.go
│   │   │   ├── customer_handler.go
│   │   │   ├── enum_handler.go
│   │   │   └── system_handler.go
│   │   ├── middleware/          # 中间件
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   ├── i18n.go
│   │   │   └── logging.go
│   │   └── routes/              # 路由定义
│   │       └── router.go
│   ├── service/                 # 业务逻辑层
│   │   ├── auth_service.go
│   │   ├── customer_service.go
│   │   ├── enum_service.go
│   │   ├── system_service.go
│   │   └── interfaces.go        # 接口定义
│   ├── repository/              # 数据访问层
│   │   ├── customer_repository.go
│   │   ├── errors.go
│   │   └── interfaces.go
│   ├── models/                  # 数据模型
│   │   ├── auth.go
│   │   ├── customer.go
│   │   ├── system.go
│   │   └── common.go
│   ├── config/                  # 配置管理
│   │   └── config.go
│   └── database/                # 数据库连接和迁移
│       ├── connection.go
│       └── migration.go
├── pkg/
│   ├── context/                 # 上下文封装
│   │   └── context.go
│   ├── i18n/                    # 国际化
│   │   ├── errors.go
│   │   └── manager.go
│   ├── utils/                   # 工具函数
│   │   ├── crypto.go
│   │   └── jwt.go
│   └── logger/                  # 日志封装
│       └── logger.go
├── configs/
│   ├── config.dev.yaml          # 开发环境配置
│   ├── config.prod.yaml         # 生产环境配置
│   ├── config.yaml              # 默认配置
│   ├── config.example.yaml      # 配置示例
│   └── i18n/                    # 国际化配置
│       └── errors/
│           ├── en-US.yaml
│           ├── ja-JP.yaml
│           └── zh-CN.yaml
├── docs/                        # 文档目录
│   ├── api_development_guide.md
│   ├── backend_design.md
│   ├── design_philosophy.md
│   ├── error_i18n_design.md
│   ├── install/
│   │   ├── db_install.md
│   │   └── swagger_install.md
│   ├── modules_design/
│   │   ├── customer-api.md
│   │   ├── customer.md
│   │   ├── login&auth.md
│   │   ├── user_story.md
│   │   ├── 产品原型设计方案.md
│   │   └── 授权模块属性设计.md
│   └── swagger/                 # Swagger文档
│       ├── docs.go
│       ├── swagger.json
│       └── swagger.yaml
├── migrations/                  # 数据库迁移
│   ├── 001_create_customers_table.sql
│   ├── 002_insert_sample_data.sql
│   └── README.md
├── go.mod
└── go.sum
```

**核心设计原则：**
- **分层架构**：Handler → Service → Repository
- **依赖注入**：通过接口解耦各层
- **配置驱动**：Viper管理所有配置
- **缓存抽象**：支持内存/Redis切换

