package service

import (
	"context"
	"license-manager/internal/models"
	"license-manager/internal/repository"
	pkgcontext "license-manager/pkg/context"
	"license-manager/pkg/i18n"
	"time"

	"gorm.io/gorm"
)

// LeadService 线索服务接口
type LeadService interface {
	CreateLead(ctx context.Context, req *models.LeadCreateRequest) (*models.Lead, error)
	UpdateLead(ctx context.Context, id string, req *models.LeadUpdateRequest) (*models.Lead, error)
	DeleteLead(ctx context.Context, id string) error
	GetLead(ctx context.Context, id string) (*models.Lead, error)
	GetLeadList(ctx context.Context, req *models.LeadListRequest) (*models.LeadListResponse, error)
}

type leadService struct {
	repo repository.LeadRepository
	db   *gorm.DB
}

// NewLeadService 创建线索服务
func NewLeadService(repo repository.LeadRepository, db *gorm.DB) LeadService {
	return &leadService{
		repo: repo,
		db:   db,
	}
}

func (s *leadService) CreateLead(ctx context.Context, req *models.LeadCreateRequest) (*models.Lead, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	lead := &models.Lead{
		CompanyName:  req.CompanyName,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		ContactEmail: req.ContactEmail,
		Requirement:  req.Requirement,
		ExtraInfo:    req.ExtraInfo,
		Status:       string(models.LeadStatusPending),
	}

	if err := s.repo.Create(lead); err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	return lead, nil
}

func (s *leadService) UpdateLead(ctx context.Context, id string, req *models.LeadUpdateRequest) (*models.Lead, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	lead, err := s.repo.GetByID(id)
	if err != nil {
		return nil, i18n.NewI18nError("900002", lang)
	}

	// 更新字段
	if req.CompanyName != "" {
		lead.CompanyName = req.CompanyName
	}
	if req.ContactName != "" {
		lead.ContactName = req.ContactName
	}
	if req.ContactPhone != "" {
		lead.ContactPhone = req.ContactPhone
	}
	if req.ContactEmail != "" {
		lead.ContactEmail = req.ContactEmail
	}
	if req.Requirement != "" {
		lead.Requirement = req.Requirement
	}
	if req.ExtraInfo != "" {
		lead.ExtraInfo = req.ExtraInfo
	}
	if req.Status != "" {
		lead.Status = req.Status
	}
	lead.FollowUpDate = req.FollowUpDate
	if req.FollowUpRecord != "" {
		lead.FollowUpRecord = req.FollowUpRecord
	}
	if req.InternalNote != "" {
		lead.InternalNote = req.InternalNote
	}
	lead.UpdatedAt = time.Now()

	if err := s.repo.Update(lead); err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	return lead, nil
}

func (s *leadService) DeleteLead(ctx context.Context, id string) error {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 检查是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return i18n.NewI18nError("900002", lang)
	}

	if err := s.repo.Delete(id); err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	return nil
}

func (s *leadService) GetLead(ctx context.Context, id string) (*models.Lead, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	lead, err := s.repo.GetByID(id)
	if err != nil {
		return nil, i18n.NewI18nError("900002", lang)
	}

	return lead, nil
}

func (s *leadService) GetLeadList(ctx context.Context, req *models.LeadListRequest) (*models.LeadListResponse, error) {
	leads, total, err := s.repo.GetList(ctx, req)
	if err != nil {
		lang := pkgcontext.GetLanguageFromContext(ctx)
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 转换为响应
	leadResponses := make([]*models.LeadResponse, len(leads))
	for i, lead := range leads {
		leadResponses[i] = lead.ToResponse()
	}

	page := 1
	pageSize := 10
	if req != nil {
		if req.Page > 0 {
			page = req.Page
		}
		if req.PageSize > 0 {
			pageSize = req.PageSize
		}
	}

	return &models.LeadListResponse{
		Leads:      leadResponses,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}
