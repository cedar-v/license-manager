# SQL 数据库迁移文件

本目录包含许可证管理系统的数据库结构定义和数据迁移文件。

## 目录结构

```
sql/
├── migrations/          # 数据库迁移文件
│   ├── 001_create_customers_table.sql    # 客户表结构
│   └── 002_insert_sample_data.sql        # 示例数据
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