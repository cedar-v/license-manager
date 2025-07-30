# API开发规范

## 目录结构规范

### 分层架构
```
internal/
├── models/          # 数据模型层
├── service/         # 业务逻辑层  
├── repository/      # 数据访问层
└── api/handlers/    # HTTP处理层
```

### 职责划分
- **Models**: 定义请求/响应结构体，数据验证规则
- **Service**: 实现业务逻辑，处理业务规则
- **Repository**: 数据库操作，缓存操作
- **Handler**: HTTP请求处理，参数绑定，响应格式化

## 文件命名规范

### 模块分组
按业务模块组织文件：
```
models/
├── auth.go         # 认证相关模型
├── customer.go     # 客户相关模型
├── license.go      # 授权相关模型
└── common.go       # 通用模型

service/
├── interfaces.go   # 服务接口定义
├── auth_service.go
├── customer_service.go
└── license_service.go

api/handlers/
├── auth_handler.go
├── customer_handler.go
└── license_handler.go
```

## 代码结构规范

### 1. 模型定义 (models/)
```go
// 请求模型
type CreateCustomerRequest struct {
    Name    string `json:"name" binding:"required"`
    Email   string `json:"email" binding:"required,email"`
    Company string `json:"company"`
}

// 响应模型
type CustomerResponse struct {
    ID      int    `json:"id"`
    Name    string `json:"name"`
    Email   string `json:"email"`
    Company string `json:"company"`
}
```

### 2. 服务接口 (service/interfaces.go)
```go
type CustomerService interface {
    CreateCustomer(req *models.CreateCustomerRequest) (*models.CustomerResponse, error)
    GetCustomer(id int) (*models.CustomerResponse, error)
    UpdateCustomer(id int, req *models.UpdateCustomerRequest) error
    DeleteCustomer(id int) error
}
```

### 3. 服务实现 (service/)
```go
type customerService struct {
    repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
    return &customerService{repo: repo}
}

func (s *customerService) CreateCustomer(req *models.CreateCustomerRequest) (*models.CustomerResponse, error) {
    // 业务逻辑处理
    // 数据验证
    // 调用repository
}
```

### 4. 处理器实现 (api/handlers/)
```go
type CustomerHandler struct {
    service service.CustomerService
}

func NewCustomerHandler(service service.CustomerService) *CustomerHandler {
    return &CustomerHandler{service: service}
}

// @Summary 创建客户
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param request body models.CreateCustomerRequest true "客户信息"
// @Success 200 {object} models.APIResponse
// @Router /api/v1/customers [post]
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
    var req models.CreateCustomerRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Code: http.StatusBadRequest,
            Error: "INVALID_REQUEST",
            Message: "请求参数无效",
        })
        return
    }

    result, err := h.service.CreateCustomer(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Code: http.StatusInternalServerError,
            Error: "CREATE_FAILED",
            Message: err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, models.APIResponse{
        Code: http.StatusOK,
        Message: "创建成功",
        Data: result,
    })
}
```

## API设计规范

### 1. 路由设计
```
GET    /api/v1/customers       # 获取客户列表
POST   /api/v1/customers       # 创建客户
GET    /api/v1/customers/:id   # 获取客户详情
PUT    /api/v1/customers/:id   # 更新客户
DELETE /api/v1/customers/:id   # 删除客户
```

### 2. 响应格式
```go
// 成功响应
type APIResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// 错误响应
type ErrorResponse struct {
    Code      int    `json:"code"`
    Error     string `json:"error"`
    Message   string `json:"message"`
    Timestamp string `json:"timestamp"`
}
```

### 3. HTTP状态码
- `200` - 请求成功
- `201` - 创建成功
- `400` - 请求参数错误
- `401` - 未认证
- `403` - 权限不足
- `404` - 资源不存在
- `500` - 服务器内部错误

## Swagger文档规范

### 必需注解
```go
// @Summary 接口简要描述
// @Description 接口详细描述
// @Tags 标签分组
// @Accept json
// @Produce json
// @Param 参数名 参数位置 参数类型 是否必需 "参数描述"
// @Success 状态码 {object} 响应类型 "成功描述"
// @Failure 状态码 {object} 错误类型 "失败描述"
// @Security BearerAuth (需要认证的接口)
// @Router 路由路径 [HTTP方法]
```

### 标签分组
- `认证` - 登录、登出、Token相关
- `客户管理` - 客户CRUD操作
- `授权管理` - 授权码、许可证相关
- `系统` - 系统信息、健康检测

## 错误处理规范

### 错误码定义
```go
const (
    // 通用错误
    ErrInvalidRequest = "INVALID_REQUEST"
    ErrUnauthorized   = "UNAUTHORIZED"
    ErrForbidden      = "FORBIDDEN"
    ErrNotFound       = "NOT_FOUND"
    ErrInternalError  = "INTERNAL_ERROR"
    
    // 业务错误
    ErrCustomerExists    = "CUSTOMER_EXISTS"
    ErrLicenseExpired    = "LICENSE_EXPIRED"
    ErrInvalidLicense    = "INVALID_LICENSE"
)
```

### 统一错误处理
```go
func handleError(c *gin.Context, err error) {
    var statusCode int
    var errorCode string
    
    switch err {
    case ErrCustomerExists:
        statusCode = http.StatusConflict
        errorCode = "CUSTOMER_EXISTS"
    default:
        statusCode = http.StatusInternalServerError
        errorCode = "INTERNAL_ERROR"
    }
    
    c.JSON(statusCode, models.ErrorResponse{
        Code:      statusCode,
        Error:     errorCode,
        Message:   err.Error(),
        Timestamp: time.Now().Format(time.RFC3339),
    })
}
```

## 开发流程

### 1. 新增API步骤
1. 在 `models/` 定义请求/响应结构体
2. 在 `service/interfaces.go` 添加服务接口
3. 在 `service/` 实现业务逻辑
4. 在 `handlers/` 实现HTTP处理器
5. 在 `routes/router.go` 注册路由
6. 添加Swagger注解
7. 生成文档：`swag init -g cmd/main.go -o ./docs/swagger`

### 2. 代码检查清单
- [ ] 模型定义是否合理
- [ ] 业务逻辑是否在service层
- [ ] 错误处理是否完整
- [ ] Swagger注解是否完整
- [ ] 路由是否正确注册
- [ ] HTTP状态码是否合适

## 最佳实践

1. **职责单一**：每个函数只做一件事
2. **接口优先**：先定义接口，再实现
3. **错误优雅**：统一错误处理和响应格式
4. **文档完整**：所有API都要有Swagger文档
5. **命名规范**：使用有意义的命名
6. **代码复用**：提取公共逻辑到utils包