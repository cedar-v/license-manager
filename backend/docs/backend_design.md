# 后端架构设计

## 技术栈选择

  Web框架: Gin
  ORM: GORM,开发环境使用Auto Migration
  配置: Viper
  日志: logrus
  认证: JWT (支持管理员和C端用户双JWT体系)
  缓存: Redis (可选)
  数据库: PostgreSQL/MySQL
  支付集成: 支付宝
  加密算法: RSA
  

## 目录结构


基于Gin + GORM技术栈，采用**Clean Architecture**风格的目录结构：

```
backend/
├── cmd/
│   ├── main.go                  # 程序入口
│   ├── license_manager.sh       # 启动脚本
│   ├── client-demo/             # C端客户端演示程序
│   │   ├── main.go
│   │   ├── rsa.go
│   │   ├── client_config.json
│   │   ├── README.md
│   │   └── license_code/
│   └── gen-rsa-keys/            # RSA密钥生成工具
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/            # HTTP处理器
│   │   │   ├── auth_handler.go
│   │   │   ├── authorization_code_handler.go
│   │   │   ├── customer_handler.go
│   │   │   ├── dashboard_handler.go
│   │   │   ├── enum_handler.go
│   │   │   ├── license_handler.go
│   │   │   ├── payment_handler.go
│   │   │   ├── system_handler.go
│   │   │   ├── cu_auth_handler.go
│   │   │   ├── cu_authorization_handler.go
│   │   │   ├── cu_device_handler.go
│   │   │   ├── cu_order_handler.go
│   │   │   ├── cu_order_handler_test.go
│   │   │   └── cu_profile_handler.go
│   │   ├── middleware/          # 中间件
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   ├── i18n.go
│   │   │   └── logging.go
│   │   └── routes/              # 路由定义
│   │       └── router.go
│   ├── service/                 # 业务逻辑层
│   │   ├── auth_service.go
│   │   ├── authorization_code_service.go
│   │   ├── customer_service.go
│   │   ├── dashboard_service.go
│   │   ├── enum_service.go
│   │   ├── license_service.go
│   │   ├── payment_service.go
│   │   ├── system_service.go
│   │   ├── cu_device_service.go
│   │   ├── cu_order_service.go
│   │   ├── cu_user_service.go
│   │   └── interfaces.go        # 接口定义
│   ├── repository/              # 数据访问层
│   │   ├── authorization_code_repository.go
│   │   ├── customer_repository.go
│   │   ├── dashboard_repository.go
│   │   ├── license_repository.go
│   │   ├── payment_repository.go
│   │   ├── user_repository.go
│   │   ├── cu_order_repository.go
│   │   ├── cu_user_repository.go
│   │   ├── gorm/                 # GORM相关扩展
│   │   ├── errors.go
│   │   └── interfaces.go
│   ├── models/                  # 数据模型
│   │   ├── auth.go
│   │   ├── common.go
│   │   ├── customer.go
│   │   ├── cu_order.go
│   │   ├── cu_user.go
│   │   ├── dashboard.go
│   │   ├── license.go
│   │   ├── payment.go
│   │   ├── system.go
│   │   └── user.go
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
│   │   ├── alipay.go            # 支付宝支付集成
│   │   ├── crypto.go            # 加密工具
│   │   ├── cu_jwt.go            # C端用户JWT工具
│   │   ├── jwt.go               # JWT工具
│   │   └── license_codec.go     # 授权码编解码
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
│   │   ├── auth_api.md
│   │   ├── customer-api.md
│   │   ├── customer.md
│   │   ├── login&auth.md
│   │   ├── user.md
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
│   ├── 003_create_authorization_codes_table.sql
│   ├── 004_create_licenses_table.sql
│   ├── 005_create_authorization_changes_table.sql
│   ├── 006_create_users_table.sql
│   ├── 007_create_cu_users_table.sql
│   ├── 008_create_cu_orders_table.sql
│   ├── 009_update_authorization_codes_code_length.sql
│   ├── 010_update_cu_orders_authorization_code_length.sql
│   ├── 011_create_payments_table.sql
│   └── README.md
├── go.mod
└── go.sum
```

**核心设计原则：**
- **分层架构**：Handler → Service → Repository
- **依赖注入**：通过接口解耦各层
- **配置驱动**：Viper管理所有配置
- **缓存抽象**：支持内存/Redis切换
- **双JWT体系**：支持管理员和C端用户独立认证体系
- **模块化设计**：支持客户管理、订单管理、支付集成等独立模块
- **国际化支持**：完整的多语言错误信息系统

