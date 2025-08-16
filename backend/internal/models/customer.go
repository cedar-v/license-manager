package models

import (
	"time"

	"gorm.io/gorm"
)

// Customer 客户信息模型
type Customer struct {
	ID                   string         `gorm:"type:varchar(36);primaryKey;default:(UUID())" json:"id"`
	CustomerCode         string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"customer_code"`
	CustomerName         string         `gorm:"type:varchar(200);not null;index" json:"customer_name"`
	CustomerType         string         `gorm:"type:varchar(20);not null;default:'enterprise';index" json:"customer_type"`
	CustomerTypeDisplay  string         `gorm:"-" json:"customer_type_display,omitempty"`
	ContactPerson        string         `gorm:"type:varchar(100);not null" json:"contact_person"`
	ContactTitle         *string        `gorm:"type:varchar(100)" json:"contact_title"`
	Email                *string        `gorm:"type:varchar(255)" json:"email"`
	Phone                *string        `gorm:"type:varchar(20)" json:"phone"`
	Address              *string        `gorm:"type:text" json:"address"`
	CompanySize          *string        `gorm:"type:varchar(20)" json:"company_size"`
	CompanySizeDisplay   string         `gorm:"-" json:"company_size_display,omitempty"`
	CustomerLevel        string         `gorm:"type:varchar(20);not null;default:'normal';index" json:"customer_level"`
	CustomerLevelDisplay string         `gorm:"-" json:"customer_level_display,omitempty"`
	Status               string         `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	StatusDisplay        string         `gorm:"-" json:"status_display,omitempty"`
	Description          *string        `gorm:"type:text" json:"description"`
	CreatedAt            time.Time      `gorm:"not null;index" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"not null" json:"updated_at"`
	CreatedBy            string         `gorm:"type:varchar(36);not null;index" json:"created_by"`
	UpdatedBy            *string        `gorm:"type:varchar(36)" json:"updated_by"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

// CustomerCodeSequence 客户编码序列模型
type CustomerCodeSequence struct {
	Year           int `gorm:"primaryKey" json:"year"`
	SequenceNumber int `gorm:"not null;default:0" json:"sequence_number"`
}

// TableName 指定表名
func (Customer) TableName() string {
	return "customers"
}

// BeforeCreate 创建前自动设置时间戳
// 确保数据完整性，避免依赖数据库默认值（符合失败快速原则）
func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	if c.UpdatedAt.IsZero() {
		c.UpdatedAt = now
	}
	return nil
}

// BeforeUpdate 更新前自动刷新更新时间
// 保持数据时间戳的准确性和一致性
func (c *Customer) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}

// TableName 指定表名
func (CustomerCodeSequence) TableName() string {
	return "customer_code_sequence"
}

// CustomerListRequest 客户列表查询请求结构
type CustomerListRequest struct {
	Page         int    `form:"page" binding:"omitempty,min=1"`                              // 页码，默认1
	PageSize     int    `form:"page_size" binding:"omitempty,min=1,max=100"`               // 每页条数，默认20，最大100
	Search       string `form:"search" binding:"omitempty,max=100"`                        // 搜索关键词(支持客户编码、名称、联系人、邮箱)
	CustomerType string `form:"customer_type" binding:"omitempty,oneof=individual enterprise government education"` // 客户类型筛选
	CustomerLevel string `form:"customer_level" binding:"omitempty,oneof=normal vip enterprise strategic"` // 客户等级筛选
	Status       string `form:"status" binding:"omitempty,oneof=active disabled"`          // 状态筛选
	Sort         string `form:"sort" binding:"omitempty,oneof=created_at updated_at customer_name customer_code"` // 排序字段，默认created_at
	Order        string `form:"order" binding:"omitempty,oneof=asc desc"`                  // 排序方向，默认desc
}

// CustomerListItem 客户列表项结构（用于列表展示，包含主要字段）
type CustomerListItem struct {
	ID                   string  `json:"id"`
	CustomerCode         string  `json:"customer_code"`
	CustomerName         string  `json:"customer_name"`
	CustomerType         string  `json:"customer_type"`
	CustomerTypeDisplay  string  `json:"customer_type_display,omitempty"`
	ContactPerson        string  `json:"contact_person"`
	Email                *string `json:"email"`
	CustomerLevel        string  `json:"customer_level"`
	CustomerLevelDisplay string  `json:"customer_level_display,omitempty"`
	Status               string  `json:"status"`
	StatusDisplay        string  `json:"status_display,omitempty"`
	CreatedAt            string  `json:"created_at"`
}

// CustomerListResponse 客户列表响应结构
type CustomerListResponse struct {
	List       []CustomerListItem `json:"list"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}
