# SQL 数据库迁移文件

本目录包含许可证管理系统的数据库结构定义和数据迁移文件。

## 目录结构

```
sql/
├── migrations/          # 数据库迁移文件
│   ├── 001_create_customers_table.sql       # 客户表结构
│   ├── 002_insert_sample_data.sql           # 示例数据
│   ├── 003_create_authorization_codes_table.sql  # 授权码表结构
│   ├── 004_create_licenses_table.sql        # 许可证表结构
│   ├── 005_create_authorization_changes_table.sql # 授权变更表结构
│   ├── 006_create_users_table.sql           # 管理员用户表结构
│   ├── 007_create_cu_users_table.sql        # 客户用户表结构
│   ├── 008_create_cu_orders_table.sql       # 客户用户订单表结构
│   ├── 009_update_authorization_codes_code_length.sql # 更新授权码字段长度
│   └── 010_update_cu_orders_authorization_code_length.sql # 更新订单表授权码字段长度
└── README.md           # 本文件
```

## 数据库支持

- **MySQL** (推荐版本 8.0+)
- 兼容 MariaDB (版本 10.3+)

## 迁移文件说明

### 001_create_customers_table.sql
创建客户相关的数据库结构，包括：

- **客户主表** (`customers`)：存储客户的所有基本信息，使用VARCHAR类型替代枚举
- **客户编码序列表** (`customer_code_sequence`)：用于自动生成客户编码
- **自动时间戳**：MySQL原生支持自动更新时间戳功能
- **索引优化**：为常用查询字段创建索引

### 002_insert_sample_data.sql
插入示例数据，包括：
- VIP企业客户
- 普通个人客户
- 政府客户
- 教育客户
- 禁用状态客户

### 008_create_cu_orders_table.sql
创建客户用户订单表结构，包括：

- **订单主表** (`cu_orders`)：存储客户用户购买套餐的订单信息
- **许可管理**：记录许可数量、价格、折扣等购买信息
- **状态跟踪**：订单状态从创建到支付完成的生命周期
- **授权关联**：存储生成的授权码和过期时间

### 009_update_authorization_codes_code_length.sql
更新授权码表结构以支持自包含配置授权码：

- **字段更新**：`code VARCHAR(100) → VARCHAR(1000)`
- **功能增强**：支持包含完整配置信息的自包含授权码
- **向后兼容**：保持对现有短授权码的支持
- **安全考虑**：RSA签名确保配置信息完整性

**执行注意事项**：
- 此迁移会修改现有表结构
- 生产环境执行前务必备份数据
- 现有短授权码继续有效
- 新授权码可包含完整配置信息

### 010_update_cu_orders_authorization_code_length.sql
更新客户订单表授权码字段长度以支持新的HMAC签名授权码：

- **字段更新**：`authorization_code VARCHAR(50) → VARCHAR(500)`
- **兼容性**：支持约150-200字符的新格式授权码
- **安全升级**：从RSA签名升级到HMAC签名，显著缩短授权码长度

**执行注意事项**：
- 扩展字段长度以适应新的授权码格式
- 不会影响现有数据
- 支持更高效的离线验证

## 执行顺序

按文件名的数字前缀顺序执行：

```bash
# MySQL
mysql -u username -p database_name < migrations/001_create_customers_table.sql
mysql -u username -p database_name < migrations/002_insert_sample_data.sql
```

## 客户编码规则

- 格式：`CUS-YYYY-NNNN`
- 示例：`CUS-2025-0001`
- 自动递增，按年度重置

## 主要功能特性

1. **自动时间戳**：创建和更新时间自动管理（使用MySQL原生功能）
2. **UUID主键**：使用UUID作为主键，保证全局唯一性
3. **字符串约束**：使用VARCHAR类型替代枚举，通过注释说明取值范围
4. **索引优化**：为常用查询字段创建索引
5. **简化设计**：移除复杂的触发器和函数，提高可维护性
6. **注释完整**：表和字段都有详细的中文注释

## 数据类型映射

| 设计文档字段 | MySQL类型 | 说明 |
|-------------|-----------|------|
| id | VARCHAR(36) | 主键，UUID格式 |
| customer_code | VARCHAR(20) | 唯一客户编码 |
| customer_type | VARCHAR(20) | 字符串：individual/enterprise/government/education |
| email | VARCHAR(255) | 唯一邮箱地址 |
| created_at | TIMESTAMP | 时间戳，自动生成 |
| status | VARCHAR(20) | 字符串：active/disabled |

## 注意事项

1. 执行前请确保MySQL版本支持 `UUID()` 函数（MySQL 8.0+）
2. 客户编码生成逻辑简化，建议在应用层实现自动生成
3. 字段约束通过应用层验证，确保数据一致性
4. 示例数据仅用于开发和测试，生产环境请谨慎使用