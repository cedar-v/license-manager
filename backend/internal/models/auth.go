package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginData 登录成功返回的数据
type LoginData struct {
	Token     string   `json:"token"`
	ExpiresIn int      `json:"expires_in"`
	UserInfo  UserInfo `json:"user_info"`
}

// UserInfo 用户信息
type UserInfo struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

// AuthorizationCode 授权码模型
type AuthorizationCode struct {
	ID                   string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Code                 string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"code"`
	CustomerID           string         `gorm:"type:varchar(36);not null;index" json:"customer_id"`
	CustomerName         string         `gorm:"-" json:"customer_name,omitempty"`
	CustomerNameDisplay  string         `gorm:"-" json:"customer_name_display,omitempty"`
	CreatedBy            string         `gorm:"type:varchar(36);not null" json:"created_by"`
	SoftwareID           *string        `gorm:"type:varchar(50)" json:"software_id"`
	Description          *string        `gorm:"type:text" json:"description"`
	StartDate            time.Time      `gorm:"type:datetime(3);not null" json:"start_date"`
	EndDate              time.Time      `gorm:"type:datetime(3);not null" json:"end_date"`
	DeploymentType       string         `gorm:"type:varchar(20);not null;default:'standalone'" json:"deployment_type"`
	DeploymentTypeDisplay string         `gorm:"-" json:"deployment_type_display,omitempty"`
	EncryptionType       *string        `gorm:"type:varchar(20);default:'standard'" json:"encryption_type"`
	EncryptionTypeDisplay string         `gorm:"-" json:"encryption_type_display,omitempty"`
	SoftwareVersion      *string        `gorm:"type:varchar(50)" json:"software_version"`
	MaxActivations       int            `gorm:"not null;default:1" json:"max_activations"`
	CurrentActivations   int            `gorm:"-" json:"current_activations,omitempty"`
	IsLocked             bool           `gorm:"not null;default:false" json:"is_locked"`
	LockReason           *string        `gorm:"type:text" json:"lock_reason"`
	LockedAt             *time.Time     `gorm:"type:datetime(3)" json:"locked_at"`
	LockedBy             *string        `gorm:"type:varchar(36)" json:"locked_by"`
	FeatureConfig        JSON           `gorm:"type:json" json:"feature_config"`
	UsageLimits          JSON           `gorm:"type:json" json:"usage_limits"`
	CustomParameters     JSON           `gorm:"type:json" json:"custom_parameters"`
	Status               string         `gorm:"-" json:"status,omitempty"`
	StatusDisplay        string         `gorm:"-" json:"status_display,omitempty"`
	CreatedAt            time.Time      `gorm:"type:datetime(3);not null" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"type:datetime(3);not null" json:"updated_at"`
}

// TableName 指定表名
func (AuthorizationCode) TableName() string {
	return "authorization_codes"
}

// BeforeCreate 创建前自动设置时间戳和ID
func (a *AuthorizationCode) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}
	if a.UpdatedAt.IsZero() {
		a.UpdatedAt = now
	}
	return nil
}

// BeforeUpdate 更新前自动刷新更新时间
func (a *AuthorizationCode) BeforeUpdate(tx *gorm.DB) error {
	a.UpdatedAt = time.Now()
	return nil
}

// AuthorizationCodeCreateRequest 创建授权码请求结构
type AuthorizationCodeCreateRequest struct {
	CustomerID       string      `json:"customer_id" binding:"required"`
	SoftwareID       *string     `json:"software_id" binding:"omitempty"`
	Description      *string     `json:"description" binding:"omitempty,max=1000"`
	ValidityDays     int         `json:"validity_days" binding:"required,min=1,max=36500"`
	DeploymentType   string      `json:"deployment_type" binding:"required,oneof=standalone cloud hybrid"`
	EncryptionType   *string     `json:"encryption_type" binding:"omitempty,oneof=standard advanced"`
	SoftwareVersion  *string     `json:"software_version" binding:"omitempty"`
	MaxActivations   int         `json:"max_activations" binding:"required,min=1"`
	FeatureConfig    interface{} `json:"feature_config" binding:"omitempty"`
	UsageLimits      interface{} `json:"usage_limits" binding:"omitempty"`
	CustomParameters interface{} `json:"custom_parameters" binding:"omitempty"`
}

// AuthorizationCodeCreateResponse 创建授权码响应结构
type AuthorizationCodeCreateResponse struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

// AuthorizationCodeListRequest 授权码列表查询请求结构
type AuthorizationCodeListRequest struct {
	Page         int    `form:"page" binding:"omitempty,min=1"`                              // 页码，默认1
	PageSize     int    `form:"page_size" binding:"omitempty,min=1,max=100"`               // 每页条数，默认20，最大100
	CustomerID   string `form:"customer_id" binding:"omitempty"`                           // 客户ID筛选
	Status       string `form:"status" binding:"omitempty,oneof=normal locked expired"`    // 状态筛选
	StartDate    string `form:"start_date" binding:"omitempty"`                            // 创建开始时间
	EndDate      string `form:"end_date" binding:"omitempty"`                              // 创建结束时间
	Sort         string `form:"sort" binding:"omitempty,oneof=created_at updated_at code"` // 排序字段，默认created_at
	Order        string `form:"order" binding:"omitempty,oneof=asc desc"`                  // 排序方向，默认desc
}

// AuthorizationCodeListItem 授权码列表项结构
type AuthorizationCodeListItem struct {
	ID                     string  `json:"id"`
	Code                   string  `json:"code"`
	CustomerID             string  `json:"customer_id"`
	CustomerName           string  `json:"customer_name"`
	CustomerNameDisplay    string  `json:"customer_name_display,omitempty"`
	Status                 string  `json:"status"`
	StatusDisplay          string  `json:"status_display,omitempty"`
	StartDate              string  `json:"start_date"`
	EndDate                string  `json:"end_date"`
	MaxActivations         int     `json:"max_activations"`
	CurrentActivations     int     `json:"current_activations"`
	DeploymentType         string  `json:"deployment_type"`
	DeploymentTypeDisplay  string  `json:"deployment_type_display,omitempty"`
	IsLocked               bool    `json:"is_locked"`
	Description            *string `json:"description"`
	CreatedAt              string  `json:"created_at"`
}

// AuthorizationCodeListResponse 授权码列表响应结构
type AuthorizationCodeListResponse struct {
	List       []AuthorizationCodeListItem `json:"list"`
	Total      int64                       `json:"total"`
	Page       int                         `json:"page"`
	PageSize   int                         `json:"page_size"`
	TotalPages int                         `json:"total_pages"`
}