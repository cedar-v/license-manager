package service

import (
	"context"

	"license-manager/internal/models"
	"license-manager/internal/repository"
	pkgcontext "license-manager/pkg/context"
	"license-manager/pkg/i18n"
)

type cuDeviceService struct {
	licenseRepo repository.LicenseRepository
}

// NewCuDeviceService 创建客户设备服务实例
func NewCuDeviceService(licenseRepo repository.LicenseRepository) CuDeviceService {
	return &cuDeviceService{
		licenseRepo: licenseRepo,
	}
}

// GetDeviceList 获取设备列表
func (s *cuDeviceService) GetDeviceList(ctx context.Context, customerID string, req *models.DeviceListRequest) (*models.DeviceListResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 参数验证
	if customerID == "" {
		return nil, i18n.NewI18nError("900001", lang) // 业务错误，不覆盖多语言message
	}
	if req == nil {
		req = &models.DeviceListRequest{}
	}

	// 委托给Repository层进行数据访问
	return s.licenseRepo.GetCustomerDeviceList(ctx, customerID, req)
}

// UnbindDevice 解绑设备
func (s *cuDeviceService) UnbindDevice(ctx context.Context, customerID, licenseID string) error {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 参数验证
	if customerID == "" || licenseID == "" {
		return i18n.NewI18nError("900001", lang) // 业务错误，不覆盖多语言message
	}

	// 检查许可证是否存在且属于当前客户
	belongs, err := s.licenseRepo.CheckLicenseBelongsToCustomer(ctx, licenseID, customerID)
	if err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}
	if !belongs {
		return i18n.NewI18nError("620001", lang) // 设备不存在或无权限访问
	}

	// 物理删除许可证记录
	err = s.licenseRepo.DeleteLicenseByID(ctx, licenseID)
	if err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	return nil
}
