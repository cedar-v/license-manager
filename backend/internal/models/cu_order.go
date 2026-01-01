package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CuOrder 客户用户订单模型
type CuOrder struct {
	ID                string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	OrderNo           string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"order_no"`
	CustomerID        string         `gorm:"type:varchar(36);not null;index" json:"customer_id"`
	CuUserID          string         `gorm:"type:varchar(36);not null;index" json:"cu_user_id"`
	PackageID         string         `gorm:"type:varchar(50);not null;index" json:"package_id"`
	PackageName       string         `gorm:"type:varchar(100);not null" json:"package_name"`
	LicenseCount      int            `gorm:"not null" json:"license_count"`
	UnitPrice         float64        `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	DiscountRate      float64        `gorm:"type:decimal(3,2);not null;default:1.0" json:"discount_rate"`
	TotalAmount       float64        `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status            string         `gorm:"type:varchar(20);not null;default:'paid'" json:"status"`
	AuthorizationCode *string        `gorm:"type:varchar(500)" json:"authorization_code"`
	ExpiredAt         *time.Time     `gorm:"index" json:"expired_at"`
	CreatedAt         time.Time      `gorm:"not null;index" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (CuOrder) TableName() string {
	return "cu_orders"
}

// BeforeCreate 创建前自动设置时间戳和ID
func (o *CuOrder) BeforeCreate(tx *gorm.DB) error {
	// 生成UUID作为主键ID
	if o.ID == "" {
		o.ID = uuid.New().String()
	}

	// 设置时间戳
	now := time.Now()
	if o.CreatedAt.IsZero() {
		o.CreatedAt = now
	}
	if o.UpdatedAt.IsZero() {
		o.UpdatedAt = now
	}

	return nil
}

// BeforeUpdate 更新前自动设置时间戳
func (o *CuOrder) BeforeUpdate(tx *gorm.DB) error {
	o.UpdatedAt = time.Now()
	return nil
}

// CuOrderCreateRequest 创建订单请求结构
type CuOrderCreateRequest struct {
	PackageID    string `json:"package_id" binding:"required"`
	LicenseCount int    `json:"license_count" binding:"required,min=1,max=1000"`
}

// CuOrderResponse 订单响应结构
type CuOrderResponse struct {
	ID                string     `json:"id"`
	OrderNo           string     `json:"order_no"`
	CustomerID        string     `json:"customer_id"`
	CuUserID          string     `json:"cu_user_id"`
	PackageID         string     `json:"package_id"`
	PackageName       string     `json:"package_name"`
	LicenseCount      int        `json:"license_count"`
	UnitPrice         float64    `json:"unit_price"`
	DiscountRate      float64    `json:"discount_rate"`
	TotalAmount       float64    `json:"total_amount"`
	Status            string     `json:"status"`
	AuthorizationCode *string    `json:"authorization_code"`
	ExpiredAt         *time.Time `json:"expired_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// ToResponse 转换为响应结构
func (o *CuOrder) ToResponse() *CuOrderResponse {
	return &CuOrderResponse{
		ID:                o.ID,
		OrderNo:           o.OrderNo,
		CustomerID:        o.CustomerID,
		CuUserID:          o.CuUserID,
		PackageID:         o.PackageID,
		PackageName:       o.PackageName,
		LicenseCount:      o.LicenseCount,
		UnitPrice:         o.UnitPrice,
		DiscountRate:      o.DiscountRate,
		TotalAmount:       o.TotalAmount,
		Status:            o.Status,
		AuthorizationCode: o.AuthorizationCode,
		ExpiredAt:         o.ExpiredAt,
		CreatedAt:         o.CreatedAt,
		UpdatedAt:         o.UpdatedAt,
	}
}

// CuOrderListResponse 订单列表响应结构
type CuOrderListResponse struct {
	Orders     []*CuOrderResponse `json:"orders"`
	TotalCount int64              `json:"total_count"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
}
