package repository

import (
	"context"
	"license-manager/internal/models"

	"gorm.io/gorm"
)

// LeadRepository 线索仓储接口
type LeadRepository interface {
	Create(lead *models.Lead) error
	Update(lead *models.Lead) error
	Delete(id string) error
	GetByID(id string) (*models.Lead, error)
	GetList(ctx context.Context, req *models.LeadListRequest) ([]*models.Lead, int64, error)
}

type leadRepository struct {
	db *gorm.DB
}

// NewLeadRepository 创建线索仓储
func NewLeadRepository(db *gorm.DB) LeadRepository {
	return &leadRepository{db: db}
}

func (r *leadRepository) Create(lead *models.Lead) error {
	return r.db.Create(lead).Error
}

func (r *leadRepository) Update(lead *models.Lead) error {
	return r.db.Save(lead).Error
}

func (r *leadRepository) Delete(id string) error {
	return r.db.Unscoped().Where("id = ?", id).Delete(&models.Lead{}).Error
}

func (r *leadRepository) GetByID(id string) (*models.Lead, error) {
	var lead models.Lead
	err := r.db.Where("id = ?", id).First(&lead).Error
	if err != nil {
		return nil, err
	}
	return &lead, nil
}

func (r *leadRepository) GetList(ctx context.Context, req *models.LeadListRequest) ([]*models.Lead, int64, error) {
	var leads []*models.Lead
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
		search = req.Search
		status = req.Status
	}
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&models.Lead{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("(company_name LIKE ? OR contact_name LIKE ? OR contact_phone LIKE ?)", like, like, like)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表，按创建时间降序
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&leads).Error
	if err != nil {
		return nil, 0, err
	}

	return leads, total, nil
}
