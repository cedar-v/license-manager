package service

import (
	"context"
	"encoding/json"
	"fmt"
	"license-manager/internal/models"
	"license-manager/internal/repository"
	pkgcontext "license-manager/pkg/context"
	"license-manager/pkg/i18n"
	"time"

	"gorm.io/gorm"
)

// PackageService 套餐服务接口
type PackageService interface {
	// 管理端接口
	CreatePackage(ctx context.Context, req *models.PackageCreateRequest) (*models.Package, error)
	UpdatePackage(ctx context.Context, id string, req *models.PackageUpdateRequest) (*models.Package, error)
	DeletePackage(ctx context.Context, id string) error
	GetPackageList(ctx context.Context, req *models.PackageListRequest) (*models.PackageListResponse, error)
	GetPackageDetail(ctx context.Context, id string) (*models.Package, error)

	// 用户端接口
	GetCuPackageList(ctx context.Context) ([]*models.CuPackageResponse, error)
}

type packageService struct {
	repo repository.PackageRepository
	db   *gorm.DB
}

// NewPackageService 创建套餐服务
func NewPackageService(repo repository.PackageRepository, db *gorm.DB) PackageService {
	return &packageService{
		repo: repo,
		db:   db,
	}
}

// ========== 管理端接口实现 ==========

func (s *packageService) CreatePackage(ctx context.Context, req *models.PackageCreateRequest) (*models.Package, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	pkg := &models.Package{
		Name:                req.Name,
		Type:                req.Type,
		Price:               req.Price,
		PriceDescription:    req.PriceDescription,
		DurationDescription: req.DurationDescription,
		Description:         req.Description,
		Features:            req.Features,
		Status:              req.Status,
		SortOrder:           req.SortOrder,
		Remark:              req.Remark,
	}

	if err := s.repo.Create(pkg); err != nil {
		return nil, i18n.NewI18nError("800001", lang)
	}

	return pkg, nil
}

func (s *packageService) UpdatePackage(ctx context.Context, id string, req *models.PackageUpdateRequest) (*models.Package, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	pkg, err := s.repo.GetByID(id)
	if err != nil {
		return nil, i18n.NewI18nError("800002", lang)
	}

	// 更新字段
	if req.Name != "" {
		pkg.Name = req.Name
	}
	if req.Type != "" {
		pkg.Type = req.Type
	}
	if req.Price >= 0 {
		pkg.Price = req.Price
	}
	if req.PriceDescription != "" {
		pkg.PriceDescription = req.PriceDescription
	}
	if req.DurationDescription != "" {
		pkg.DurationDescription = req.DurationDescription
	}
	if req.Description != "" {
		pkg.Description = req.Description
	}
	if req.Features != "" {
		pkg.Features = req.Features
	}
	if req.Status != nil {
		pkg.Status = *req.Status
	}
	if req.SortOrder != nil {
		pkg.SortOrder = *req.SortOrder
	}
	if req.Remark != "" {
		pkg.Remark = req.Remark
	}

	pkg.UpdatedAt = time.Now()

	if err := s.repo.Update(pkg); err != nil {
		return nil, i18n.NewI18nError("800003", lang)
	}

	return pkg, nil
}

func (s *packageService) DeletePackage(ctx context.Context, id string) error {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 检查是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return i18n.NewI18nError("800002", lang)
	}

	if err := s.repo.Delete(id); err != nil {
		return i18n.NewI18nError("800004", lang)
	}

	return nil
}

func (s *packageService) GetPackageList(ctx context.Context, req *models.PackageListRequest) (*models.PackageListResponse, error) {
	packages, total, err := s.repo.GetList(ctx, req)
	if err != nil {
		lang := pkgcontext.GetLanguageFromContext(ctx)
		return nil, i18n.NewI18nError("800005", lang)
	}

	responses := make([]*models.PackageResponse, len(packages))
	for i, pkg := range packages {
		responses[i] = pkg.ToResponse()
	}

	return &models.PackageListResponse{
		Packages:   responses,
		TotalCount: total,
	}, nil
}

func (s *packageService) GetPackageDetail(ctx context.Context, id string) (*models.Package, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	pkg, err := s.repo.GetByID(id)
	if err != nil {
		return nil, i18n.NewI18nError("800002", lang)
	}

	return pkg, nil
}

// ========== 用户端接口实现 ==========

func (s *packageService) GetCuPackageList(ctx context.Context) ([]*models.CuPackageResponse, error) {
	packages, err := s.repo.GetEnabledList(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*models.CuPackageResponse, len(packages))
	for i, pkg := range packages {
		responses[i] = s.toCuPackageResponse(pkg)
	}

	return responses, nil
}

// toCuPackageResponse 转换为用户端套餐响应（兼容原有结构）
func (s *packageService) toCuPackageResponse(pkg *models.Package) *models.CuPackageResponse {
	// 从features中解析max_devices
	maxDevices := 1000 // 默认值
	if pkg.Features != "" {
		var features []string
		if err := json.Unmarshal([]byte(pkg.Features), &features); err == nil {
			for _, f := range features {
				// 尝试解析"X个许可"格式
				var n int
				if _, err := fmt.Sscanf(f, "%d个许可", &n); err == nil {
					maxDevices = n
					break
				}
			}
		}
	}

	// 构建details字符串
	details := pkg.Description
	if pkg.DurationDescription != "" {
		if details != "" {
			details += "|"
		}
		details += pkg.DurationDescription
	}
	if maxDevices > 0 {
		if details != "" {
			details += "|"
		}
		details += fmt.Sprintf("%d个许可", maxDevices)
	}

	return &models.CuPackageResponse{
		ID:          pkg.ID,
		Name:        pkg.Name,
		Type:        pkg.Type,
		Price:       pkg.Price,
		MaxDevices:  maxDevices,
		Description: pkg.Description,
		Features:    pkg.Features,
		Details:     details,
	}
}
