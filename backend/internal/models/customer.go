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
	ContactPerson        string         `gorm:"type:varchar(100);not null" json:"contact_person"`
	ContactTitle         *string        `gorm:"type:varchar(100)" json:"contact_title"`
	Email                string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Phone                *string        `gorm:"type:varchar(20)" json:"phone"`
	Address              *string        `gorm:"type:text" json:"address"`
	CompanySize          *string        `gorm:"type:varchar(20)" json:"company_size"`
	PreferredLicenseType string         `gorm:"type:varchar(20);default:'online'" json:"preferred_license_type"`
	CustomerLevel        string         `gorm:"type:varchar(20);not null;default:'normal';index" json:"customer_level"`
	Status               string         `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
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
