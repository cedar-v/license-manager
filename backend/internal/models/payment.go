package models

import (
	"time"

	"gorm.io/gorm"
)

// Payment 支付订单模型
type Payment struct {
	ID              int        `gorm:"primaryKey;autoIncrement" json:"id"`
	PaymentNo       string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"payment_no"`
	BusinessType    string     `gorm:"type:varchar(50);not null;index:business_type_id" json:"business_type"`
	BusinessID      *string    `gorm:"type:varchar(36);index:business_type_id" json:"business_id"`
	CustomerID      string     `gorm:"type:varchar(36);not null;index" json:"customer_id"`
	CuUserID        string     `gorm:"type:varchar(36);not null;index" json:"cu_user_id"`
	Amount          float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency        string     `gorm:"type:varchar(3);default:'CNY'" json:"currency"`
	PaymentMethod   string     `gorm:"type:varchar(20);default:'alipay'" json:"payment_method"`
	PaymentProvider string     `gorm:"type:varchar(20);default:'alipay'" json:"payment_provider"`
	Status          string     `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	TradeNo         *string    `gorm:"type:varchar(64)" json:"trade_no"`
	PaymentTime     *time.Time `gorm:"type:datetime(3)" json:"payment_time"`
	ExpireTime      time.Time  `gorm:"type:datetime(3);not null" json:"expire_time"`
	PaymentURL      *string    `gorm:"type:text" json:"payment_url"`
	NotifyData      *string    `gorm:"type:json" json:"notify_data"`
	ExtraData       *string    `gorm:"type:json" json:"extra_data"`
	CreatedAt       time.Time  `gorm:"type:datetime(3);not null;index" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"type:datetime(3);not null" json:"updated_at"`
}

// TableName 指定表名
func (Payment) TableName() string {
	return "payments"
}

// BeforeCreate 创建前自动设置时间戳
func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = now
	}
	return nil
}

// BeforeUpdate 更新前自动设置时间戳
func (p *Payment) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

// PaymentCreateRequest 创建支付请求结构
type PaymentCreateRequest struct {
	BusinessType  string  `json:"business_type" binding:"required"`
	BusinessID    *string `json:"business_id"`
	CustomerID    string  `json:"customer_id" binding:"required"`
	CuUserID      string  `json:"cu_user_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,min=0.01"`
	Currency      string  `json:"currency"`
	PaymentMethod string  `json:"payment_method"`
}

// PaymentResponse 支付响应结构
type PaymentResponse struct {
	ID              int        `json:"id"`
	PaymentNo       string     `json:"payment_no"`
	BusinessType    string     `json:"business_type"`
	BusinessID      *string    `json:"business_id"`
	CustomerID      string     `json:"customer_id"`
	CuUserID        string     `json:"cu_user_id"`
	Amount          float64    `json:"amount"`
	Currency        string     `json:"currency"`
	PaymentMethod   string     `json:"payment_method"`
	PaymentProvider string     `json:"payment_provider"`
	Status          string     `json:"status"`
	TradeNo         *string    `json:"trade_no"`
	PaymentTime     *time.Time `json:"payment_time"`
	ExpireTime      time.Time  `json:"expire_time"`
	PaymentURL      *string    `json:"payment_url"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ToResponse 转换为响应结构
func (p *Payment) ToResponse() *PaymentResponse {
	return &PaymentResponse{
		ID:              p.ID,
		PaymentNo:       p.PaymentNo,
		BusinessType:    p.BusinessType,
		BusinessID:      p.BusinessID,
		CustomerID:      p.CustomerID,
		CuUserID:        p.CuUserID,
		Amount:          p.Amount,
		Currency:        p.Currency,
		PaymentMethod:   p.PaymentMethod,
		PaymentProvider: p.PaymentProvider,
		Status:          p.Status,
		TradeNo:         p.TradeNo,
		PaymentTime:     p.PaymentTime,
		ExpireTime:      p.ExpireTime,
		PaymentURL:      p.PaymentURL,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

// PaymentStatusResponse 支付状态响应结构
type PaymentStatusResponse struct {
	PaymentNo     string                `json:"payment_no"`
	Status        string                `json:"status"`
	Amount        float64               `json:"amount"`
	PaymentTime   *time.Time            `json:"payment_time"`
	TradeNo       *string               `json:"trade_no"`
	BusinessOrder *PaymentBusinessOrder `json:"business_order,omitempty"`
}

// PaymentBusinessOrder 支付关联的业务订单信息
type PaymentBusinessOrder struct {
	OrderNo           string  `json:"order_no"`
	AuthorizationCode *string `json:"authorization_code"`
	Status            string  `json:"status"`
}

// PaymentListResponse 支付列表响应结构
type PaymentListResponse struct {
	Payments   []*PaymentResponse `json:"payments"`
	TotalCount int64              `json:"total_count"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
}

// 业务类型常量
const (
	BusinessTypePackageOrder   = "package_order"   // 套餐购买订单
	BusinessTypeLicenseExtend  = "license_extend"  // 许可延期
	BusinessTypeServiceUpgrade = "service_upgrade" // 服务升级
)

// 支付状态常量
const (
	PaymentStatusPending   = "pending"   // 待支付
	PaymentStatusPaid      = "paid"      // 已支付
	PaymentStatusCancelled = "cancelled" // 已取消
	PaymentStatusExpired   = "expired"   // 已过期
	PaymentStatusFailed    = "failed"    // 支付失败
	PaymentStatusRefunded  = "refunded"  // 已退款
)

// 支付方式常量
const (
	PaymentMethodAlipay = "alipay" // 支付宝
	PaymentMethodWechat = "wechat" // 微信支付
	PaymentMethodBank   = "bank"   // 银行卡
)
