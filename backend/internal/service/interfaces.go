package service

import (
	"context"
	"license-manager/internal/models"
)

// AuthService 认证服务接口
type AuthService interface {
	Login(ctx context.Context, req *models.LoginRequest, clientIP string) (*models.LoginData, error)
	RefreshToken(token string) (string, error)
	ValidateToken(token string) error
}

// SystemService 系统服务接口
type SystemService interface {
	GetHealthStatus() *models.HealthResponse
	GetSystemInfo() map[string]interface{}
}

// CustomerService 客户服务接口
type CustomerService interface {
	GetCustomerList(ctx context.Context, req *models.CustomerListRequest) (*models.CustomerListResponse, error)
	GetCustomer(ctx context.Context, id string) (*models.Customer, error)
	CreateCustomer(ctx context.Context, req *models.CustomerCreateRequest) (*models.Customer, error)
	UpdateCustomer(ctx context.Context, id string, req *models.CustomerUpdateRequest) (*models.Customer, error)
	DeleteCustomer(ctx context.Context, id string) error
	UpdateCustomerStatus(ctx context.Context, id string, req *models.CustomerStatusUpdateRequest) (*models.Customer, error)
}

// EnumService 枚举服务接口
type EnumService interface {
	GetAllEnums(ctx context.Context) (*models.EnumListResponse, error)
	GetEnumsByType(ctx context.Context, enumType string) (*models.EnumTypeResponse, error)
}
