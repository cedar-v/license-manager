package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"license-manager/internal/models"
	"license-manager/internal/repository"
	pkgcontext "license-manager/pkg/context"
	"license-manager/pkg/i18n"
	"time"

	"gorm.io/gorm"
)

// AdminInvoiceService 管理端发票服务接口
type AdminInvoiceService interface {
	// ========== 管理端发票操作 ==========
	RejectInvoice(ctx context.Context, invoiceID string, adminID string, req *models.InvoiceRejectRequest) (*models.Invoice, error)
	IssueInvoice(ctx context.Context, invoiceID string, adminID string, req *models.InvoiceIssueRequest) (*models.Invoice, error)
	GetAdminInvoices(ctx context.Context, req *models.InvoiceListRequest) (*models.InvoiceListResponse, error)
	GetAdminInvoiceSummary(ctx context.Context) (*models.InvoiceSummaryResponse, error)
	GetInvoiceDetail(ctx context.Context, invoiceID string, isAdmin bool, userID string) (*models.InvoiceDetailResponse, error)
	GenerateDownloadToken(ctx context.Context, invoiceID string) (string, error)

	// ========== 公共发票操作 ==========
	DownloadByToken(ctx context.Context, token string) (string, error)
}

type adminInvoiceService struct {
	repo       repository.AdminInvoiceRepository
	orderRepo  repository.CuOrderRepository
	userRepo   repository.UserRepository
	cuUserRepo repository.CuUserRepository
	db         *gorm.DB
}

func NewAdminInvoiceService(repo repository.AdminInvoiceRepository, orderRepo repository.CuOrderRepository, userRepo repository.UserRepository, cuUserRepo repository.CuUserRepository, db *gorm.DB) AdminInvoiceService {
	return &adminInvoiceService{
		repo:       repo,
		orderRepo:  orderRepo,
		userRepo:   userRepo,
		cuUserRepo: cuUserRepo,
		db:         db,
	}
}

// ========== 管理端发票操作实现 ==========

func (s *adminInvoiceService) RejectInvoice(ctx context.Context, invoiceID string, adminID string, req *models.InvoiceRejectRequest) (*models.Invoice, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	invoice, err := s.repo.GetByID(invoiceID)
	if err != nil {
		return nil, i18n.NewI18nError("700001", lang)
	}

	// 验证状态
	if invoice.Status != models.InvoiceStatusPending {
		return nil, i18n.NewI18nError("700004", lang)
	}

	now := time.Now()
	invoice.Status = models.InvoiceStatusRejected
	invoice.RejectReason = &req.RejectReason
	invoice.Suggestion = &req.Suggestion
	invoice.RejectedAt = &now
	invoice.RejectedBy = &adminID
	invoice.UpdatedAt = now

	err = s.repo.Update(invoice)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *adminInvoiceService) IssueInvoice(ctx context.Context, invoiceID string, adminID string, req *models.InvoiceIssueRequest) (*models.Invoice, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	invoice, err := s.repo.GetByID(invoiceID)
	if err != nil {
		return nil, i18n.NewI18nError("700001", lang)
	}

	// 验证状态：允许从pending或rejected状态开票
	if invoice.Status != models.InvoiceStatusPending && invoice.Status != models.InvoiceStatusRejected {
		return nil, i18n.NewI18nError("700004", lang)
	}

	// 解析开票时间
	issuedAt, err := time.Parse(time.RFC3339, req.IssuedAt)
	if err != nil {
		return nil, i18n.NewI18nError("700004", lang)
	}

	now := time.Now()
	invoice.Status = models.InvoiceStatusIssued
	invoice.InvoiceFileURL = &req.InvoiceFileURL
	invoice.IssuedAt = &issuedAt
	invoice.UploadedAt = &now
	invoice.UploadedBy = &adminID
	invoice.UpdatedAt = now

	// 清除驳回信息（如果是从rejected状态开票）
	if invoice.Status == models.InvoiceStatusRejected {
		invoice.RejectReason = nil
		invoice.Suggestion = nil
		invoice.RejectedAt = nil
		invoice.RejectedBy = nil
	}

	err = s.repo.Update(invoice)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *adminInvoiceService) GetAdminInvoices(ctx context.Context, req *models.InvoiceListRequest) (*models.InvoiceListResponse, error) {
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

	invoices, total, err := s.repo.GetAdminList(ctx, req, createdAtStart, createdAtEnd)
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

func (s *adminInvoiceService) GetAdminInvoiceSummary(ctx context.Context) (*models.InvoiceSummaryResponse, error) {
	return s.repo.GetAdminInvoiceSummary(ctx)
}

func (s *adminInvoiceService) GetInvoiceDetail(ctx context.Context, invoiceID string, isAdmin bool, userID string) (*models.InvoiceDetailResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	invoice, err := s.repo.GetByID(invoiceID)
	if err != nil {
		return nil, i18n.NewI18nError("700001", lang)
	}

	// 权限验证
	if !isAdmin && invoice.CuUserID != userID {
		return nil, i18n.NewI18nError("700001", lang)
	}

	response := &models.InvoiceDetailResponse{
		InvoiceResponse: invoice.ToResponse(),
	}

	// 填充多语言显示字段
	s.fillDetailDisplayFields(response, lang)

	// 获取关联信息
	if isAdmin {
		// 管理端需要更多关联信息
		order, err := s.orderRepo.GetByID(invoice.OrderID)
		if err == nil {
			response.OrderPackageName = order.PackageName
		}

		cuUser, err := s.cuUserRepo.GetByID(invoice.CuUserID)
		if err == nil && cuUser.RealName != nil {
			response.ApplicantName = *cuUser.RealName
			response.ApplicantPhone = cuUser.Phone
		}

		if invoice.UploadedBy != nil {
			user, err := s.userRepo.GetUserByID(ctx, *invoice.UploadedBy)
			if err == nil {
				response.UploaderName = user.FullName
			}
		}

		if invoice.RejectedBy != nil {
			user, err := s.userRepo.GetUserByID(ctx, *invoice.RejectedBy)
			if err == nil {
				response.RejecterName = user.FullName
			}
		}
	}

	return response, nil
}

func (s *adminInvoiceService) GenerateDownloadToken(ctx context.Context, invoiceID string) (string, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	invoice, err := s.repo.GetByID(invoiceID)
	if err != nil {
		return "", i18n.NewI18nError("700001", lang)
	}

	if invoice.Status != models.InvoiceStatusIssued {
		return "", i18n.NewI18nError("700005", lang)
	}

	if invoice.InvoiceFileURL == nil || *invoice.InvoiceFileURL == "" {
		return "", i18n.NewI18nError("700005", lang)
	}

	// 生成随机token
	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	// 更新token
	err = s.repo.UpdateDownloadToken(invoiceID, &token)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ========== 公共发票操作实现 ==========

func (s *adminInvoiceService) DownloadByToken(ctx context.Context, token string) (string, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	invoice, err := s.repo.GetByDownloadToken(token)
	if err != nil {
		return "", i18n.NewI18nError("700001", lang)
	}

	if invoice.Status != models.InvoiceStatusIssued {
		return "", i18n.NewI18nError("700005", lang)
	}

	if invoice.InvoiceFileURL == nil || *invoice.InvoiceFileURL == "" {
		return "", i18n.NewI18nError("700005", lang)
	}

	return *invoice.InvoiceFileURL, nil
}

// ========== 私有辅助方法 ==========

// fillDisplayFields 填充响应结构的多语言显示字段
func (s *adminInvoiceService) fillDisplayFields(response *models.InvoiceResponse, lang string) {
	response.StatusDisplay = i18n.GetEnumMessage("invoice_status", response.Status, lang)
	response.InvoiceTypeDisplay = i18n.GetEnumMessage("invoice_type", response.InvoiceType, lang)
}

// fillDetailDisplayFields 填充详情响应结构的多语言显示字段
func (s *adminInvoiceService) fillDetailDisplayFields(response *models.InvoiceDetailResponse, lang string) {
	response.StatusDisplay = i18n.GetEnumMessage("invoice_status", response.Status, lang)
	response.InvoiceTypeDisplay = i18n.GetEnumMessage("invoice_type", response.InvoiceType, lang)
}

// parseDateRange 解析日期范围
func (s *adminInvoiceService) parseDateRange(dateStr string) (*time.Time, *time.Time, error) {
	// 假设输入格式为 YYYY-MM-DD
	start, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, nil, err
	}

	end := start.AddDate(0, 0, 1) // 下一天

	return &start, &end, nil
}
