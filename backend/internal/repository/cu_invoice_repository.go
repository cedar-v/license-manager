package repository

import (
	"context"
	"license-manager/internal/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

// CuInvoiceRepository 客户端发票仓储接口
type CuInvoiceRepository interface {
	// ========== 客户端专用方法 ==========
	GetByCuUserID(ctx context.Context, cuUserID string, req *models.InvoiceListRequest, createdAtStart, createdAtEnd *time.Time) ([]*models.Invoice, int64, error)
	GetCuUserInvoiceSummary(ctx context.Context, cuUserID string) (*models.InvoiceSummaryResponse, error)

	// ========== 通用工具方法 ==========
	Create(invoice *models.Invoice) error
	Update(invoice *models.Invoice) error
	GetByID(id string) (*models.Invoice, error)
	CheckOrderInvoiceExists(orderID string) (bool, error)
}

type cuInvoiceRepository struct {
	db *gorm.DB
}

func NewCuInvoiceRepository(db *gorm.DB) CuInvoiceRepository {
	return &cuInvoiceRepository{db: db}
}

func (r *cuInvoiceRepository) Create(invoice *models.Invoice) error {
	return r.db.Create(invoice).Error
}

func (r *cuInvoiceRepository) Update(invoice *models.Invoice) error {
	// 使用 Select 指定更新字段，确保可以清空某些字段（虽然此处目前主要是业务更新和状态重置）
	// 或者直接使用 Save 更新全量字段
	return r.db.Save(invoice).Error
}

func (r *cuInvoiceRepository) GetByID(id string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *cuInvoiceRepository) GetByCuUserID(ctx context.Context, cuUserID string, req *models.InvoiceListRequest, createdAtStart, createdAtEnd *time.Time) ([]*models.Invoice, int64, error) {
	var invoices []*models.Invoice
	var total int64

	// 默认值
	page := 1
	pageSize := 10
	search := ""
	status := ""
	if req != nil {
		if req.Page > 0 {
			page = req.Page
		}
		if req.PageSize > 0 {
			pageSize = req.PageSize
		}
		if pageSize > 100 {
			pageSize = 100
		}
		search = strings.TrimSpace(req.Search)
		status = strings.TrimSpace(req.Status)
	}
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&models.Invoice{}).Where("cu_user_id = ? AND deleted_at IS NULL", cuUserID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("(invoice_no LIKE ? OR order_no LIKE ?)", like, like)
	}

	if createdAtStart != nil && createdAtEnd != nil {
		query = query.Where("created_at >= ? AND created_at < ?", *createdAtStart, *createdAtEnd)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&invoices).Error
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func (r *cuInvoiceRepository) GetCuUserInvoiceSummary(ctx context.Context, cuUserID string) (*models.InvoiceSummaryResponse, error) {
	// 查询发票总数
	var totalCount int64
	err := r.db.Model(&models.Invoice{}).
		Where("cu_user_id = ? AND deleted_at IS NULL", cuUserID).
		Count(&totalCount).Error
	if err != nil {
		return nil, err
	}

	// 查询待处理发票数
	var pendingCount int64
	err = r.db.Model(&models.Invoice{}).
		Where("cu_user_id = ? AND status = ? AND deleted_at IS NULL", cuUserID, models.InvoiceStatusPending).
		Count(&pendingCount).Error
	if err != nil {
		return nil, err
	}

	// 查询已开票数
	var issuedCount int64
	err = r.db.Model(&models.Invoice{}).
		Where("cu_user_id = ? AND status = ? AND deleted_at IS NULL", cuUserID, models.InvoiceStatusIssued).
		Count(&issuedCount).Error
	if err != nil {
		return nil, err
	}

	// 查询已驳回数
	var rejectedCount int64
	err = r.db.Model(&models.Invoice{}).
		Where("cu_user_id = ? AND status = ? AND deleted_at IS NULL", cuUserID, models.InvoiceStatusRejected).
		Count(&rejectedCount).Error
	if err != nil {
		return nil, err
	}

	return &models.InvoiceSummaryResponse{
		TotalCount:    totalCount,
		PendingCount:  pendingCount,
		IssuedCount:   issuedCount,
		RejectedCount: rejectedCount,
	}, nil
}

func (r *cuInvoiceRepository) CheckOrderInvoiceExists(orderID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Invoice{}).
		Where("order_id = ? AND deleted_at IS NULL", orderID).
		Count(&count).Error
	return count > 0, err
}
