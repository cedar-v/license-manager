package repository

import (
	"context"
	"license-manager/internal/models"
)

// CustomerRepository 客户数据访问接口
type CustomerRepository interface {
	// GetCustomerList 查询客户列表
	GetCustomerList(ctx context.Context, req *models.CustomerListRequest) (*models.CustomerListResponse, error)
	
	// GetCustomerByID 根据ID获取客户信息
	GetCustomerByID(ctx context.Context, id string) (*models.Customer, error)
	
	// CreateCustomer 创建客户
	CreateCustomer(ctx context.Context, customer *models.Customer) error
	
	// UpdateCustomer 更新客户信息
	UpdateCustomer(ctx context.Context, customer *models.Customer) error
	
	// DeleteCustomer 删除客户
	DeleteCustomer(ctx context.Context, id string) error
	
	// GetCustomerCount 获取客户总数（用于统计）
	GetCustomerCount(ctx context.Context, filters map[string]interface{}) (int64, error)
}

// UserRepository 用户数据访问接口
type UserRepository interface {
	// GetUserByUsername 根据用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	
	// GetUserByID 根据ID获取用户
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	
	// UpdateUser 更新用户信息
	UpdateUser(ctx context.Context, user *models.User) error
	
	// IncrementLoginAttempts 增加登录失败次数
	IncrementLoginAttempts(ctx context.Context, id string) error
	
	// ResetLoginAttempts 重置登录失败次数
	ResetLoginAttempts(ctx context.Context, id string) error
	
	// LockUser 锁定用户账号
	LockUser(ctx context.Context, id string, lockDuration int) error
}

// AuthorizationCodeRepository 授权码数据访问接口
type AuthorizationCodeRepository interface {
	// CreateAuthorizationCode 创建授权码
	CreateAuthorizationCode(ctx context.Context, authCode *models.AuthorizationCode) error
	
	// GetAuthorizationCodeByID 根据ID获取授权码
	GetAuthorizationCodeByID(ctx context.Context, id string) (*models.AuthorizationCode, error)
	
	// GetAuthorizationCodeList 查询授权码列表
	GetAuthorizationCodeList(ctx context.Context, req *models.AuthorizationCodeListRequest) (*models.AuthorizationCodeListResponse, error)
	
	// CheckCustomerExists 检查客户是否存在
	CheckCustomerExists(ctx context.Context, customerID string) (bool, error)
}