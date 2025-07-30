# 登录与认证模块设计

## 概述

本模块负责系统的用户认证、登录状态管理和权限控制，基于JWT（JSON Web Token）实现无状态认证机制。

## 技术选型

- **认证方式**: JWT (JSON Web Token)
- **加密算法**: HS256 (HMAC-SHA256)
- **密码加密**: bcrypt
- **Session存储**: 内存缓存/Redis（可选）

## 社区版登录设计

用户名和密码配置在配置文件中，默认用户名和密码为 `admin` / `admin@123`

## 功能模块

### 1. 用户登录 (Login)

#### 1.1 登录流程
```
1. 用户提交用户名/密码
2. 验证用户凭证
3. 生成JWT Token
4. 返回Token和用户信息
5. 前端存储Token（localStorage/sessionStorage）
```

#### 1.2 API接口设计
```http
POST /api/v1/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin@123"
}
```

**响应格式：**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600,
    "user_info": {
      "username": "admin",
      "role": "administrator"
    }
  }
}
```

#### 1.3 实现细节
- 密码使用bcrypt进行哈希验证
- JWT包含用户ID、用户名、角色、过期时间
- Token默认有效期1小时，可配置
- 支持记住我功能（延长Token有效期）

### 2. 用户登出 (Logout)

#### 2.1 登出流程
```
1. 客户端删除本地Token
2. 可选：服务端将Token加入黑名单
3. 返回登出成功响应
```

#### 2.2 API接口设计
```http
POST /api/v1/logout
Authorization: Bearer <token>
```

**响应格式：**
```json
{
  "code": 200,
  "message": "登出成功"
}
```

### 3. Token刷新

#### 3.1 刷新机制
- 当Token即将过期时（剩余时间<30分钟），自动刷新
- 刷新Token与访问Token分离（可选实现）

#### 3.2 API接口设计
```http
POST /api/v1/auth/refresh
Authorization: Bearer <token>
```

## JWT Token设计

### 1. Token结构
```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "user_id": 1,
    "username": "admin",
    "role": "administrator",
    "iat": 1640995200,
    "exp": 1641002400
  },
  "signature": "signature_string"
}
```

### 2. Claims定义
- `user_id`: 用户ID
- `username`: 用户名
- `role`: 用户角色
- `iat`: 签发时间
- `exp`: 过期时间

## 中间件认证设计

### 1. 认证中间件 (AuthMiddleware)

#### 1.1 功能职责
- 验证HTTP请求头中的JWT Token
- 解析Token获取用户信息
- 处理Token过期和无效情况
- 将用户信息注入到请求上下文

#### 1.2 认证流程
```
1. 检查Authorization头是否存在
2. 提取Bearer Token
3. 验证Token签名和有效性
4. 解析用户信息
5. 将用户信息存储到Context
6. 继续处理请求或返回401错误
```

#### 1.3 实现伪代码
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 获取Token
        token := extractToken(c)
        if token == "" {
            c.JSON(401, gin.H{"error": "Missing token"})
            c.Abort()
            return
        }
        
        // 2. 验证Token
        claims, err := validateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        // 3. 设置用户信息到Context
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        
        c.Next()
    }
}
```

### 2. 可选中间件 (OptionalAuth)

对于某些接口可能需要可选认证（如公开API），提供可选认证中间件：

```go
func OptionalAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c)
        if token != "" {
            if claims, err := validateToken(token); err == nil {
                c.Set("user_id", claims.UserID)
                c.Set("username", claims.Username)
                c.Set("role", claims.Role)
            }
        }
        c.Next()
    }
}
```

## 路由保护策略

### 1. 路由分组
```go
// 公开路由
public := r.Group("/api/v1")
{
    public.POST("/login", handlers.Login)
}

// 需要认证的路由
auth := r.Group("/api/v1")
auth.Use(middleware.AuthMiddleware())
{
    auth.POST("/logout", handlers.Logout)
    auth.GET("/customers", handlers.GetCustomers)
    auth.POST("/licenses", handlers.CreateLicense)
}

// 管理员路由
admin := r.Group("/api/v1/admin")
admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
{
    admin.GET("/system/info", handlers.GetSystemInfo)
}
```

### 2. 权限级别
- **Public**: 无需认证
- **Authenticated**: 需要有效Token
- **Admin**: 需要管理员权限

## 安全措施

### 1. Token安全
- 使用强密钥签名（至少32字节）
- 设置合理的过期时间

### 2. 密码安全
- bcrypt加密，cost factor >= 10
- 登录失败限制（防暴力破解）


## 配置参数

```yaml
auth:
  jwt:
    secret: "your-256-bit-secret"
    expire_hours: 1
    refresh_threshold_minutes: 30
  admin:
    username: "admin"
    password: "admin@123" # 明文密码 部署人员可修改
  security:
    max_login_attempts: 5
    lockout_duration_minutes: 15
```

## 错误处理

### 1. 认证错误码
- `AUTH_001`: Token缺失
- `AUTH_002`: Token无效
- `AUTH_003`: Token过期
- `AUTH_004`: 权限不足
- `AUTH_005`: 用户名或密码错误
- `AUTH_006`: 账户被锁定

### 2. 统一错误响应
```json
{
  "code": 401,
  "error": "AUTH_002",
  "message": "Token无效",
  "timestamp": "2024-01-01T12:00:00Z"
}
```
