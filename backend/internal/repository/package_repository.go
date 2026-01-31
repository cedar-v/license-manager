package repository

import (
	"context"
	"license-manager/internal/models"

	"gorm.io/gorm"
)

// PackageRepository 套餐仓储接口
type PackageRepository interface {
	Create(pkg *models.Package) error
	Update(pkg *models.Package) error
	Delete(id string) error
	GetByID(id string) (*models.Package, error)
	GetList(ctx context.Context, req *models.PackageListRequest) ([]*models.Package, int64, error)
	GetEnabledList(ctx context.Context) ([]*models.Package, error)
	GetByType(ctx context.Context, pkgType string) ([]*models.Package, error)
}

type packageRepository struct {
	db *gorm.DB
}

// NewPackageRepository 创建套餐仓储
func NewPackageRepository(db *gorm.DB) PackageRepository {
	return &packageRepository{db: db}
}

func (r *packageRepository) Create(pkg *models.Package) error {
	return r.db.Create(pkg).Error
}

func (r *packageRepository) Update(pkg *models.Package) error {
	return r.db.Save(pkg).Error
}

func (r *packageRepository) Delete(id string) error {
	return r.db.Model(&models.Package{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *packageRepository) GetByID(id string) (*models.Package, error) {
	var pkg models.Package
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&pkg).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *packageRepository) GetList(ctx context.Context, req *models.PackageListRequest) ([]*models.Package, int64, error) {
	var packages []*models.Package
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Package{}).Where("deleted_at IS NULL")

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表，按sort_order降序排序
	err := query.Order("sort_order DESC, created_at DESC").Find(&packages).Error
	if err != nil {
		return nil, 0, err
	}

	return packages, total, nil
}

func (r *packageRepository) GetEnabledList(ctx context.Context) ([]*models.Package, error) {
	var packages []*models.Package
	err := r.db.WithContext(ctx).
		Where("status = ? AND deleted_at IS NULL", 1).
		Order("sort_order ASC, created_at ASC").
		Find(&packages).Error
	if err != nil {
		return nil, err
	}
	return packages, nil
}

func (r *packageRepository) GetByType(ctx context.Context, pkgType string) ([]*models.Package, error) {
	var packages []*models.Package
	err := r.db.WithContext(ctx).
		Where("type = ? AND status = ? AND deleted_at IS NULL", pkgType, 1).
		Order("sort_order DESC").
		Find(&packages).Error
	if err != nil {
		return nil, err
	}
	return packages, nil
}
