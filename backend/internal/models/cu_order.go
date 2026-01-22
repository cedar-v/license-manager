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
	PackageID     string `json:"package_id" binding:"required"`
	LicenseCount  int    `json:"license_count" binding:"required,min=1,max=1000"`
	PaymentMethod string `json:"payment_method,omitempty"` // 可选：支付方式，不传则为免费订单 支持：alipay，wechat
}

// CuOrderListRequest 订单列表查询请求结构
type CuOrderListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`              // 页码，默认1
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"` // 每页条数，默认10，最大100
	Search   string `form:"search" binding:"omitempty"`                  // 订单号或授权码模糊匹配
	Status   string `form:"status" binding:"omitempty"`                  // 订单状态筛选
	Time     string `form:"time" binding:"omitempty"`                    // 时间筛选：today/week/month/three_months 或 今天/本周/本月/近三个月
}

// CuOrderResponse 订单响应结构
type CuOrderResponse struct {
	ID                string     `json:"id"`                 // 订单ID
	OrderNo           string     `json:"order_no"`           // 订单号，格式如ORD202601210123456789
	CustomerID        string     `json:"customer_id"`        // 客户ID
	CuUserID          string     `json:"cu_user_id"`         // 客户用户ID
	PackageID         string     `json:"package_id"`         // 套餐ID
	PackageName       string     `json:"package_name"`       // 套餐名称
	LicenseCount      int        `json:"license_count"`      // 许可数量
	UnitPrice         float64    `json:"unit_price"`         // 单价，单位元
	DiscountRate      float64    `json:"discount_rate"`      // 折扣率，0.0-1.0之间
	TotalAmount       float64    `json:"total_amount"`       // 订单总金额，单位元
	Status            string     `json:"status"`             // 订单状态，pending-待支付/paid-已支付
	AuthorizationCode *string    `json:"authorization_code"` // 授权码，已支付订单生成
	ExpiredAt         *time.Time `json:"expired_at"`         // 订单过期时间
	CreatedAt         time.Time  `json:"created_at"`         // 创建时间
	UpdatedAt         time.Time  `json:"updated_at"`         // 更新时间
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

// OrderSummaryResponse 订单汇总响应结构
type OrderSummaryResponse struct {
	TotalOrders   int64 `json:"total_orders"`   // 订单总数
	PendingOrders int64 `json:"pending_orders"` // 待支付订单数
	PaidOrders    int64 `json:"paid_orders"`    // 已支付订单数
}
