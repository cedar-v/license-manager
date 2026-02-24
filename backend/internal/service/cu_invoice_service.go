package service

import (
	"context"
	"fmt"
	"license-manager/internal/models"
	"license-manager/internal/repository"
	pkgcontext "license-manager/pkg/context"
	"license-manager/pkg/i18n"
	"strings"
	"time"

	"gorm.io/gorm"
)

// CuInvoiceService 客户端发票服务接口
type CuInvoiceService interface {
	// ========== 客户端发票操作 ==========
	CreateInvoice(ctx context.Context, cuUserID string, req *models.InvoiceCreateRequest) (*models.Invoice, error)
	GetInvoice(ctx context.Context, invoiceID, cuUserID string) (*models.Invoice, error)
	GetUserInvoices(ctx context.Context, cuUserID string, req *models.InvoiceListRequest) (*models.InvoiceListResponse, error)
	GetInvoiceSummary(ctx context.Context, cuUserID string) (*models.InvoiceSummaryResponse, error)
	DownloadInvoice(ctx context.Context, invoiceID, cuUserID string) (string, error)
	UpdateInvoice(ctx context.Context, invoiceID, cuUserID string, req *models.InvoiceUpdateRequest) (*models.Invoice, error)
}

type cuInvoiceService struct {
	repo      repository.CuInvoiceRepository
	orderRepo repository.CuOrderRepository
	db        *gorm.DB
}

func NewCuInvoiceService(repo repository.CuInvoiceRepository, orderRepo repository.CuOrderRepository, db *gorm.DB) CuInvoiceService {
	return &cuInvoiceService{
		repo:      repo,
		orderRepo: orderRepo,
		db:        db,
	}
}

// ========== 客户端发票操作实现 ==========

func (s *cuInvoiceService) CreateInvoice(ctx context.Context, cuUserID string, req *models.InvoiceCreateRequest) (*models.Invoice, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 验证订单是否存在且属于当前用户
	order, err := s.orderRepo.GetByID(req.OrderID)
	if err != nil {
		return nil, i18n.NewI18nError("700002", lang)
	}

	if order.CuUserID != cuUserID {
		return nil, i18n.NewI18nError("700002", lang)
	}

	// 验证订单状态必须为已支付
	if order.Status != "paid" {
		return nil, i18n.NewI18nError("700002", lang)
	}

	// 检查是否已存在发票
	exists, err := s.repo.CheckOrderInvoiceExists(req.OrderID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, i18n.NewI18nError("700003", lang)
	}

	// 验证发票类型和纳税人识别号
	if req.InvoiceType != models.InvoiceTypePersonal && (req.TaxpayerID == nil || *req.TaxpayerID == "") {
		return nil, i18n.NewI18nError("700004", lang)
	}

	// 生成发票申请号
	invoiceNo := s.generateInvoiceNo()

	// 创建发票记录
	now := time.Now()
	invoice := &models.Invoice{
		InvoiceNo:     invoiceNo,
		OrderID:       req.OrderID,
		OrderNo:       order.OrderNo,
		CustomerID:    order.CustomerID,
		CuUserID:      cuUserID,
		Amount:        order.TotalAmount,
		Status:        models.InvoiceStatusPending,
		InvoiceType:   req.InvoiceType,
		Title:         req.Title,
		TaxpayerID:    req.TaxpayerID,
		Content:       req.Content,
		ReceiverEmail: req.ReceiverEmail,
		Remark:        req.Remark,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err = s.repo.Create(invoice)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *cuInvoiceService) GetInvoice(ctx context.Context, invoiceID, cuUserID string) (*models.Invoice, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	invoice, err := s.repo.GetByID(invoiceID)
	if err != nil {
		return nil, i18n.NewI18nError("700001", lang)
	}

	// 验证权限：只能查看自己的发票
	if invoice.CuUserID != cuUserID {
		return nil, i18n.NewI18nError("700001", lang)
	}

	return invoice, nil
}

func (s *cuInvoiceService) GetUserInvoices(ctx context.Context, cuUserID string, req *models.InvoiceListRequest) (*models.InvoiceListResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 解析日期筛选
	var createdAtStart, createdAtEnd *time.Time
	if req.ApplyDate != "" {
		start, end, err := s.parseDateRange(req.ApplyDate)
		if err != nil {
			return nil, i18n.NewI18nError("900001", lang)
		}
		createdAtStart = start
		createdAtEnd = end
	}

	invoices, total, err := s.repo.GetByCuUserID(ctx, cuUserID, req, createdAtStart, createdAtEnd)
	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 转换为响应结构并填充多语言显示字段
	invoiceResponses := make([]*models.InvoiceResponse, len(invoices))
	for i, invoice := range invoices {
		invoiceResponses[i] = invoice.ToResponse()
		s.fillDisplayFields(invoiceResponses[i], lang)
	}

	return &models.InvoiceListResponse{
		Invoices:   invoiceResponses,
		TotalCount: total,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

func (s *cuInvoiceService) GetInvoiceSummary(ctx context.Context, cuUserID string) (*models.InvoiceSummaryResponse, error) {
	return s.repo.GetCuUserInvoiceSummary(ctx, cuUserID)
}

func (s *cuInvoiceService) DownloadInvoice(ctx context.Context, invoiceID, cuUserID string) (string, error) {
	invoice, err := s.GetInvoice(ctx, invoiceID, cuUserID)
	if err != nil {
		return "", err
	}

	lang := pkgcontext.GetLanguageFromContext(ctx)

	if invoice.Status != models.InvoiceStatusIssued {
		return "", i18n.NewI18nError("700005", lang)
	}

	if invoice.InvoiceFileURL == nil || *invoice.InvoiceFileURL == "" {
		return "", i18n.NewI18nError("700005", lang)
	}

	return *invoice.InvoiceFileURL, nil
}

func (s *cuInvoiceService) UpdateInvoice(ctx context.Context, invoiceID, cuUserID string, req *models.InvoiceUpdateRequest) (*models.Invoice, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 获取原发票
	invoice, err := s.repo.GetByID(invoiceID)
	if err != nil {
		return nil, i18n.NewI18nError("700001", lang)
	}

	// 权限校验
	if invoice.CuUserID != cuUserID {
		return nil, i18n.NewI18nError("700001", lang)
	}

	// 状态校验：只有被驳回的发票可以修改
	if invoice.Status != models.InvoiceStatusRejected {
		return nil, i18n.NewI18nError("700004", lang)
	}

	// 业务逻辑校验：发票类型和纳税人识别号
	if req.InvoiceType != models.InvoiceTypePersonal && (req.TaxpayerID == nil || *req.TaxpayerID == "") {
		return nil, i18n.NewI18nError("700004", lang)
	}

	// 更新字段
	invoice.InvoiceType = req.InvoiceType
	invoice.Title = req.Title
	invoice.TaxpayerID = req.TaxpayerID
	invoice.Content = req.Content
	invoice.ReceiverEmail = req.ReceiverEmail
	invoice.Remark = req.Remark

	// 重置状态和驳回信息
	invoice.Status = models.InvoiceStatusPending
	invoice.RejectReason = nil
	invoice.Suggestion = nil
	invoice.RejectedAt = nil
	invoice.RejectedBy = nil

	// 更新数据库
	err = s.repo.Update(invoice)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

// ========== 私有辅助方法 ==========

// fillDisplayFields 填充响应结构的多语言显示字段
func (s *cuInvoiceService) fillDisplayFields(response *models.InvoiceResponse, lang string) {
	response.StatusDisplay = i18n.GetEnumMessage("invoice_status", response.Status, lang)
	response.InvoiceTypeDisplay = i18n.GetEnumMessage("invoice_type", response.InvoiceType, lang)
}

// generateInvoiceNo 生成发票申请号
func (s *cuInvoiceService) generateInvoiceNo() string {
	now := time.Now()
	timestamp := now.UnixNano() / int64(time.Millisecond) % 1000000000 // 取纳秒后9位
	return fmt.Sprintf("INV%04d%02d%02d%09d",
		now.Year(), now.Month(), now.Day(), timestamp)
}

// parseDateRange 解析日期范围
func (s *cuInvoiceService) parseDateRange(dateStr string) (*time.Time, *time.Time, error) {
	parts := strings.Split(dateStr, ",")
	if len(parts) == 1 {
		start, _, err := s.parseFlexibleTime(strings.TrimSpace(dateStr))
		if err != nil {
			return nil, nil, err
		}
		end := start.AddDate(0, 0, 1)
		return &start, &end, nil
	}
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid date range")
	}

	startText := strings.TrimSpace(parts[0])
	endText := strings.TrimSpace(parts[1])
	if startText == "" || endText == "" {
		return nil, nil, fmt.Errorf("invalid date range")
	}

	start, _, err := s.parseFlexibleTime(startText)
	if err != nil {
		return nil, nil, err
	}
	end, endDateOnly, err := s.parseFlexibleTime(endText)
	if err != nil {
		return nil, nil, err
	}
	if endDateOnly {
		end = end.AddDate(0, 0, 1)
	}
	if end.Before(start) {
		start, end = end, start
	}

	return &start, &end, nil
}

func (s *cuInvoiceService) parseFlexibleTime(value string) (time.Time, bool, error) {
	if t, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return t, false, nil
	}
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t, false, nil
	}
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, false, err
	}
	return t, true, nil
}
