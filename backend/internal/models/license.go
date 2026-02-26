package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// License 许可证模型
type License struct {
	ID                  string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	LicenseKey          string         `gorm:"type:varchar(200);uniqueIndex;not null" json:"license_key"`
	AuthorizationCodeID string         `gorm:"type:varchar(36);not null;index" json:"authorization_code_id"`
	CustomerID          string         `gorm:"type:varchar(36);not null;index" json:"customer_id"`
	HardwareFingerprint string         `gorm:"type:varchar(200);not null;index" json:"hardware_fingerprint"`
	DeviceInfo          JSON           `gorm:"type:json" json:"device_info,omitempty" swaggertype:"object"`
	ActivationIP        *string        `gorm:"type:varchar(45)" json:"activation_ip"`
	Status              string         `gorm:"type:varchar(20);not null;default:'inactive';index" json:"status"`
	StatusDisplay       string         `gorm:"-" json:"status_display,omitempty"`
	ActivatedAt         *time.Time     `gorm:"index" json:"activated_at"`
	LastHeartbeat       *time.Time     `gorm:"index" json:"last_heartbeat"`
	LastOnlineIP        *string        `gorm:"type:varchar(45)" json:"last_online_ip"`
	ConfigUpdatedAt     *time.Time     `gorm:"index" json:"config_updated_at"`
	UsageData           JSON           `gorm:"type:json" json:"usage_data,omitempty" swaggertype:"object"`
	CreatedAt           time.Time      `gorm:"not null;index" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联字段（用于查询时的JOIN）
	AuthorizationCode *AuthorizationCode `gorm:"foreignKey:AuthorizationCodeID" json:"authorization_code,omitempty"`
	Customer          *Customer          `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	IsOnline          bool               `gorm:"-" json:"is_online"`
	IsOnlineDisplay   string             `gorm:"-" json:"is_online_display,omitempty"`
}

// TableName 指定表名
func (License) TableName() string {
	return "licenses"
}

// BeforeCreate 创建前自动设置时间戳和ID
func (l *License) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}

	now := time.Now()
	if l.CreatedAt.IsZero() {
		l.CreatedAt = now
	}
	if l.UpdatedAt.IsZero() {
		l.UpdatedAt = now
	}
	return nil
}

// BeforeUpdate 更新前自动刷新更新时间
func (l *License) BeforeUpdate(tx *gorm.DB) error {
	l.UpdatedAt = time.Now()
	return nil
}

// LicenseListRequest 许可证列表查询请求结构
type LicenseListRequest struct {
	Page                int     `form:"page" binding:"omitempty,min=1"`                                                   // 页码，默认1
	PageSize            int     `form:"page_size" binding:"omitempty,min=1,max=100"`                                      // 每页条数，默认20，最大100
	AuthorizationCodeID string  `form:"authorization_code_id" binding:"omitempty"`                                        // 授权码ID筛选
	CustomerID          string  `form:"customer_id" binding:"omitempty"`                                                  // 客户ID筛选
	Status              string  `form:"status" binding:"omitempty,oneof=active inactive revoked"`                         // 状态筛选
	IsOnline            *string `form:"is_online" binding:"omitempty"`                                                    // 在线状态筛选
	Sort                string  `form:"sort" binding:"omitempty,oneof=created_at updated_at activated_at last_heartbeat"` // 排序字段，默认created_at
	Order               string  `form:"order" binding:"omitempty,oneof=asc desc"`                                         // 排序方向，默认desc
}

// LicenseListItem 许可证列表项结构
type LicenseListItem struct {
	ID                  string  `json:"id"`                          // 许可证ID
	LicenseKey          string  `json:"license_key"`                 // 许可证密钥
	AuthorizationCodeID string  `json:"authorization_code_id"`       // 授权码ID
	AuthorizationCode   string  `json:"authorization_code"`          // 授权码
	CustomerName        string  `json:"customer_name"`               // 客户名称
	HardwareFingerprint string  `json:"hardware_fingerprint"`        // 硬件指纹
	Status              string  `json:"status"`                      // 许可证状态
	StatusDisplay       string  `json:"status_display,omitempty"`    // 状态显示名称
	IsOnline            bool    `json:"is_online"`                   // 是否在线
	IsOnlineDisplay     string  `json:"is_online_display,omitempty"` // 在线状态显示名称
	ActivationIP        *string `json:"activation_ip"`               // 激活IP
	LastOnlineIP        *string `json:"last_online_ip"`              // 最后在线IP
	ActivatedAt         *string `json:"activated_at"`                // 激活时间
	LastHeartbeat       *string `json:"last_heartbeat"`              // 最后心跳时间
}

// LicenseListResponse 许可证列表响应结构
type LicenseListResponse struct {
	List       []LicenseListItem `json:"list"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// LicenseDetailResponse 许可证详情响应结构
type LicenseDetailResponse struct {
	ID                  string                 `json:"id"`                          // 许可证ID
	LicenseKey          string                 `json:"license_key"`                 // 许可证密钥
	AuthorizationCodeID string                 `json:"authorization_code_id"`       // 授权码ID
	AuthorizationCode   string                 `json:"authorization_code"`          // 授权码
	CustomerID          string                 `json:"customer_id"`                 // 客户ID
	CustomerName        string                 `json:"customer_name"`               // 客户名称
	HardwareFingerprint string                 `json:"hardware_fingerprint"`        // 硬件指纹
	DeviceInfo          map[string]interface{} `json:"device_info,omitempty"`       // 设备信息
	ActivationIP        *string                `json:"activation_ip"`               // 激活IP
	Status              string                 `json:"status"`                      // 许可证状态
	StatusDisplay       string                 `json:"status_display,omitempty"`    // 状态显示名称
	IsOnline            bool                   `json:"is_online"`                   // 是否在线
	IsOnlineDisplay     string                 `json:"is_online_display,omitempty"` // 在线状态显示名称
	ActivatedAt         *string                `json:"activated_at"`                // 激活时间
	LastHeartbeat       *string                `json:"last_heartbeat"`              // 最后心跳时间
	LastOnlineIP        *string                `json:"last_online_ip"`              // 最后在线IP
	ConfigUpdatedAt     *string                `json:"config_updated_at"`           // 客户端配置更新时间
	UsageData           map[string]interface{} `json:"usage_data,omitempty"`        // 使用数据
	CreatedAt           string                 `json:"created_at"`                  // 创建时间
	UpdatedAt           string                 `json:"updated_at"`                  // 更新时间
}

// LicenseCreateRequest 手动添加许可证请求结构
type LicenseCreateRequest struct {
	AuthorizationCodeID string                 `json:"authorization_code_id" binding:"required"` // 授权码ID，必填
	HardwareFingerprint string                 `json:"hardware_fingerprint" binding:"required"`  // 硬件指纹，必填
	DeviceInfo          map[string]interface{} `json:"device_info,omitempty"`                    // 设备信息，可选
	ActivationIP        *string                `json:"activation_ip" binding:"omitempty,ip"`     // 激活IP，可选
}

// LicenseRevokeRequest 撤销许可证请求结构
type LicenseRevokeRequest struct {
	Reason string `json:"reason" binding:"omitempty,max=500"` // 撤销原因，可选
}

// ActivateRequest 软件激活请求结构
type ActivateRequest struct {
	AuthorizationCode   string                 `json:"authorization_code" binding:"required"`   // 授权码，必填
	HardwareFingerprint string                 `json:"hardware_fingerprint" binding:"required"` // 硬件指纹，必填
	DeviceInfo          map[string]interface{} `json:"device_info,omitempty"`                   // 设备信息，可选
	SoftwareVersion     *string                `json:"software_version" binding:"omitempty"`    // 软件版本，可选
}

// ActivateResponse 软件激活响应结构
type ActivateResponse struct {
	LicenseKey        string `json:"license_key"`        // 许可证密钥
	LicenseFile       string `json:"license_file"`       // base64编码的加密许可证文件
	HeartbeatInterval int    `json:"heartbeat_interval"` // 心跳间隔(秒)
}

// HeartbeatRequest 心跳检测请求结构
type HeartbeatRequest struct {
	LicenseKey          string                 `json:"license_key" binding:"required"`          // 许可证密钥，必填
	HardwareFingerprint string                 `json:"hardware_fingerprint" binding:"required"` // 硬件指纹，必填
	ConfigUpdatedAt     *string                `json:"config_updated_at,omitempty"`             // 客户端配置更新时间，可选
	UsageData           map[string]interface{} `json:"usage_data,omitempty"`                    // 使用数据，可选
	SoftwareVersion     *string                `json:"software_version,omitempty"`              // 软件版本，可选
}

// HeartbeatResponse 心跳检测响应结构
type HeartbeatResponse struct {
	Status            string  `json:"status"`             // 许可证状态
	ConfigUpdated     bool    `json:"config_updated"`     // 配置是否有更新
	LicenseFile       *string `json:"license_file"`       // base64编码的新许可证文件(如有更新)
	HeartbeatInterval int     `json:"heartbeat_interval"` // 下次心跳间隔(秒)
}

// StatsOverviewResponse 授权概览统计响应结构
type StatsOverviewResponse struct {
	TotalAuthCodes int64      `json:"total_auth_codes"` // 总授权码数量（所有授权码的累计总数）
	ActiveLicenses int64      `json:"active_licenses"`  // 活跃许可证数量（状态为active的许可证总数）
	ExpiringSoon   int64      `json:"expiring_soon"`    // 即将过期数量（30天内即将过期的授权码，且未被锁定）
	AbnormalAlerts int64      `json:"abnormal_alerts"`  // 异常告警数量（活跃许可证中心跳超时的数量，超过配置的心跳超时时间未收到心跳）
	GrowthRate     GrowthRate `json:"growth_rate"`      // 增长率（同比上月）
}

// GrowthRate 增长率结构（同比上月：当前值相比一个月前同一时刻值的增长率）
type GrowthRate struct {
	AuthCodes float64 `json:"auth_codes"` // 授权码增长率(%)（当前总授权码数相比一个月前总授权码数的增长率）
	Licenses  float64 `json:"licenses"`   // 许可证增长率(%)（当前活跃许可证数相比一个月前活跃许可证数的增长率）
}

// DeviceListRequest 设备列表查询请求结构（客户用户接口）
type DeviceListRequest struct {
	Page                int     `form:"page" binding:"omitempty,min=1"`                   // 页码，默认1
	PageSize            int     `form:"page_size" binding:"omitempty,min=1,max=100"`      // 每页数量，默认20，最大100
	DeviceName          string  `form:"device_name" binding:"omitempty,max=100"`          // 设备名称模糊搜索（匹配device_info.name）
	AuthorizationCodeID string  `form:"authorization_code_id" binding:"omitempty,len=36"` // 按授权码ID筛选设备
	IsOnline            *string `form:"is_online" binding:"omitempty"`                    // 是否在线筛选（true在线，false离线；不传不过滤）
}

// DeviceListItem 设备列表项结构（客户用户接口）
type DeviceListItem struct {
	ID                string                 `json:"id"`                 // 许可证ID
	DeviceInfo        map[string]interface{} `json:"device_info"`        // 设备信息
	IsOnline          bool                   `json:"is_online"`          // 是否在线
	LastOnlineIP      *string                `json:"last_online_ip"`     // 最后在线IP
	LastHeartbeat     *string                `json:"last_heartbeat"`     // 最后心跳时间
	ActivatedAt       *string                `json:"activated_at"`       // 激活时间
	AuthorizationInfo AuthorizationInfo      `json:"authorization_info"` // 授权信息
}

// AuthorizationInfo 授权信息结构
type AuthorizationInfo struct {
	AuthorizationCode   string `json:"authorization_code"`    // 授权码
	AuthorizationCodeID string `json:"authorization_code_id"` // 授权码ID
	EndDate             string `json:"end_date"`              // 授权到期时间
	Description         string `json:"description"`           // 授权描述
}

// DeviceListResponse 设备列表响应结构（客户用户接口）
type DeviceListResponse struct {
	List     []DeviceListItem `json:"list"`      // 设备列表
	Total    int64            `json:"total"`     // 总数量
	Page     int              `json:"page"`      // 当前页码
	PageSize int              `json:"page_size"` // 每页数量
}

// DeviceSummaryResponse 设备汇总响应结构
type DeviceSummaryResponse struct {
	TotalDevices   int64 `json:"total_devices"`   // 设备总数
	OnlineDevices  int64 `json:"online_devices"`  // 在线设备数
	OfflineDevices int64 `json:"offline_devices"` // 离线设备数
}
