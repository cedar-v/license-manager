package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LeadStatus 线索状态
type LeadStatus string

const (
	LeadStatusPending   LeadStatus = "pending"   // 待联系
	LeadStatusContacted LeadStatus = "contacted" // 已联系
	LeadStatusConverted LeadStatus = "converted" // 已成交
	LeadStatusInvalid   LeadStatus = "invalid"   // 已失效
)

// Lead 企业线索
type Lead struct {
	ID             string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	LeadNo         string     `gorm:"type:varchar(50);not null;uniqueIndex" json:"lead_no"`
	CompanyName    string     `gorm:"type:varchar(200);not null" json:"company_name"`
	ContactName    string     `gorm:"type:varchar(100);not null" json:"contact_name"`
	ContactPhone   string     `gorm:"type:varchar(20);not null" json:"contact_phone"`
	ContactEmail   string     `gorm:"type:varchar(100)" json:"contact_email"`
	Requirement    string     `gorm:"type:text;not null" json:"requirement"`
	ExtraInfo      string     `gorm:"type:text" json:"extra_info"`
	Status         string     `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	FollowUpDate   *time.Time `gorm:"" json:"follow_up_date"`
	FollowUpRecord string     `gorm:"type:text" json:"follow_up_record"`
	InternalNote   string     `gorm:"type:text" json:"internal_note"`
	CreatedAt      time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"not null" json:"updated_at"`
}

// TableName 指定表名
func (Lead) TableName() string {
	return "leads"
}

// BeforeCreate 创建前自动设置ID和时间戳
func (l *Lead) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	now := time.Now()
	if l.CreatedAt.IsZero() {
		l.CreatedAt = now
	}
	if l.UpdatedAt.IsZero() {
		l.UpdatedAt = now
	}
	// 生成线索编号
	if l.LeadNo == "" {
		l.LeadNo = generateLeadNo()
	}
	return nil
}

// generateLeadNo 生成线索编号
func generateLeadNo() string {
	return "LEAD" + time.Now().Format("20060102150405") + randomString(6)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	now := time.Now().UnixNano()
	for i := range result {
		result[i] = charset[(now+int64(i)*17)%int64(len(charset))]
	}
	return string(result)
}

// ToResponse 转换为响应结构
func (l *Lead) ToResponse() *LeadResponse {
	return &LeadResponse{
		ID:             l.ID,
		LeadNo:         l.LeadNo,
		CompanyName:    l.CompanyName,
		ContactName:    l.ContactName,
		ContactPhone:   l.ContactPhone,
		ContactEmail:   l.ContactEmail,
		Requirement:    l.Requirement,
		ExtraInfo:      l.ExtraInfo,
		Status:         l.Status,
		FollowUpDate:   l.FollowUpDate,
		FollowUpRecord: l.FollowUpRecord,
		InternalNote:   l.InternalNote,
		CreatedAt:      l.CreatedAt,
		UpdatedAt:      l.UpdatedAt,
	}
}

// LeadResponse 线索响应结构
type LeadResponse struct {
	ID             string     `json:"id"`
	LeadNo         string     `json:"lead_no"`
	CompanyName    string     `json:"company_name"`
	ContactName    string     `json:"contact_name"`
	ContactPhone   string     `json:"contact_phone"`
	ContactEmail   string     `json:"contact_email"`
	Requirement    string     `json:"requirement"`
	ExtraInfo      string     `json:"extra_info"`
	Status         string     `json:"status"`
	FollowUpDate   *time.Time `json:"follow_up_date"`
	FollowUpRecord string     `json:"follow_up_record"`
	InternalNote   string     `json:"internal_note"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// LeadCreateRequest 创建线索请求
type LeadCreateRequest struct {
	CompanyName  string `json:"company_name" binding:"required,min=1,max=200"`
	ContactName  string `json:"contact_name" binding:"required,min=1,max=100"`
	ContactPhone string `json:"contact_phone" binding:"required,min=1,max=20"`
	ContactEmail string `json:"contact_email" binding:"omitempty,email,max=100"`
	Requirement  string `json:"requirement" binding:"required,min=1"`
	ExtraInfo    string `json:"extra_info"`
}

// LeadUpdateRequest 更新线索请求
type LeadUpdateRequest struct {
	CompanyName    string     `json:"company_name" binding:"omitempty,min=1,max=200"`
	ContactName    string     `json:"contact_name" binding:"omitempty,min=1,max=100"`
	ContactPhone   string     `json:"contact_phone" binding:"omitempty,min=1,max=20"`
	ContactEmail   string     `json:"contact_email" binding:"omitempty,email,max=100"`
	Requirement    string     `json:"requirement" binding:"omitempty,min=1"`
	ExtraInfo      string     `json:"extra_info"`
	Status         string     `json:"status" binding:"omitempty,oneof=pending contacted converted invalid"`
	FollowUpDate   *time.Time `json:"follow_up_date"`
	FollowUpRecord string     `json:"follow_up_record"`
	InternalNote   string     `json:"internal_note"`
}

// LeadListRequest 线索列表请求
type LeadListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search"` // 关键词检索
	Status   string `form:"status"` // 状态筛选
}

// LeadListResponse 线索列表响应
type LeadListResponse struct {
	Leads      []*LeadResponse `json:"leads"`
	TotalCount int64           `json:"total_count"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
}
