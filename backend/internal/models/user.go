package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID            string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Username      string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email         string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash  string    `gorm:"type:varchar(255);not null" json:"-"` // 不在JSON中显示密码哈希
	FullName      string    `gorm:"type:varchar(100);not null" json:"full_name"`
	Phone         *string   `gorm:"type:varchar(20)" json:"phone"`
	Role          string    `gorm:"type:varchar(20);not null;default:'viewer';index" json:"role"`
	Status        string    `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	LastLoginAt   *time.Time `gorm:"index" json:"last_login_at"`
	LastLoginIP   *string   `gorm:"type:varchar(45)" json:"last_login_ip"`
	LoginAttempts int       `gorm:"not null;default:0" json:"login_attempts"`
	LockedUntil   *time.Time `gorm:"index" json:"locked_until"`
	CreatedAt     time.Time `gorm:"not null;index" json:"created_at"`
	UpdatedAt     time.Time `gorm:"not null" json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前自动设置时间戳和ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
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
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// IsAccountLocked 检查账号是否被锁定
func (u *User) IsAccountLocked() bool {
	if u.Status == "locked" {
		return true
	}
	if u.LockedUntil != nil && time.Now().Before(*u.LockedUntil) {
		return true
	}
	return false
}

// UserCreateRequest 创建用户请求结构
type UserCreateRequest struct {
	Username string  `json:"username" binding:"required,min=3,max=50"`
	Email    string  `json:"email" binding:"required,email,max=255"`
	Password string  `json:"password" binding:"required,min=8,max=50"`
	FullName string  `json:"full_name" binding:"required,max=100"`
	Phone    *string `json:"phone" binding:"omitempty,max=20"`
	Role     string  `json:"role" binding:"required,oneof=admin operator viewer"`
	Status   string  `json:"status" binding:"required,oneof=active disabled"`
}

// UserUpdateRequest 更新用户请求结构
type UserUpdateRequest struct {
	Email    *string `json:"email" binding:"omitempty,email,max=255"`
	FullName *string `json:"full_name" binding:"omitempty,max=100"`
	Phone    *string `json:"phone" binding:"omitempty,max=20"`
	Role     *string `json:"role" binding:"omitempty,oneof=admin operator viewer"`
	Status   *string `json:"status" binding:"omitempty,oneof=active disabled locked"`
}

// ChangePasswordRequest 修改密码请求结构
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=50"`
}

// UserListRequest 用户列表查询请求结构
type UserListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search" binding:"omitempty,max=100"`
	Role     string `form:"role" binding:"omitempty,oneof=admin operator viewer"`
	Status   string `form:"status" binding:"omitempty,oneof=active disabled locked"`
	Sort     string `form:"sort" binding:"omitempty,oneof=created_at updated_at username full_name"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// UserListItem 用户列表项结构
type UserListItem struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	FullName    string     `json:"full_name"`
	Phone       *string    `json:"phone"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// UserListResponse 用户列表响应结构
type UserListResponse struct {
	List       []UserListItem `json:"list"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}