package repository

import (
	"context"
	"license-manager/internal/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

// AdminInvoiceRepository 管理端发票仓储接口
type AdminInvoiceRepository interface {
	// ========== 管理端专用方法 ==========
	GetByCustomerID(customerID string, offset, limit int) ([]*models.Invoice, int64, error)
	GetAdminList(ctx context.Context, req *models.InvoiceListRequest, createdAtStart, createdAtEnd *time.Time) ([]*models.Invoice, int64, error)
	GetAdminInvoiceSummary(ctx context.Context) (*models.InvoiceSummaryResponse, error)

	// ========== 通用工具方法 ==========
	GetByID(id string) (*models.Invoice, error)
	GetByInvoiceNo(invoiceNo string) (*models.Invoice, error)
	GetByOrderID(orderID string) (*models.Invoice, error)
	Update(invoice *models.Invoice) error
	Delete(id string) error
	GetByDownloadToken(token string) (*models.Invoice, error)
	UpdateDownloadToken(id string, token *string) error
}

type adminInvoiceRepository struct {
	db *gorm.DB
}

func NewAdminInvoiceRepository(db *gorm.DB) AdminInvoiceRepository {
	return &adminInvoiceRepository{db: db}
}

func (r *adminInvoiceRepository) GetByID(id string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *adminInvoiceRepository) GetByInvoiceNo(invoiceNo string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.db.Where("invoice_no = ? AND deleted_at IS NULL", invoiceNo).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *adminInvoiceRepository) GetByOrderID(orderID string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.db.Where("order_id = ? AND deleted_at IS NULL", orderID).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *adminInvoiceRepository) GetByCustomerID(customerID string, offset, limit int) ([]*models.Invoice, int64, error) {
	var invoices []*models.Invoice
	var total int64

	query := r.db.Model(&models.Invoice{}).Where("customer_id = ? AND deleted_at IS NULL", customerID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&invoices).Error
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func (r *adminInvoiceRepository) GetAdminList(ctx context.Context, req *models.InvoiceListRequest, createdAtStart, createdAtEnd *time.Time) ([]*models.Invoice, int64, error) {
	var invoices []*models.Invoice
	var total int64

	// 默认值
	page := 1
	pageSize := 10
	search := ""
	status := ""
	customerID := ""
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
		customerID = strings.TrimSpace(req.CustomerID)
	}
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&models.Invoice{}).Where("deleted_at IS NULL")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
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

func (r *adminInvoiceRepository) Update(invoice *models.Invoice) error {
	return r.db.Save(invoice).Error
}

func (r *adminInvoiceRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Invoice{}).Error
}

func (r *adminInvoiceRepository) GetAdminInvoiceSummary(ctx context.Context) (*models.InvoiceSummaryResponse, error) {
	// 查询发票总数
	var totalCount int64
	err := r.db.Model(&models.Invoice{}).
		Where("deleted_at IS NULL").
		Count(&totalCount).Error
	if err != nil {
		return nil, err
	}

	// 查询待处理发票数
	var pendingCount int64
	err = r.db.Model(&models.Invoice{}).
		Where("status = ? AND deleted_at IS NULL", models.InvoiceStatusPending).
		Count(&pendingCount).Error
	if err != nil {
		return nil, err
	}

	// 查询已开票数
	var issuedCount int64
	err = r.db.Model(&models.Invoice{}).
		Where("status = ? AND deleted_at IS NULL", models.InvoiceStatusIssued).
		Count(&issuedCount).Error
	if err != nil {
		return nil, err
	}

	// 查询已驳回数
	var rejectedCount int64
	err = r.db.Model(&models.Invoice{}).
		Where("status = ? AND deleted_at IS NULL", models.InvoiceStatusRejected).
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

func (r *adminInvoiceRepository) GetByDownloadToken(token string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.db.Where("download_token = ? AND deleted_at IS NULL", token).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *adminInvoiceRepository) UpdateDownloadToken(id string, token *string) error {
	return r.db.Model(&models.Invoice{}).Where("id = ?", id).Update("download_token", token).Error
}
