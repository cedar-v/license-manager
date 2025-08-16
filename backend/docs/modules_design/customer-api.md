# 客户管理API设计

## 1. API概览

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/customers` | 查询客户列表 |
| GET | `/api/customers/{id}` | 获取客户详情 |
| POST | `/api/customers` | 创建客户 |
| PUT | `/api/customers/{id}` | 更新客户信息 |
| DELETE | `/api/customers/{id}` | 删除客户(软删除) |
| PATCH | `/api/customers/{id}/status` | 修改客户状态 |

---

## 2. 数据模型

### 客户信息结构
```json
{
  "id": "uuid",
  "customer_code": "CUS-2025-0001",
  "customer_name": "北京科技有限公司",
  "customer_type": "enterprise",
  "customer_type_display": "企业客户",
  "contact_person": "张三",
  "contact_title": "技术总监",
  "email": "zhangsan@example.com",
  "phone": "13800138000",
  "address": "北京市朝阳区科技园区100号",
  "company_size": "medium",
  "company_size_display": "中型企业",
  "customer_level": "vip",
  "customer_level_display": "VIP客户",
  "status": "active",
  "status_display": "激活",
  "description": "重要客户，需要优先技术支持",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-20T15:45:00Z",
  "created_by": "admin_uuid",
  "updated_by": "admin_uuid"
}
```

### 字段说明
- `customer_type`: `individual` | `enterprise` | `government` | `education`
- `customer_type_display`: 客户类型的多语言显示名称
- `company_size`: `small` | `medium` | `large` | `enterprise`
- `company_size_display`: 企业规模的多语言显示名称
- `customer_level`: `normal` | `vip` | `enterprise` | `strategic`
- `customer_level_display`: 客户等级的多语言显示名称
- `status`: `active` | `disabled`
- `status_display`: 状态的多语言显示名称

> **注意**: `*_display` 字段根据请求头 `Accept-Language` 或查询参数 `lang` 返回对应语言的显示文本

---

## 3. API详细设计

### 3.1 查询客户列表
```
GET /api/customers
```

**Query参数**
- `page`: 页码，默认1
- `page_size`: 每页条数，默认20，最大100
- `search`: 搜索关键词(支持客户编码、名称、联系人、邮箱)
- `customer_type`: 客户类型筛选
- `customer_level`: 客户等级筛选
- `status`: 状态筛选
- `sort`: 排序字段，默认`created_at`
- `order`: 排序方向，默认`desc`

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": "uuid",
        "customer_code": "CUS-2025-0001",
        "customer_name": "北京科技有限公司",
        "customer_type": "enterprise",
        "customer_type_display": "企业客户",
        "contact_person": "张三",
        "email": "zhangsan@example.com",
        "customer_level": "vip",
        "customer_level_display": "VIP客户",
        "status": "active",
        "status_display": "激活",
        "created_at": "2025-01-15T10:30:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5
  }
}
```

### 3.2 获取客户详情
```
GET /api/customers/{id}
```

**响应示例**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    // 完整的客户信息结构
  }
}
```

### 3.3 创建客户
```
POST /api/customers
```

**请求体**
```json
{
  "customer_name": "北京科技有限公司",
  "customer_type": "enterprise",
  "contact_person": "张三",
  "contact_title": "技术总监",
  "email": "zhangsan@example.com",
  "phone": "13800138000",
  "address": "北京市朝阳区科技园区100号",
  "company_size": "medium",
  "customer_level": "normal",
  "description": "新客户",
  "status": "active"
}
```

**必填字段**
- `customer_name`: 客户名称
- `contact_person`: 联系人姓名
- `customer_level`: 客户等级
- `customer_type`: 客户类型
- `company_size`: 企业规模
- `status`: 状态


**响应示例**
```json
{
  "code": 0,
  "message": "客户创建成功",
  "data": {
    // 创建后的完整客户信息
  }
}
```

### 3.4 更新客户信息
```
PUT /api/customers/{id}
```

**请求体**: 与创建客户相同，但所有字段都是可选的

**响应示例**
```json
{
  "code": 0,
  "message": "客户信息更新成功",
  "data": {
    // 更新后的完整客户信息
  }
}
```

### 3.5 删除客户
```
DELETE /api/customers/{id}
```

**响应示例**
```json
{
  "code": 0,
  "message": "客户删除成功",
  "data": null
}
```

### 3.6 修改客户状态
```
PATCH /api/customers/{id}/status
```

**请求体**
```json
{
  "status": "disabled"
}
```

**响应示例**
```json
{
  "code": 0,
  "message": "客户状态更新成功",
  "data": {
    "id": "uuid",
    "status": "disabled"
  }
}
```
