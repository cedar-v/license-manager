package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ========== 共享发票数据模型 ==========

// Invoice 发票模型
type Invoice struct {
	ID             string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	InvoiceNo      string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"invoice_no"`
	OrderID        string         `gorm:"type:varchar(36);not null" json:"order_id"`
	OrderNo        string         `gorm:"type:varchar(50);not null;index" json:"order_no"`
	CustomerID     string         `gorm:"type:varchar(36);not null;index" json:"customer_id"`
	CuUserID       string         `gorm:"type:varchar(36);not null;index" json:"cu_user_id"`
	Amount         float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status         string         `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	InvoiceType    string         `gorm:"type:varchar(20);not null" json:"invoice_type"`
	Title          string         `gorm:"type:varchar(200);not null" json:"title"`
	TaxpayerID     *string        `gorm:"type:varchar(50)" json:"taxpayer_id"`
	Content        string         `gorm:"type:varchar(200);not null" json:"content"`
	ReceiverEmail  string         `gorm:"type:varchar(255);not null" json:"receiver_email"`
	Remark         *string        `gorm:"type:varchar(1000)" json:"remark"`
	InvoiceFileURL *string        `gorm:"type:varchar(500)" json:"invoice_file_url"`
	UploadedAt     *time.Time     `gorm:"" json:"uploaded_at"`
	UploadedBy     *string        `gorm:"type:varchar(36)" json:"uploaded_by"`
	IssuedAt       *time.Time     `gorm:"" json:"issued_at"`
	RejectReason   *string        `gorm:"type:varchar(500)" json:"reject_reason"`
	Suggestion     *string        `gorm:"type:varchar(500)" json:"suggestion"`
	RejectedAt     *time.Time     `gorm:"" json:"rejected_at"`
	RejectedBy     *string        `gorm:"type:varchar(36)" json:"rejected_by"`
	DownloadToken  *string        `gorm:"type:varchar(64)" json:"download_token"`
	CreatedAt      time.Time      `gorm:"not null;index" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Invoice) TableName() string {
	return "invoices"
}

// BeforeCreate 创建前自动设置时间戳和ID
func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	// 生成UUID作为主键ID
	if i.ID == "" {
		i.ID = uuid.New().String()
	}

	// 设置时间戳
	now := time.Now()
	if i.CreatedAt.IsZero() {
		i.CreatedAt = now
	}
	if i.UpdatedAt.IsZero() {
		i.UpdatedAt = now
	}

	return nil
}

// BeforeUpdate 更新前自动设置时间戳
func (i *Invoice) BeforeUpdate(tx *gorm.DB) error {
	i.UpdatedAt = time.Now()
	return nil
}

// ========== 发票枚举常量 ==========

// 发票状态枚举
const (
	InvoiceStatusPending  = "pending"  // 待处理
	InvoiceStatusIssued   = "issued"   // 已开票
	InvoiceStatusRejected = "rejected" // 已驳回
)

// 发票类型枚举
const (
	InvoiceTypePersonal   = "personal"    // 个人
	InvoiceTypeEnterprise = "enterprise"  // 企业普票
	InvoiceTypeVATSpecial = "vat_special" // 增值税专用发票
)

// ========== 共享响应结构体 ==========

// InvoiceResponse 发票响应结构
type InvoiceResponse struct {
	ID                 string     `json:"id"`                   // 发票ID
	InvoiceNo          string     `json:"invoice_no"`           // 发票申请号，格式如INV202601210123456789
	OrderID            string     `json:"order_id"`             // 关联的订单ID
	OrderNo            string     `json:"order_no"`             // 关联的订单号
	CustomerID         string     `json:"customer_id"`          // 客户ID
	CuUserID           string     `json:"cu_user_id"`           // 申请人客户用户ID
	Amount             float64    `json:"amount"`               // 发票金额，单位元，保留2位小数
	Status             string     `json:"status"`               // 发票状态，pending-待处理/issued-已开票/rejected-已驳回
	StatusDisplay      string     `json:"status_display"`       // 发票状态显示文本，根据当前语言显示
	InvoiceType        string     `json:"invoice_type"`         // 发票类型，personal-个人/enterprise-企业普票/vat_special-增值税专用发票
	InvoiceTypeDisplay string     `json:"invoice_type_display"` // 发票类型显示文本，根据当前语言显示
	Title              string     `json:"title"`                // 发票抬头
	TaxpayerID         *string    `json:"taxpayer_id"`          // 纳税人识别号，企业发票填写
	Content            string     `json:"content"`              // 开票内容
	ReceiverEmail      string     `json:"receiver_email"`       // 收票邮箱
	Remark             *string    `json:"remark"`               // 备注信息
	InvoiceFileURL     *string    `json:"invoice_file_url"`     // 发票文件下载URL，已开票后填写
	UploadedAt         *time.Time `json:"uploaded_at"`          // 发票文件上传时间
	UploadedBy         *string    `json:"uploaded_by"`          // 发票文件上传人ID（管理员）
	IssuedAt           *time.Time `json:"issued_at"`            // 发票开票完成时间
	RejectReason       *string    `json:"reject_reason"`        // 驳回原因，已驳回时填写
	Suggestion         *string    `json:"suggestion"`           // 修改建议，已驳回时填写
	RejectedAt         *time.Time `json:"rejected_at"`          // 驳回时间
	RejectedBy         *string    `json:"rejected_by"`          // 驳回人ID（管理员）
	DownloadToken      *string    `json:"download_token"`       // 下载令牌，用于邮件链接下载
	CreatedAt          time.Time  `json:"created_at"`           // 创建时间
	UpdatedAt          time.Time  `json:"updated_at"`           // 更新时间
}

// ToResponse 转换为响应结构（不包含多语言显示字段）
func (i *Invoice) ToResponse() *InvoiceResponse {
	return &InvoiceResponse{
		ID:             i.ID,
		InvoiceNo:      i.InvoiceNo,
		OrderID:        i.OrderID,
		OrderNo:        i.OrderNo,
		CustomerID:     i.CustomerID,
		CuUserID:       i.CuUserID,
		Amount:         i.Amount,
		Status:         i.Status,
		InvoiceType:    i.InvoiceType,
		Title:          i.Title,
		TaxpayerID:     i.TaxpayerID,
		Content:        i.Content,
		ReceiverEmail:  i.ReceiverEmail,
		Remark:         i.Remark,
		InvoiceFileURL: i.InvoiceFileURL,
		UploadedAt:     i.UploadedAt,
		UploadedBy:     i.UploadedBy,
		IssuedAt:       i.IssuedAt,
		RejectReason:   i.RejectReason,
		Suggestion:     i.Suggestion,
		RejectedAt:     i.RejectedAt,
		RejectedBy:     i.RejectedBy,
		DownloadToken:  i.DownloadToken,
		CreatedAt:      i.CreatedAt,
		UpdatedAt:      i.UpdatedAt,
	}
}
