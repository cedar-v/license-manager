package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Code    string      `json:"code"`           // 响应代码
	Message string      `json:"message"`        // 响应消息
	Data    interface{} `json:"data,omitempty"` // 响应数据
}

// LoginData 登录成功返回的数据
type LoginData struct {
	Token     string   `json:"token"`      // JWT令牌
	ExpiresIn int      `json:"expires_in"` // 过期时间（秒）
	UserInfo  UserInfo `json:"user_info"`  // 用户信息
}

// UserInfo 用户信息
type UserInfo struct {
	Username string `json:"username"` // 用户名
	Role     string `json:"role"`     // 角色
}

// AuthorizationCode 授权码模型
type AuthorizationCode struct {
	ID                     string                   `gorm:"type:varchar(36);primaryKey" json:"id"`                                 // 授权码ID
	Code                   string                   `gorm:"type:varchar(1000);not null" json:"code"`                               // 授权码
	CustomerID             string                   `gorm:"type:varchar(36);not null;index" json:"customer_id"`                    // 客户ID
	CustomerName           string                   `gorm:"-" json:"customer_name,omitempty"`                                      // 客户名称
	CustomerNameDisplay    string                   `gorm:"-" json:"customer_name_display,omitempty"`                              // 客户名称显示（多语言）
	CreatedBy              string                   `gorm:"type:varchar(36);not null" json:"created_by"`                           // 创建人ID
	SoftwareID             *string                  `gorm:"type:varchar(50)" json:"software_id"`                                   // 软件ID
	Description            *string                  `gorm:"type:text" json:"description"`                                          // 描述
	StartDate              time.Time                `gorm:"type:datetime(3);not null" json:"start_date"`                           // 生效日期
	EndDate                time.Time                `gorm:"type:datetime(3);not null" json:"end_date"`                             // 失效日期
	DeploymentType         string                   `gorm:"type:varchar(20);not null;default:'standalone'" json:"deployment_type"` // 部署类型：standalone/cloud/hybrid
	DeploymentTypeDisplay  string                   `gorm:"-" json:"deployment_type_display,omitempty"`                            // 部署类型显示（多语言）
	EncryptionType         *string                  `gorm:"type:varchar(20);default:'standard'" json:"encryption_type"`            // 加密类型：standard/advanced
	EncryptionTypeDisplay  string                   `gorm:"-" json:"encryption_type_display,omitempty"`                            // 加密类型显示（多语言）
	SoftwareVersion        *string                  `gorm:"type:varchar(50)" json:"software_version"`                              // 软件版本
	MaxActivations         int                      `gorm:"not null;default:1" json:"max_activations"`                             // 最大激活次数
	CurrentActivations     int                      `gorm:"-" json:"current_activations,omitempty"`                                // 当前激活次数
	IsLocked               bool                     `gorm:"not null;default:false" json:"is_locked"`                               // 是否锁定
	LockReason             *string                  `gorm:"type:text" json:"lock_reason"`                                          // 锁定原因
	LockedAt               *time.Time               `gorm:"type:datetime(3)" json:"locked_at"`                                     // 锁定时间
	LockedBy               *string                  `gorm:"type:varchar(36)" json:"locked_by"`                                     // 锁定人ID
	FeatureConfig          JSON                     `gorm:"type:json" json:"feature_config" swaggertype:"object"`                  // 功能配置（JSON对象）
	UsageLimits            JSON                     `gorm:"type:json" json:"usage_limits" swaggertype:"object"`                    // 使用限制（JSON对象）
	CustomParameters       JSON                     `gorm:"type:json" json:"custom_parameters" swaggertype:"object"`               // 自定义参数（JSON对象）
	Status                 string                   `gorm:"-" json:"status,omitempty"`                                             // 状态：normal/locked/expired
	StatusDisplay          string                   `gorm:"-" json:"status_display,omitempty"`                                     // 状态显示（多语言）
	CreatedAt              time.Time                `gorm:"type:datetime(3);not null" json:"created_at"`                           // 创建时间
	UpdatedAt              time.Time                `gorm:"type:datetime(3);not null" json:"updated_at"`                           // 更新时间
	CustomerInfo           *CustomerInfoForAuthCode `gorm:"-" json:"customer_info,omitempty"`                                      // 客户信息（仅在详情接口返回）
	ActivatedLicensesCount int64                    `gorm:"-" json:"activated_licenses_count,omitempty"`                           // 该授权码下已激活的许可证数量
}

// CustomerInfoForAuthCode 授权码详情中的客户信息结构
type CustomerInfoForAuthCode struct {
	ID                  string `json:"id"`                              // 客户ID
	CustomerCode        string `json:"customer_code"`                   // 客户编号
	CustomerName        string `json:"customer_name"`                   // 客户名称
	CustomerType        string `json:"customer_type"`                   // 客户类型
	CustomerTypeDisplay string `json:"customer_type_display,omitempty"` // 客户类型显示（多语言）
	Status              string `json:"status"`                          // 客户状态
	StatusDisplay       string `json:"status_display,omitempty"`        // 客户状态显示（多语言）
	CreatedAt           string `json:"created_at"`                      // 客户创建时间
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
	CustomerID       string      `json:"customer_id" binding:"required"`                                   // 客户ID
	SoftwareID       *string     `json:"software_id" binding:"omitempty"`                                  // 软件ID
	Description      *string     `json:"description" binding:"omitempty,max=1000"`                         // 描述
	ValidityDays     int         `json:"validity_days" binding:"required,min=1,max=365000"`                // 有效天数（1-365000天，365000代表永久有效）
	DeploymentType   string      `json:"deployment_type" binding:"required,oneof=standalone cloud hybrid"` // 部署类型：standalone/cloud/hybrid
	EncryptionType   *string     `json:"encryption_type" binding:"omitempty,oneof=standard advanced"`      // 加密类型：standard/advanced
	SoftwareVersion  *string     `json:"software_version" binding:"omitempty"`                             // 软件版本
	MaxActivations   int         `json:"max_activations" binding:"required,min=1"`                         // 最大激活次数
	FeatureConfig    interface{} `json:"feature_config" binding:"omitempty"`                               // 功能配置（JSON对象）
	UsageLimits      interface{} `json:"usage_limits" binding:"omitempty"`                                 // 使用限制（JSON对象）
	CustomParameters interface{} `json:"custom_parameters" binding:"omitempty"`                            // 自定义参数（JSON对象）
}

// AuthorizationCodeCreateResponse 创建授权码响应结构
type AuthorizationCodeCreateResponse struct {
	ID   string `json:"id"`   // 授权码ID
	Code string `json:"code"` // 授权码
}

// AuthorizationCodeListRequest 授权码列表查询请求结构
type AuthorizationCodeListRequest struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`                            // 页码，默认1
	PageSize   int    `form:"page_size" binding:"omitempty,min=1,max=100"`               // 每页条数，默认20，最大100
	CustomerID string `form:"customer_id" binding:"omitempty"`                           // 客户ID筛选
	Status     string `form:"status" binding:"omitempty,oneof=normal locked expired"`    // 状态筛选
	StartDate  string `form:"start_date" binding:"omitempty"`                            // 创建开始时间
	EndDate    string `form:"end_date" binding:"omitempty"`                              // 创建结束时间
	Sort       string `form:"sort" binding:"omitempty,oneof=created_at updated_at code"` // 排序字段，默认created_at
	Order      string `form:"order" binding:"omitempty,oneof=asc desc"`                  // 排序方向，默认desc
}

// AuthorizationCodeListItem 授权码列表项结构
type AuthorizationCodeListItem struct {
	ID                    string  `json:"id"`                                // 授权码ID
	Code                  string  `json:"code"`                              // 授权码
	CustomerID            string  `json:"customer_id"`                       // 客户ID
	CustomerName          string  `json:"customer_name"`                     // 客户名称
	CustomerNameDisplay   string  `json:"customer_name_display,omitempty"`   // 客户名称显示（多语言）
	Status                string  `json:"status"`                            // 状态：normal/locked/expired
	StatusDisplay         string  `json:"status_display,omitempty"`          // 状态显示（多语言）
	StartDate             string  `json:"start_date"`                        // 生效日期
	EndDate               string  `json:"end_date"`                          // 失效日期
	MaxActivations        int     `json:"max_activations"`                   // 最大激活次数
	CurrentActivations    int     `json:"current_activations"`               // 当前激活次数
	DeploymentType        string  `json:"deployment_type"`                   // 部署类型：standalone/cloud/hybrid
	DeploymentTypeDisplay string  `json:"deployment_type_display,omitempty"` // 部署类型显示（多语言）
	IsLocked              bool    `json:"is_locked"`                         // 是否锁定
	Description           *string `json:"description"`                       // 描述
	CreatedAt             string  `json:"created_at"`                        // 创建时间
}

// AuthorizationCodeListResponse 授权码列表响应结构
type AuthorizationCodeListResponse struct {
	List       []AuthorizationCodeListItem `json:"list"`        // 授权码列表
	Total      int64                       `json:"total"`       // 总记录数
	Page       int                         `json:"page"`        // 当前页码
	PageSize   int                         `json:"page_size"`   // 每页条数
	TotalPages int                         `json:"total_pages"` // 总页数
}

// AuthorizationCodeUpdateRequest 更新授权码请求结构
type AuthorizationCodeUpdateRequest struct {
	SoftwareID       *string     `json:"software_id" binding:"omitempty"`                                   // 软件ID
	Description      *string     `json:"description" binding:"omitempty,max=1000"`                          // 描述
	ValidityDays     *int        `json:"validity_days" binding:"omitempty,min=1,max=365000"`                // 有效天数（1-365000天，365000代表永久有效）
	DeploymentType   *string     `json:"deployment_type" binding:"omitempty,oneof=standalone cloud hybrid"` // 部署类型：standalone/cloud/hybrid
	EncryptionType   *string     `json:"encryption_type" binding:"omitempty,oneof=standard advanced"`       // 加密类型：standard/advanced
	SoftwareVersion  *string     `json:"software_version" binding:"omitempty"`                              // 软件版本
	MaxActivations   *int        `json:"max_activations" binding:"omitempty,min=1"`                         // 最大激活次数
	FeatureConfig    interface{} `json:"feature_config" binding:"omitempty"`                                // 功能配置
	UsageLimits      interface{} `json:"usage_limits" binding:"omitempty"`                                  // 使用限制
	CustomParameters interface{} `json:"custom_parameters" binding:"omitempty"`                             // 自定义参数
	// 可选的起止时间（优先于 validity_days），格式：YYYY-MM-DD
	StartDate  *string `json:"start_date" binding:"omitempty"`                                                      // 生效日期（YYYY-MM-DD）
	EndDate    *string `json:"end_date" binding:"omitempty"`                                                        // 失效日期（YYYY-MM-DD）
	ChangeType string  `json:"change_type" binding:"required,oneof=renewal feature_limit_change lock unlock other"` // 变更类型：renewal/upgrade/limit_change/feature_toggle/lock/unlock/other
	Reason     *string `json:"reason" binding:"omitempty,max=500"`                                                  // 变更原因
}

// AuthorizationCodeLockRequest 锁定/解锁授权码请求结构
type AuthorizationCodeLockRequest struct {
	IsLocked   bool    `json:"is_locked" binding:""`                    // true-锁定，false-解锁
	LockReason *string `json:"lock_reason" binding:"omitempty,max=500"` // 锁定原因
	Reason     *string `json:"reason" binding:"omitempty,max=500"`      // 变更原因（记录到变更历史）
}

// AuthorizationChange 授权变更历史模型
type AuthorizationChange struct {
	ID                  string    `gorm:"type:varchar(36);primaryKey" json:"id"`                        // 变更记录ID
	AuthorizationCodeID string    `gorm:"type:varchar(36);not null;index" json:"authorization_code_id"` // 授权码ID
	ChangeType          string    `gorm:"type:varchar(30);not null;index" json:"change_type"`           // 变更类型
	ChangeTypeDisplay   string    `gorm:"-" json:"change_type_display,omitempty"`                       // 变更类型显示（多语言）
	OldConfig           JSON      `gorm:"type:json" json:"old_config,omitempty" swaggertype:"object"`   // 变更前配置（JSON对象）
	NewConfig           JSON      `gorm:"type:json" json:"new_config,omitempty" swaggertype:"object"`   // 变更后配置（JSON对象）
	OperatorID          string    `gorm:"type:varchar(36);not null;index" json:"operator_id"`           // 操作人ID
	OperatorName        string    `gorm:"-" json:"operator_name,omitempty"`                             // 操作人名称
	Reason              *string   `gorm:"type:text" json:"reason"`                                      // 变更原因
	CreatedAt           time.Time `gorm:"type:datetime(3);not null;index" json:"created_at"`            // 创建时间

	// 关联字段
	AuthorizationCode *AuthorizationCode `gorm:"foreignKey:AuthorizationCodeID" json:"authorization_code,omitempty"` // 关联的授权码
	Operator          *User              `gorm:"foreignKey:OperatorID" json:"operator,omitempty"`                    // 关联的操作人
}

// TableName 指定表名
func (AuthorizationChange) TableName() string {
	return "authorization_changes"
}

// BeforeCreate 创建前自动设置时间戳和ID
func (a *AuthorizationChange) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}

	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}
	return nil
}

// AuthorizationChangeListRequest 授权变更历史列表查询请求结构
type AuthorizationChangeListRequest struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`                        // 页码，默认1
	PageSize   int    `form:"page_size" binding:"omitempty,min=1,max=100"`           // 每页条数，默认20，最大100
	ChangeType string `form:"change_type" binding:"omitempty"`                       // 变更类型筛选
	OperatorID string `form:"operator_id" binding:"omitempty"`                       // 操作人筛选
	StartDate  string `form:"start_date" binding:"omitempty"`                        // 开始时间
	EndDate    string `form:"end_date" binding:"omitempty"`                          // 结束时间
	Sort       string `form:"sort" binding:"omitempty,oneof=created_at change_type"` // 排序字段，默认created_at
	Order      string `form:"order" binding:"omitempty,oneof=asc desc"`              // 排序方向，默认desc
}

// AuthorizationChangeListItem 授权变更历史列表项结构
type AuthorizationChangeListItem struct {
	ID                string  `json:"id"`                            // 变更记录ID
	ChangeType        string  `json:"change_type"`                   // 变更类型
	ChangeTypeDisplay string  `json:"change_type_display,omitempty"` // 变更类型显示（多语言）
	OperatorID        string  `json:"operator_id"`                   // 操作人ID
	OperatorName      string  `json:"operator_name,omitempty"`       // 操作人名称
	Reason            *string `json:"reason"`                        // 变更原因
	CreatedAt         string  `json:"created_at"`                    // 创建时间
}

// AuthorizationChangeListResponse 授权变更历史列表响应结构
type AuthorizationChangeListResponse struct {
	List       []AuthorizationChangeListItem `json:"list"`        // 变更历史列表
	Total      int64                         `json:"total"`       // 总记录数
	Page       int                           `json:"page"`        // 当前页码
	PageSize   int                           `json:"page_size"`   // 每页条数
	TotalPages int                           `json:"total_pages"` // 总页数
}

// AuthorizationCodeShareRequest 授权码分享请求结构
type AuthorizationCodeShareRequest struct {
	TargetUserID string `json:"target_user_id" binding:"required"`    // 受赠用户ID
	ShareCount   int    `json:"share_count" binding:"required,min=1"` // 分享激活次数
}

// AuthorizationCodeShareResponse 授权码分享响应结构
type AuthorizationCodeShareResponse struct {
	NewAuthorizationCode AuthorizationCodeShareResponseItem `json:"new_authorization_code"` // 新生成的授权码信息
}

// AuthorizationCodeShareResponseItem 分享响应中的授权码信息
type AuthorizationCodeShareResponseItem struct {
	ID             string `json:"id"`              // 授权码ID
	Code           string `json:"code"`            // 授权码
	StartDate      string `json:"start_date"`      // 开始时间
	EndDate        string `json:"end_date"`        // 结束时间
	MaxActivations int    `json:"max_activations"` // 最大激活次数
}
