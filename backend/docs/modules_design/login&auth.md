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

