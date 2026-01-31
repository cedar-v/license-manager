package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PackageType 套餐类型
type PackageType string

const (
	PackageTypeTrial        PackageType = "trial"        // 试用版
	PackageTypeBasic        PackageType = "basic"        // 基础版
	PackageTypeProfessional PackageType = "professional" // 专业版
	PackageTypeCustom       PackageType = "custom"       // 定制版
)

// Package 套餐模型
type Package struct {
	ID                  string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name                string         `gorm:"type:varchar(100);not null" json:"name"`
	Type                string         `gorm:"type:varchar(20);not null" json:"type"`
	Price               float64        `gorm:"type:decimal(10,2);not null;default:0" json:"price"`
	PriceDescription    string         `gorm:"type:varchar(100);default:''" json:"price_description"`
	DurationDescription string         `gorm:"type:varchar(200);default:''" json:"duration_description"`
	Description         string         `gorm:"type:varchar(500);default:''" json:"description"`
	Features            string         `gorm:"type:text" json:"features"` // JSON数组格式
	Status              int            `gorm:"type:tinyint(1);not null;default:1" json:"status"`
	SortOrder           int            `gorm:"type:int;not null;default:0" json:"sort_order"`
	Remark              string         `gorm:"type:varchar(500);default:''" json:"remark"`
	CreatedAt           time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Package) TableName() string {
	return "packages"
}

// BeforeCreate 创建前自动设置ID和时间戳
func (p *Package) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
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
func (p *Package) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

// ToResponse 转换为响应结构
func (p *Package) ToResponse() *PackageResponse {
	return &PackageResponse{
		ID:                  p.ID,
		Name:                p.Name,
		Type:                p.Type,
		Price:               p.Price,
		PriceDescription:    p.PriceDescription,
		DurationDescription: p.DurationDescription,
		Description:         p.Description,
		Features:            p.Features,
		Status:              p.Status,
		SortOrder:           p.SortOrder,
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
	}
}

// PackageResponse 套餐响应结构
type PackageResponse struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Type                string    `json:"type"`
	Price               float64   `json:"price"`
	PriceDescription    string    `json:"price_description"`
	DurationDescription string    `json:"duration_description"`
	Description         string    `json:"description"`
	Features            string    `json:"features"`
	Status              int       `json:"status"`
	SortOrder           int       `json:"sort_order"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// PackageCreateRequest 创建套餐请求
type PackageCreateRequest struct {
	Name                string  `json:"name" binding:"required,min=1,max=100"`
	Type                string  `json:"type" binding:"required,oneof=trial basic professional custom"`
	Price               float64 `json:"price" binding:"gte=0"`
	PriceDescription    string  `json:"price_description" binding:"max=100"`
	DurationDescription string  `json:"duration_description" binding:"max=200"`
	Description         string  `json:"description" binding:"max=500"`
	Features            string  `json:"features"` // JSON格式
	Status              int     `json:"status" binding:"oneof=0 1"`
	SortOrder           int     `json:"sort_order"`
	Remark              string  `json:"remark" binding:"max=500"`
}

// PackageUpdateRequest 更新套餐请求
type PackageUpdateRequest struct {
	Name                string  `json:"name" binding:"omitempty,min=1,max=100"`
	Type                string  `json:"type" binding:"omitempty,oneof=trial basic professional custom"`
	Price               float64 `json:"price" binding:"omitempty,gte=0"`
	PriceDescription    string  `json:"price_description" binding:"omitempty,max=100"`
	DurationDescription string  `json:"duration_description" binding:"omitempty,max=200"`
	Description         string  `json:"description" binding:"omitempty,max=500"`
	Features            string  `json:"features"` // JSON格式
	Status              *int    `json:"status" binding:"omitempty,oneof=0 1"`
	SortOrder           *int    `json:"sort_order"`
	Remark              string  `json:"remark" binding:"omitempty,max=500"`
}

// PackageListRequest 套餐列表请求
type PackageListRequest struct {
	Type   string `form:"type"`   // 套餐类型筛选
	Status *int   `form:"status"` // 状态筛选
}

// PackageListResponse 套餐列表响应
type PackageListResponse struct {
	Packages   []*PackageResponse `json:"packages"`
	TotalCount int64              `json:"total_count"`
}

// CuPackageResponse 用户端套餐响应（兼容原有结构）
type CuPackageResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
	MaxDevices  int     `json:"max_devices"` // 从features中解析
	Description string  `json:"description"`
	Features    string  `json:"features"`
	Details     string  `json:"details"`
}
