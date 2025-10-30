package repository

import (
	"context"
	"license-manager/internal/models"
	"time"
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

	// UpdateAuthorizationCode 更新授权码
	UpdateAuthorizationCode(ctx context.Context, authCode *models.AuthorizationCode) error

	// DeleteAuthorizationCode 删除授权码（软删除）
	DeleteAuthorizationCode(ctx context.Context, id string) error

	// CheckCustomerExists 检查客户是否存在
	CheckCustomerExists(ctx context.Context, customerID string) (bool, error)

	// GetAuthorizationChangeList 查询授权变更历史列表
	GetAuthorizationChangeList(ctx context.Context, authCodeID string, req *models.AuthorizationChangeListRequest) (*models.AuthorizationChangeListResponse, error)

	// RecordAuthorizationChange 记录授权变更历史
	RecordAuthorizationChange(ctx context.Context, change *models.AuthorizationChange) error
}

// LicenseRepository 许可证数据访问接口
type LicenseRepository interface {
	// GetLicenseList 查询许可证列表
	GetLicenseList(ctx context.Context, req *models.LicenseListRequest) (*models.LicenseListResponse, error)

	// GetLicenseByID 根据ID获取许可证信息
	GetLicenseByID(ctx context.Context, id string) (*models.License, error)

	// CreateLicense 创建许可证
	CreateLicense(ctx context.Context, license *models.License) error

	// UpdateLicense 更新许可证信息
	UpdateLicense(ctx context.Context, license *models.License) error

	// CheckAuthorizationCodeExists 检查授权码是否存在
	CheckAuthorizationCodeExists(ctx context.Context, authCodeID string) (bool, error)

	// GetAuthorizationCodeByID 根据ID获取授权码信息
	GetAuthorizationCodeByID(ctx context.Context, authCodeID string) (*models.AuthorizationCode, error)

	// GetAuthorizationCodeByCode 根据授权码获取授权码信息
	GetAuthorizationCodeByCode(ctx context.Context, code string) (*models.AuthorizationCode, error)

	// GetLicenseByKey 根据许可证密钥获取许可证信息
	GetLicenseByKey(ctx context.Context, licenseKey string) (*models.License, error)

	// GetActiveLicenseCount 获取指定授权码的激活许可证数量
	GetActiveLicenseCount(ctx context.Context, authCodeID string) (int64, error)
}

// DashboardRepository 仪表盘数据访问接口
type DashboardRepository interface {
	// GetAuthorizationTrendData 获取授权趋势数据
	GetAuthorizationTrendData(ctx context.Context, startDate, endDate time.Time) ([]models.TrendData, error)

	// GetRecentAuthorizations 获取最近授权列表
	GetRecentAuthorizations(ctx context.Context, req *models.DashboardRecentAuthorizationsRequest) (*models.DashboardRecentAuthorizationsResponse, error)
}
