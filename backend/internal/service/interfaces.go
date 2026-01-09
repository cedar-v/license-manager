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

// AuthorizationCodeService 授权码服务接口
type AuthorizationCodeService interface {
	CreateAuthorizationCode(ctx context.Context, req *models.AuthorizationCodeCreateRequest) (*models.AuthorizationCodeCreateResponse, error)
	GetAuthorizationCodeList(ctx context.Context, req *models.AuthorizationCodeListRequest) (*models.AuthorizationCodeListResponse, error)
	GetAuthorizationCode(ctx context.Context, id string) (*models.AuthorizationCode, error)
	UpdateAuthorizationCode(ctx context.Context, id string, req *models.AuthorizationCodeUpdateRequest) (*models.AuthorizationCode, error)
	LockUnlockAuthorizationCode(ctx context.Context, id string, req *models.AuthorizationCodeLockRequest) (*models.AuthorizationCode, error)
	DeleteAuthorizationCode(ctx context.Context, id string) error
	GetAuthorizationChangeList(ctx context.Context, authCodeID string, req *models.AuthorizationChangeListRequest) (*models.AuthorizationChangeListResponse, error)
	GenerateAuthorizationFile(ctx context.Context, id string) ([]byte, string, string, error)

	// 用户端分享功能
	ShareAuthorizationCode(ctx context.Context, authCodeID, userID string, req *models.AuthorizationCodeShareRequest) (*models.AuthorizationCodeShareResponse, error)
}

// LicenseService 许可证服务接口
type LicenseService interface {
	GetLicenseList(ctx context.Context, req *models.LicenseListRequest) (*models.LicenseListResponse, error)
	GetLicense(ctx context.Context, id string) (*models.LicenseDetailResponse, error)
	CreateLicense(ctx context.Context, req *models.LicenseCreateRequest) (*models.License, error)
	RevokeLicense(ctx context.Context, id string, req *models.LicenseRevokeRequest) (*models.License, error)
	GenerateLicenseFile(ctx context.Context, id string) ([]byte, string, string, error)

	// 客户端激活和心跳接口
	ActivateLicense(ctx context.Context, req *models.ActivateRequest, clientIP string) (*models.ActivateResponse, error)
	Heartbeat(ctx context.Context, req *models.HeartbeatRequest, clientIP string) (*models.HeartbeatResponse, error)

	// 统计接口
	GetStatsOverview(ctx context.Context) (*models.StatsOverviewResponse, error)
}

// DashboardService 仪表盘服务接口
type DashboardService interface {
	// 获取授权趋势数据
	GetAuthorizationTrend(ctx context.Context, req *models.DashboardAuthorizationTrendRequest) (*models.DashboardAuthorizationTrendResponse, error)

	// 获取最近授权列表
	GetRecentAuthorizations(ctx context.Context, req *models.DashboardRecentAuthorizationsRequest) (*models.DashboardRecentAuthorizationsResponse, error)
}

// CuDeviceService 客户设备服务接口
type CuDeviceService interface {
	// 获取设备列表
	GetDeviceList(ctx context.Context, customerID string, req *models.DeviceListRequest) (*models.DeviceListResponse, error)

	// 解绑设备
	UnbindDevice(ctx context.Context, customerID, licenseID string) error
}

