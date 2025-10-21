# 授权模块API设计

## 1. 授权码管理 API

### 1.1 授权码列表
```http
GET /api/v1/authorization-codes
```

**查询参数**
- `customer_id` - 客户ID
- `status` - 状态 (normal/locked/expired)
- `page` - 页码
- `page_size` - 每页数量（默认20，最大100）
- `start_date` - 创建开始时间
- `end_date` - 创建结束时间
- `sort` - 排序字段 (created_at/updated_at/code)
- `order` - 排序方向 (asc/desc)

**实现说明**

由于 `status` 是虚拟字段，需要在 SQL 查询中动态计算并筛选：

```sql
-- 状态计算逻辑
SELECT *,
  CASE 
    WHEN is_locked = true THEN 'locked'
    WHEN end_date < NOW() THEN 'expired'
    WHEN start_date <= NOW() AND end_date >= NOW() THEN 'normal'
    ELSE 'expired'
  END AS status
FROM authorization_codes
WHERE 1=1
  -- status 参数筛选
  AND (
    (:status = 'locked' AND is_locked = true) OR
    (:status = 'expired' AND is_locked = false AND end_date < NOW()) OR
    (:status = 'normal' AND is_locked = false AND start_date <= NOW() AND end_date >= NOW()) OR
    :status IS NULL
  )
  -- 其他筛选条件
  AND (:customer_id IS NULL OR customer_id = :customer_id)
  AND (:start_date IS NULL OR created_at >= :start_date)
  AND (:end_date IS NULL OR created_at <= :end_date)
ORDER BY created_at DESC
LIMIT :page_size OFFSET :offset;
```

**响应**
```json
{
  "code": "000000",
  "message": "成功",
  "data": {
    "list": [
      {
        "id": "uuid-string",
        "code": "LIC-COMP001-A7B9X2-C8F4",
        "customer_id": "customer-uuid",
        "customer_name": "张三公司",
        "customer_name_display": "张三公司",
        "status": "normal",
        "status_display": "正常",
        "start_date": "2024-01-01T00:00:00Z",
        "end_date": "2024-12-31T23:59:59Z",
        "max_activations": 10,
        "current_activations": 7,
        "deployment_type": "standalone",
        "deployment_type_display": "单机版",
        "is_locked": false,
        "description": "企业版授权",
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5
  }
}
```

### 1.2 创建授权码
```http
POST /api/v1/authorization-codes
```

**说明：** 有效期从当天00:00:00开始，到指定天数后的23:59:59结束。例如：validity_days=1表示今天一整天有效。

**请求体**
```json
{
  "customer_id": "customer-uuid-123",
  "software_id": "erp-system",
  "description": "企业版授权",
  "validity_days": 365,
  "deployment_type": "standalone",
  "encryption_type": "standard",
  "software_version": "1.0.0",
  "max_activations": 10,
  "feature_config": {
    "modules": ["user_mgmt", "inventory", "finance"]
  },
  "usage_limits": {
    "max_users": 100,
    "api_calls_per_day": 10000
  },
  "custom_parameters": {}
}
```

**参数说明：**
- `customer_id`: 客户ID（必填）
- `software_id`: 软件产品ID（可选）
- `description`: 授权描述（可选）
- `validity_days`: 有效期天数，范围1-36500天（必填）
- `deployment_type`: 部署类型，枚举值：standalone/cloud/hybrid（必填）
- `encryption_type`: 加密类型，枚举值：standard/advanced（可选，默认standard）
- `software_version`: 软件版本（可选）
- `max_activations`: 最大激活数量（必填，最小值1）
- `feature_config`: 功能配置JSON（可选）
- `usage_limits`: 使用限制JSON（可选）
- `custom_parameters`: 自定义参数JSON（可选）

**响应**
```json
{
  "code": "000000",
  "message": "成功",
  "data": {
    "id": "uuid-string",
    "code": "LIC-COMP001-A7B9X2-C8F4"
  }
}
```

### 1.3 授权码详情
```http
GET /api/v1/authorization-codes/{id}
```

**响应**
```json
{
  "code": "000000",
  "message": "成功",
  "data": {
    "id": "uuid-string",
    "code": "LIC-COMP001-A7B9X2-C8F4",
    "customer_id": "customer-uuid",
    "customer_name": "张三公司",
    "status": "normal",
    "status_display": "正常",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-12-31T23:59:59Z",
    "max_activations": 10,
    "current_activations": 7,
    "deployment_type": "standalone",
    "deployment_type_display": "单机版",
    "encryption_type": "standard",
    "encryption_type_display": "标准加密",
    "feature_config": {
      "modules": ["user_mgmt", "inventory", "finance"]
    },
    "usage_limits": {
      "max_users": 100,
      "api_calls_per_day": 10000
    },
    "custom_parameters": {},
    "is_locked": false,
    "lock_reason": null,
    "description": "企业版授权",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

### 1.4 更新授权码
```http
PUT /api/v1/authorization-codes/{id}
```

**请求体**
```json
{
  "start_date": "2024-01-01",
  "end_date": "2025-12-31",
  "max_activations": 20,
  "feature_config": {
    "modules": ["user_mgmt", "inventory", "finance", "crm"]
  },
  "usage_limits": {
    "max_users": 200,
    "api_calls_per_day": 20000
  },
  "change_type": "upgrade",
  "reason": "客户升级套餐"
}
```

### 1.5 锁定/解锁授权码
```http
PUT /api/v1/authorization-codes/{id}/lock
```

**请求体**
```json
{
  "is_locked": true,
  "lock_reason": "违规使用"
}
```

### 1.6 删除授权码
```http
DELETE /api/v1/authorization-codes/{id}
```

## 2. 许可证管理 API

### 2.1 许可证列表
```http
GET /api/v1/licenses
```

**查询参数**
- `authorization_code_id` - 授权码ID
- `customer_id` - 客户ID
- `status` - 状态 (active/inactive/revoked)
- `is_online` - 在线状态
- `page` - 页码（默认1）
- `page_size` - 每页数量（默认20，最大100）
- `sort` - 排序字段 (created_at/updated_at/activated_at/last_heartbeat)
- `order` - 排序方向 (asc/desc)

**响应**
```json
{
  "code": "000000",
  "message": "成功",
  "data": {
    "list": [
      {
        "id": "license-uuid",
        "license_key": "LIC-DEVICE-ABC123456789",
        "authorization_code_id": "auth-code-uuid",
        "authorization_code": "LIC-COMP001-A7B9X2-C8F4",
        "customer_name": "张三公司",
        "hardware_fingerprint": "CPU:ABC123,MB:DEF456",
        "status": "active",
        "status_display": "激活",
        "is_online": true,
        "is_online_display": "在线",
        "activation_ip": "192.168.1.100",
        "last_online_ip": "192.168.1.100",
        "activated_at": "2024-01-01T10:00:00Z",
        "last_heartbeat": "2024-01-01T14:30:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20,
    "total_pages": 3
  }
}
```

### 2.2 许可证详情
```http
GET /api/v1/licenses/{id}
```

**响应**
```json
{
  "code": "000000",
  "message": "成功",
  "data": {
    "id": "license-uuid",
    "license_key": "LIC-DEVICE-ABC123456789",
    "authorization_code_id": "auth-code-uuid",
    "authorization_code": "LIC-COMP001-A7B9X2-C8F4",
    "customer_id": "customer-uuid",
    "customer_name": "张三公司",
    "hardware_fingerprint": "CPU:ABC123,MB:DEF456,MAC:00:11:22:33:44:55",
    "device_info": {
      "cpu": "Intel i7-8700",
      "memory": "16GB",
      "disk": "512GB SSD",
      "os": "Windows 10 Pro"
    },
    "activation_ip": "192.168.1.100",
    "status": "active",
    "status_display": "激活",
    "is_online": true,
    "is_online_display": "在线",
    "activated_at": "2024-01-01T10:00:00Z",
    "last_heartbeat": "2024-01-01T14:30:00Z",
    "last_online_ip": "192.168.1.100",
    "config_updated_at": "2024-01-01T10:00:00Z",
    "usage_data": {
      "active_users": 50,
      "api_calls_today": 5000
    },
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T14:30:00Z"
  }
}
```

### 2.3 手动添加许可证
```http
POST /api/v1/licenses
```

**请求体**
```json
{
  "authorization_code_id": "auth-code-uuid",
  "hardware_fingerprint": "CPU:ABC123,MB:DEF456,MAC:00:11:22:33:44:55",
  "device_info": {
    "cpu": "Intel i7-8700",
    "memory": "16GB",
    "disk": "512GB SSD",
    "os": "Windows 10 Pro"
  },
  "activation_ip": "192.168.1.100"
}
```

**响应**
```json
{
  "code": "000000",
  "message": "成功",
  "data": {
    "id": "license-uuid",
    "license_key": "LIC-DEVICE-ABC123456789",
    "authorization_code_id": "auth-code-uuid",
    "customer_id": "customer-uuid",
    "hardware_fingerprint": "CPU:ABC123,MB:DEF456,MAC:00:11:22:33:44:55",
    "device_info": {
      "cpu": "Intel i7-8700",
      "memory": "16GB",
      "disk": "512GB SSD",
      "os": "Windows 10 Pro"
    },
    "activation_ip": "192.168.1.100",
    "status": "inactive",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

### 2.4 撤销许可证
```http
PUT /api/v1/licenses/{id}/revoke
```

**请求体**
```json
{
  "reason": "设备更换"
}
```

### 2.5 下载许可证文件
```http
GET /api/v1/licenses/{id}/download
```

**响应**: 返回加密的许可证文件

## 3. 客户端激活 API

### 3.1 软件激活
```http
POST /api/v1/activate
```

**请求体**
```json
{
  "authorization_code": "LIC-COMP001-A7B9X2-C8F4",
  "hardware_fingerprint": "CPU:ABC123,MB:DEF456,MAC:00:11:22:33:44:55",
  "device_info": {
    "cpu": "Intel i7-8700",
    "memory": "16GB",
    "os": "Windows 10 Pro"
  },
  "software_version": "1.0.0"
}
```

**响应**
```json
{
  "code": "000000",
  "message": "成功",
  "data": {
    "license_key": "LIC-DEVICE-ABC123456789",
    "license_file": "base64编码的加密许可证文件",
    "heartbeat_interval": 300
  }
}
```

### 3.2 心跳检测
```http
POST /api/v1/heartbeat
```

**请求体**
```json
{
  "license_key": "LIC-DEVICE-ABC123456789",
  "hardware_fingerprint": "CPU:ABC123,MB:DEF456,MAC:00:11:22:33:44:55",
  "config_updated_at": "2024-01-01T10:00:00Z",
  "usage_data": {
    "active_users": 50,
    "api_calls_today": 5000
  },
  "software_version": "1.0.0"
}
```

**响应**
```json
{
  "code": "000000",
  "message": "成功",
  "data": {
    "status": "active",
    "config_updated": true,
    "license_file": "base64编码的新许可证文件",
    "heartbeat_interval": 300
  }
}
```

## 4. 统计报表 API

### 4.1 授权概览统计
```http
GET /api/v1/stats/overview
```

**响应**
```json
{
  "code": "000000",
  "message": "成功",
  "data": {
    "total_auth_codes": 1234,
    "active_licenses": 856,
    "expiring_soon": 23,
    "abnormal_alerts": 5,
    "growth_rate": {
      "auth_codes": 5.2,
      "licenses": 3.1
    }
  }
}
```

### 4.2 授权变更历史
```http
GET /api/v1/authorization-codes/{id}/changes
```

**查询参数**
- `page` - 页码（默认1）
- `page_size` - 每页数量（默认20，最大100）
- `change_type` - 变更类型筛选 (renewal/upgrade/limit_change/feature_toggle/lock/unlock/other)
- `operator_id` - 操作人筛选
- `start_date` - 开始时间 (YYYY-MM-DD格式)
- `end_date` - 结束时间 (YYYY-MM-DD格式)
- `sort` - 排序字段 (created_at/change_type)
- `order` - 排序方向 (asc/desc)

**响应**
```json
{
  "code": "000000",
  "message": "成功",
  "data": {
    "list": [
      {
        "id": "change-uuid",
        "change_type": "upgrade",
        "change_type_display": "升级",
        "operator_id": "user-uuid",
        "operator_name": "张三",
        "reason": "客户升级套餐",
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 25,
    "page": 1,
    "page_size": 20,
    "total_pages": 2
  }
}
```

## 5. 错误码定义

- `300001` - 授权码不存在
- `300002` - 授权码已存在
- `300003` - 授权码已被锁定
- `300004` - 激活数量已达上限
- `300005` - 硬件指纹格式错误
- `300006` - 许可证不存在
- `300007` - 许可证已被撤销
- `300008` - 硬件指纹不匹配
- `300009` - 许可证文件生成失败
- `300010` - 配置参数错误

## 6. 状态说明

### 6.1 授权码状态
- `normal` - 正常：当前时间在有效期内且未被锁定
- `locked` - 已锁定：is_locked=true
- `expired` - 已过期：当前时间超过end_date

### 6.2 部署类型
- `standalone` - 单机版：独立部署的软件
- `cloud` - 云端版：基于云的SaaS服务
- `hybrid` - 混合版：结合本地和云端的部署

### 6.3 加密类型  
- `standard` - 标准加密：基础安全级别
- `advanced` - 高级加密：增强安全级别

### 6.4 许可证状态
- `active` - 激活：正常使用中
- `inactive` - 未激活：已创建但尚未激活
- `revoked` - 已撤销：被管理员撤销或因违规被停用