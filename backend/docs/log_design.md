# 日志系统设计方案

## 设计理念

基于系统设计理念，构建轻量级、高性能、可观测的日志系统：

- **结构化日志**：使用 Logrus 提供结构化日志输出，便于日志分析和监控
- **分层记录**：HTTP 请求日志、业务逻辑日志、错误日志分层记录
- **配置驱动**：通过配置文件控制日志级别和格式
- **性能友好**：日志输出异步化，不影响业务处理性能
- **可观测性**：支持日志聚合和监控集成

## 整体架构设计

### 1. 目录结构

```
backend/
├── configs/                    # 配置目录
│   └── config.yaml            # 日志配置
├── pkg/logger/                # 日志封装包
│   └── logger.go              # 日志管理器和自定义格式器
├── internal/
│   ├── api/
│   │   └── middleware/
│   │       └── logging.go     # HTTP请求日志中间件
│   └── service/               # 业务服务层（日志输出点）
└── docs/
    └── log_design.md          # 日志设计文档
```

### 2. 核心设计原则

- **单一职责**：日志包专门处理日志格式化和输出逻辑
- **配置解耦**：日志配置独立管理，支持运行时调整
- **性能优先**：日志输出异步，不阻塞业务逻辑
- **扩展性**：支持多种输出目标（文件、Elasticsearch等）

## 核心组件设计

### 1. 日志管理器

```go
// pkg/logger/logger.go
type CustomFormatter struct {
    *logrus.TextFormatter
}

// Init() 初始化日志系统
// GetLogger() 获取全局日志实例
```

**核心特性：**
- 基于 Logrus 的封装
- 自定义格式器，支持调用者信息显示
- 配置文件驱动的日志级别设置

### 2. HTTP 请求日志中间件

```go
// internal/api/middleware/logging.go
func CustomLoggerMiddleware() gin.HandlerFunc
func LoggingMiddleware(logger *logrus.Logger) gin.HandlerFunc
```

**记录内容：**
- 客户端IP
- 请求方法和路径
- 响应状态码
- 请求延迟时间
- User-Agent信息

### 3. 业务日志输出

```go
// 业务服务中的日志输出示例
s.logger.Infof("[GenerateLicenseFile] 开始生成许可证文件，license_id: %s", id)
s.logger.Errorf("[GenerateLicenseFile] 查询许可证失败，license_id: %s, error: %v", id, err)
```

**日志分类：**
- **Info级别**：正常业务流程记录
- **Warn级别**：异常情况但不影响业务
- **Error级别**：业务错误和系统异常

## 日志格式和输出

### 1. 日志格式

**标准格式：**
```
2024-12-25 14:30:15 INFO API Request status=200 method=GET path=/api/customers ip=127.0.0.1 latency=15.2ms
```

**带调用者信息：**
```
2024-12-25 14:30:15 INFO [GenerateLicenseFile] 开始生成许可证文件，license_id: 12345 ./backend/internal/service/license_service.go:290
```

### 2. 日志级别

| 级别 | 描述 | 使用场景 |
|------|------|----------|
| DEBUG | 调试信息 | 开发环境，详细的执行流程 |
| INFO | 一般信息 | 正常业务流程，重要状态变更 |
| WARN | 警告信息 | 异常情况但不影响业务继续 |
| ERROR | 错误信息 | 业务错误，系统异常 |

### 3. 结构化字段

**HTTP请求日志字段：**
- `status`: HTTP状态码
- `method`: 请求方法
- `path`: 请求路径
- `ip`: 客户端IP
- `latency`: 请求处理时间
- `user_agent`: 用户代理

**业务日志字段：**
- `[ModuleName]`: 模块标识
- 业务相关参数（如 license_id, customer_id 等）

## 配置系统集成

### 1. 日志配置结构

```yaml
# configs/config.yaml
log:
  level: "info"     # 日志级别: debug/info/warn/error
  format: "text"    # 日志格式: text/json
```

### 2. 配置加载逻辑

```go
// internal/config/config.go
type LogConfig struct {
    Level  string `mapstructure:"level"`
    Format string `mapstructure:"format"`
}

// 默认配置
viper.SetDefault("log.level", "info")
viper.SetDefault("log.format", "json")
```

## 集成实现方案

### 1. 应用启动集成

```go
// cmd/main.go
func main() {
    // 初始化日志系统
    logger.Init()
    log := logger.GetLogger()
    
    log.Info("启动 License Manager 服务器...")
}
```

### 2. Gin 路由集成

```go
// internal/api/routes/router.go
func SetupRouter() *gin.Engine {
    router := gin.New()
    
    // 添加日志中间件
    router.Use(middleware.CustomLoggerMiddleware())
    
    return router
}
```

### 3. 服务层集成

```go
// internal/service/license_service.go
type LicenseService struct {
    logger *logrus.Logger
}

func NewLicenseService(repo repository.LicenseRepository, logger *logrus.Logger) LicenseService {
    return LicenseService{
        logger: logger,
    }
}
```

## 日志输出示例

### 1. HTTP 请求日志

```
2024-12-25 14:30:15 INFO API Request status=200 method=GET path=/api/customers ip=127.0.0.1 latency=15.2ms user_agent=Mozilla/5.0...
```

### 2. 业务操作日志

```
2024-12-25 14:30:16 INFO [GenerateLicenseFile] 开始生成许可证文件，license_id: 12345 ./backend/internal/service/license_service.go:290
2024-12-25 14:30:16 INFO [GenerateLicenseFile] 许可证信息查询成功，license_key: LK-12345, status: active ./backend/internal/service/license_service.go:309
2024-12-25 14:30:16 WARN [GenerateLicenseFile] 许可证已被撤销，license_id: 12345 ./backend/internal/service/license_service.go:313
```

### 3. 错误日志

```
2024-12-25 14:30:17 ERROR [GenerateLicenseFile] 查询许可证失败，license_id: 12345, error: license not found ./backend/internal/service/license_service.go:302
```

## 性能优化策略

### 1. 异步日志输出

- 日志输出默认异步，不阻塞业务逻辑
- 使用缓冲区批量写入，提升I/O性能

### 2. 日志级别控制

- 生产环境使用 INFO 级别及以上
- 开发环境可启用 DEBUG 级别
- 根据需要动态调整日志级别

### 3. 结构化日志优化

- 使用预定义的字段名称
- 避免大对象序列化
- 敏感信息脱敏处理

## 前端日志处理

### 1. 开发环境日志

```typescript
// 前端调试日志
console.log('[Licenses] openChangeValidityDialog called')
console.log('[Licenses] changeDialogVisible ->', val)
```

### 2. 生产环境优化

- 开发环境的 console.log 在生产环境自动移除
- 错误日志通过统一的错误处理中间件记录
- 用户行为日志通过专门的埋点系统处理

## 监控和告警集成

### 1. 日志聚合

- 使用 ELK Stack 或 Loki 进行日志聚合
- 支持日志查询和可视化分析

### 2. 告警规则

- ERROR 级别日志自动告警
- 关键业务异常监控
- 性能指标监控（如请求延迟）

### 3. 审计日志

- 敏感操作记录（如用户登录、授权变更）
- 支持审计追溯和合规要求

## 实施策略

### 第一阶段：基础架构

1. 实现日志管理器核心功能
2. 配置基础日志格式和级别
3. 集成 HTTP 请求日志中间件
4. 在关键业务节点添加日志输出

### 第二阶段：功能完善

1. 优化日志格式和字段
2. 添加业务日志分类
3. 实现日志轮转和归档
4. 集成日志监控系统

### 第三阶段：运维优化

1. 实现分布式日志追踪
2. 添加性能监控指标
3. 建立日志分析和告警机制
4. 完善日志安全和隐私保护

## 配置示例

### 开发环境配置

```yaml
log:
  level: "debug"    # 开发环境启用debug级别
  format: "text"    # 文本格式便于阅读
```

### 生产环境配置

```yaml
log:
  level: "info"     # 生产环境使用info级别
  format: "json"    # JSON格式便于日志聚合
```

这个日志设计方案遵循系统的设计理念，提供了结构化、可观测、性能优异的日志系统，为系统的监控、调试和运维提供了有力支持。

