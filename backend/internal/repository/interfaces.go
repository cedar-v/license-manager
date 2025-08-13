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