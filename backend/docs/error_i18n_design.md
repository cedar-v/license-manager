# 错误信息多语言设计方案

## 设计理念

基于系统设计理念，构建简洁、易用、可维护的多语言错误信息响应系统：

- **配置驱动**：语言包通过配置文件管理，支持热更新
- **代码即文档**：通过结构化的语言文件体现错误信息组织
- **渐进式复杂度**：从单语言平滑升级到多语言支持
- **失败快速原则**：语言缺失时优雅降级到默认语言

## 整体架构设计

### 1. 目录结构

```
backend/
├── configs/
│   └── i18n/                    # 多语言配置目录
│       ├── errors/              # 错误信息语言包
│       │   ├── zh-CN.yaml      # 中文错误信息
│       │   ├── en-US.yaml      # 英文错误信息
│       │   └── ja-JP.yaml      # 日文错误信息
│       └── messages/            # 其他信息语言包（可选扩展）
├── pkg/
│   └── i18n/                   # 国际化工具包
│       ├── manager.go          # 多语言管理器
│       └── errors.go           # 错误信息国际化封装
```

### 2. 核心设计原则

- **单一职责**：i18n 包专门处理多语言逻辑
- **接口解耦**：通过接口与现有错误系统集成
- **配置外化**：所有语言文件独立管理
- **缓存优化**：语言包加载后内存缓存

## 语言包配置格式

### YAML 配置结构

错误信息按模块组织，保持与现有错误码体系一致：

```yaml
# configs/i18n/errors/zh-CN.yaml
system:
  name: "中文"
  locale: "zh-CN"

errors:
  # 成功状态
  "000000": "成功"
  
  # 认证模块 (10xxxx)
  auth:
    "100001": "认证已过期"
    "100002": "认证无效"
    "100003": "用户名或密码错误"
    "100004": "缺少认证信息"
    "100005": "权限不足"
  
  # 客户模块 (20xxxx) - 预留
  customer:
    "200001": "客户不存在"
    "200002": "客户已存在"
  
  # 授权模块 (30xxxx) - 预留
  license:
    "300001": "授权码无效"
    "300002": "授权已过期"
  
  # 系统通用 (90xxxx)
  system:
    "900001": "请求参数无效"
    "900002": "资源不存在"
    "900003": "资源冲突"
    "900004": "服务器内部错误"

# 默认错误信息
default_error: "未知错误"
```

```yaml
# configs/i18n/errors/en-US.yaml
system:
  name: "English"
  locale: "en-US"

errors:
  "000000": "Success"
  
  auth:
    "100001": "Authentication expired"
    "100002": "Invalid authentication"
    "100003": "Invalid username or password"
    "100004": "Authentication required"
    "100005": "Insufficient permissions"
  
  customer:
    "200001": "Customer not found"
    "200002": "Customer already exists"
  
  license:
    "300001": "Invalid license key"
    "300002": "License expired"
  
  system:
    "900001": "Invalid request parameters"
    "900002": "Resource not found"
    "900003": "Resource conflict"
    "900004": "Internal server error"

default_error: "Unknown error"
```

## 核心组件设计

### 1. 多语言管理器

```go
// pkg/i18n/manager.go
type Manager struct {
    defaultLang   string
    loadedLangs   map[string]*LanguageData
    configPath    string
}

type LanguageData struct {
    System       SystemInfo         `yaml:"system"`
    Errors       map[string]string  `yaml:"errors"`
    DefaultError string            `yaml:"default_error"`
}

type SystemInfo struct {
    Name   string `yaml:"name"`
    Locale string `yaml:"locale"`
}

// 核心接口
type I18nManager interface {
    LoadLanguage(lang string) error
    GetErrorMessage(code, lang string) string
    SupportedLanguages() []string
    SetDefaultLanguage(lang string)
}
```

### 2. 错误信息国际化封装

```go
// pkg/i18n/errors.go
type ErrorI18n struct {
    manager I18nManager
}

// 扩展现有错误响应函数
func NewI18nErrorResponse(code, lang string, customMessage ...string) (int, string, string) {
    // 获取本地化错误信息
    message := GetI18nErrorMessage(code, lang)
    
    // 支持自定义消息覆盖
    if len(customMessage) > 0 && customMessage[0] != "" {
        message = customMessage[0]
    }
    
    httpStatus := getHTTPStatusByCode(code)
    return httpStatus, code, message
}
```

## 语言检测策略

### 1. 多层级语言检测

按优先级顺序检测客户端语言偏好：

1. **URL 参数**：`?lang=en-US` (最高优先级)
2. **HTTP Header**：`Accept-Language: zh-CN,zh;q=0.9,en;q=0.8`
3. **用户配置**：数据库存储的用户语言偏好
4. **系统默认**：配置文件指定的默认语言

### 2. 语言检测中间件

```go
// internal/api/middleware/i18n.go
func I18nMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        lang := detectLanguage(c)
        c.Set("language", lang)
        c.Next()
    }
}

func detectLanguage(c *gin.Context) string {
    // 1. URL 参数检测
    if lang := c.Query("lang"); lang != "" {
        return lang
    }
    
    // 2. Accept-Language 头检测
    if lang := parseAcceptLanguage(c.GetHeader("Accept-Language")); lang != "" {
        return lang
    }
    
    // 3. 返回默认语言
    return config.DefaultLanguage
}
```

## 集成实现方案

### 1. 配置系统集成

扩展现有配置结构，添加国际化配置：

```go
// internal/config/config.go
type Config struct {
    // ... 现有配置
    
    I18n I18nConfig `yaml:"i18n"`
}

type I18nConfig struct {
    Enable       bool     `yaml:"enable"`        // 是否启用多语言
    DefaultLang  string   `yaml:"default_lang"`  // 默认语言
    ConfigPath   string   `yaml:"config_path"`   // 语言包路径
    SupportLangs []string `yaml:"support_langs"` // 支持的语言列表
}
```

### 2. Handler 层集成

修改现有错误处理流程，支持多语言响应：

```go
// 使用示例
func (h *AuthHandler) Login(c *gin.Context) {
    // ... 业务逻辑
    
    if err != nil {
        // 获取客户端语言
        lang := c.GetString("language")
        
        // 返回本地化错误响应
        httpStatus, code, message := i18n.NewI18nErrorResponse("100003", lang)
        c.JSON(httpStatus, models.ErrorResponse{
            Code:      httpStatus,
            Error:     code,
            Message:   message,
            Timestamp: time.Now().Format(time.RFC3339),
        })
        return
    }
}
```

## 性能优化策略

### 1. 缓存机制

- **启动加载**：应用启动时预加载常用语言包
- **内存缓存**：语言包加载后常驻内存
- **懒加载**：按需加载非常用语言包

### 2. 降级策略

- **语言缺失降级**：请求语言不存在时降级到默认语言
- **错误码缺失降级**：特定错误码翻译缺失时使用通用错误信息
- **完全降级**：多语言系统异常时回退到原有单语言模式

## 开发和维护流程

### 1. 语言包维护

- **版本控制**：语言文件纳入 Git 管理
- **格式校验**：CI 流程中校验 YAML 格式正确性
- **完整性检查**：确保所有语言包包含相同的错误码

### 2. 新增错误码流程

1. 在 `pkg/errors/codes.go` 中添加新错误码
2. 在所有语言包中添加对应翻译
3. 更新相关文档和测试用例

### 3. 新增语言支持流程

1. 创建新的语言包文件（如 `fr-FR.yaml`）
2. 翻译所有错误信息
3. 在配置中添加新语言到支持列表
4. 更新 API 文档

## 实施建议

### 第一阶段：基础架构

1. 实现多语言管理器核心功能
2. 创建中英文语言包
3. 集成语言检测中间件
4. 改造现有错误处理函数

### 第二阶段：功能完善

1. 添加更多语言支持
2. 实现配置热更新
3. 添加性能监控
4. 完善错误处理和降级机制

### 第三阶段：扩展应用

1. 扩展到其他响应信息的多语言支持
2. 添加用户语言偏好存储
3. 实现动态语言切换
4. 提供语言包管理接口

## 配置示例

### 应用配置

```yaml
# configs/config.yaml
i18n:
  enable: true
  default_lang: "zh-CN"
  config_path: "configs/i18n/errors"
  support_langs: ["zh-CN", "en-US", "ja-JP"]
```

### 使用示例

```go
// 初始化
i18nManager := i18n.NewManager("configs/i18n/errors", "zh-CN")
i18nManager.LoadLanguage("zh-CN")
i18nManager.LoadLanguage("en-US")

// 使用
message := i18nManager.GetErrorMessage("100003", "en-US")
// 输出：Invalid username or password
```

这个设计方案遵循系统的设计理念，提供了简洁、可维护的多语言错误信息支持，同时保持了与现有架构的良好集成。