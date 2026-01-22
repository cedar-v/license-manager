package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CuUser 客户用户模型
type CuUser struct {
	ID               string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	CustomerID       string         `gorm:"type:varchar(36);not null;index" json:"customer_id"`
	Phone            string         `gorm:"type:varchar(20);not null" json:"phone"`
	PhoneCountryCode string         `gorm:"type:varchar(10);not null;default:'+86'" json:"phone_country_code"`
	Password         *string        `gorm:"type:varchar(255)" json:"-"` // 不在JSON中显示密码
	Salt             *string        `gorm:"type:varchar(32)" json:"-"`
	UserRole         string         `gorm:"type:varchar(20);not null;default:'member'" json:"user_role"`
	RealName         *string        `gorm:"type:varchar(100)" json:"real_name"`
	Email            *string        `gorm:"type:varchar(255);index" json:"email"`
	Status           string         `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	PhoneVerified    bool           `gorm:"not null;default:true" json:"phone_verified"`
	EmailVerified    bool           `gorm:"not null;default:false" json:"email_verified"`
	LastLoginAt      *time.Time     `gorm:"index" json:"last_login_at"`
	LastLoginIP      *string        `gorm:"type:varchar(50)" json:"last_login_ip"`
	LoginAttempts    int            `gorm:"not null;default:0" json:"login_attempts"`
	LockedUntil      *time.Time     `gorm:"index" json:"locked_until"`
	AvatarURL        *string        `gorm:"type:varchar(500)" json:"avatar_url"`
	Language         string         `gorm:"type:varchar(10);not null;default:'zh-CN'" json:"language"`
	Timezone         string         `gorm:"type:varchar(50);not null;default:'Asia/Shanghai'" json:"timezone"`
	AdditionalInfo   *string        `gorm:"type:json" json:"additional_info"`
	Remark           *string        `gorm:"type:text" json:"remark"`
	CreatedAt        time.Time      `gorm:"not null;index" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (CuUser) TableName() string {
	return "cu_users"
}

// BeforeCreate 创建前自动设置时间戳和ID
func (u *CuUser) BeforeCreate(tx *gorm.DB) error {
	// 生成UUID作为主键ID
	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	// 设置时间戳
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = now
	}
	return nil
}

// BeforeUpdate 更新前自动刷新更新时间
func (u *CuUser) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// IsAccountLocked 检查账号是否被锁定
func (u *CuUser) IsAccountLocked() bool {
	if u.Status == "locked" {
		return true
	}
	if u.LockedUntil != nil && time.Now().Before(*u.LockedUntil) {
		return true
	}
	return false
}

// CuUserRegisterRequest 客户用户注册请求结构
// CuUserRegisterRequest 客户用户注册请求结构
type CuUserRegisterRequest struct {
	Phone            string  `json:"phone" binding:"required,len=11"`        // 手机号，必填，11位手机号
	PhoneCountryCode string  `json:"phone_country_code" binding:"omitempty"` // 国家区号，可选，默认+86
	SmsCode          string  `json:"sms_code" binding:"required,len=6"`      // 短信验证码，必填，6位数字
	CustomerID       string  `json:"customer_id" binding:"omitempty"`        // 客户ID，可选，不提供则自动创建客户
	Password         string  `json:"password" binding:"required,min=8,max=50"` // 密码，必填，8-50位
	RealName         *string `json:"real_name" binding:"omitempty,max=100"`  // 真实姓名，可选，最大100字符
	Email            *string `json:"email" binding:"omitempty,email,max=255"` // 邮箱，可选，格式为邮箱地址
}

// CuUserLoginRequest 客户用户登录请求结构
type CuUserLoginRequest struct {
	Phone            string `json:"phone" binding:"required,min=7,max=20"`           // 手机号，必填，7-20位
	PhoneCountryCode string `json:"phone_country_code" binding:"omitempty"`          // 国家区号，可选，默认+86
	Password         string `json:"password"`                                        // 密码，密码登录时必填
	SmsCode          string `json:"sms_code"`                                        // 短信验证码，验证码登录时必填
	LoginType        string `json:"login_type" binding:"required,oneof=password sms"` // 登录类型，必填，password-密码登录/sms-验证码登录
}

// CuUserProfileUpdateRequest 客户用户资料更新请求结构
type CuUserProfileUpdateRequest struct {
	RealName       *string `json:"real_name" binding:"omitempty,max=100"`              // 真实姓名，可选，最大100字符
	Email          *string `json:"email" binding:"omitempty,email,max=255"`           // 邮箱，可选，格式为邮箱地址
	AvatarURL      *string `json:"avatar_url" binding:"omitempty,max=500"`            // 头像URL，可选，最大500字符
	Language       *string `json:"language" binding:"omitempty,oneof=zh-CN en-US ja-JP"` // 语言偏好，可选，zh-CN-中文/en-US-英文/ja-JP-日文
	Timezone       *string `json:"timezone" binding:"omitempty,max=50"`               // 时区，可选，最大50字符
	AdditionalInfo *string `json:"additional_info" binding:"omitempty"`                // 附加信息，可选，JSON格式
	Remark         *string `json:"remark" binding:"omitempty,max=1000"`                // 备注，可选，最大1000字符
}

// CuUserPhoneUpdateRequest 客户用户手机号更新请求结构
type CuUserPhoneUpdateRequest struct {
	NewPhone            string `json:"new_phone" binding:"required,len=11"`        // 新手机号，必填，11位手机号
	NewPhoneCountryCode string `json:"new_phone_country_code" binding:"required"`  // 新手机号国家区号，必填
	CurrentSmsCode      string `json:"current_sms_code" binding:"required,len=6"` // 当前手机号验证码，必填，6位数字
	NewSmsCode          string `json:"new_sms_code" binding:"required,len=6"`     // 新手机号验证码，必填，6位数字
}

// CuUserSendRegisterSmsRequest 注册发送验证码请求结构
type CuUserSendRegisterSmsRequest struct {
	Phone            string `json:"phone" binding:"required,min=7,max=20"`      // 手机号，必填，7-20位
	PhoneCountryCode string `json:"phone_country_code" binding:"omitempty"`     // 国家区号，可选，默认+86
}

// CuUserSendLoginSmsRequest 发送登录验证码请求结构
type CuUserSendLoginSmsRequest struct {
	Phone            string `json:"phone" binding:"required,min=7,max=20"`      // 手机号，必填，7-20位
	PhoneCountryCode string `json:"phone_country_code" binding:"omitempty"`     // 国家区号，可选，默认+86
}

// CuUserSendCurrentPhoneSmsRequest 发送当前手机号验证码请求结构
type CuUserSendCurrentPhoneSmsRequest struct {
	// 无需额外参数，从JWT token中获取当前用户手机号
}

// CuUserSendNewPhoneSmsRequest 发送新手机号验证码请求结构
type CuUserSendNewPhoneSmsRequest struct {
	NewPhone            string `json:"new_phone" binding:"required,min=7,max=20"` // 新手机号，必填，7-20位
	NewPhoneCountryCode string `json:"new_phone_country_code" binding:"omitempty"`  // 新手机号国家区号，可选，默认+86
}

// CuUserForgotPasswordRequest 忘记密码请求结构（发送重置验证码）
type CuUserForgotPasswordRequest struct {
	Phone            string `json:"phone" binding:"required,min=7,max=20"`      // 手机号，必填，7-20位
	PhoneCountryCode string `json:"phone_country_code" binding:"omitempty"`     // 国家区号，可选，默认+86
}

// CuUserResetPasswordRequest 重置密码请求结构
type CuUserResetPasswordRequest struct {
	Phone            string `json:"phone" binding:"required,min=7,max=20"`      // 手机号，必填，7-20位
	PhoneCountryCode string `json:"phone_country_code" binding:"omitempty"`     // 国家区号，可选，默认+86
	SmsCode          string `json:"sms_code" binding:"required,len=6"`          // 短信验证码，必填，6位数字
	NewPassword      string `json:"new_password" binding:"required,min=8,max=50"` // 新密码，必填，8-50位
}

// CuUserChangePasswordRequest 修改密码请求结构
type CuUserChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`                 // 旧密码，必填
	NewPassword string `json:"new_password" binding:"required,min=8,max=50"`     // 新密码，必填，8-50位
}

// CuUserResponse 用户信息响应结构（用于API，避免暴露敏感字段）
type CuUserResponse struct {
	ID               string  `json:"id"`                              // 用户ID
	CustomerID       string  `json:"customer_id"`                    // 客户ID
	Phone            string  `json:"phone"`                          // 手机号
	PhoneCountryCode string  `json:"phone_country_code"`             // 国家区号
	UserRole         string  `json:"user_role"`                      // 用户角色，admin-管理员/member-成员
	RealName         *string `json:"real_name"`                      // 真实姓名
	Email            *string `json:"email"`                          // 邮箱
	Status           string  `json:"status"`                         // 账户状态，active-激活/disabled-禁用
	PhoneVerified    bool    `json:"phone_verified"`                // 手机号是否已验证
	EmailVerified    bool    `json:"email_verified"`                // 邮箱是否已验证
	LastLoginAt      *string `json:"last_login_at,omitempty"`        // 最后登录时间
	LastLoginIP      *string `json:"last_login_ip,omitempty"`        // 最后登录IP
	AvatarURL        *string `json:"avatar_url"`                    // 头像URL
	Language         string  `json:"language"`                      // 语言偏好
	Timezone         string  `json:"timezone"`                      // 时区
	AdditionalInfo   *string `json:"additional_info"`               // 附加信息，JSON格式
	Remark           *string `json:"remark"`                         // 备注
	CreatedAt        string  `json:"created_at"`                    // 创建时间
	UpdatedAt        string  `json:"updated_at"`                    // 更新时间
}

// ToResponse 转换为响应结构
func (u *CuUser) ToResponse() *CuUserResponse {
	response := &CuUserResponse{
		ID:               u.ID,
		CustomerID:       u.CustomerID,
		Phone:            u.Phone,
		PhoneCountryCode: u.PhoneCountryCode,
		UserRole:         u.UserRole,
		RealName:         u.RealName,
		Email:            u.Email,
		Status:           u.Status,
		PhoneVerified:    u.PhoneVerified,
		EmailVerified:    u.EmailVerified,
		AvatarURL:        u.AvatarURL,
		Language:         u.Language,
		Timezone:         u.Timezone,
		AdditionalInfo:   u.AdditionalInfo,
		Remark:           u.Remark,
		CreatedAt:        u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// 格式化时间字段
	if u.LastLoginAt != nil {
		timeStr := u.LastLoginAt.Format("2006-01-02T15:04:05Z07:00")
		response.LastLoginAt = &timeStr
	}
	if u.LastLoginIP != nil {
		response.LastLoginIP = u.LastLoginIP
	}

	return response
}
