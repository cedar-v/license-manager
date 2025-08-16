package service

import (
	"context"

	"license-manager/internal/models"
	pkgcontext "license-manager/pkg/context"
	"license-manager/pkg/i18n"
)

type enumService struct{}

// NewEnumService 创建枚举服务实例
func NewEnumService() EnumService {
	return &enumService{}
}

// GetAllEnums 获取所有枚举类型及其值
func (s *enumService) GetAllEnums(ctx context.Context) (*models.EnumListResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)
	
	// 从i18n管理器获取所有枚举值
	enumsMap := i18n.GetAllEnums(lang)
	
	var enums []models.EnumTypeResponse
	for enumType, enumItems := range enumsMap {
		var items []models.EnumItem
		for key, display := range enumItems {
			items = append(items, models.EnumItem{
				Key:     key,
				Display: display,
			})
		}
		
		enums = append(enums, models.EnumTypeResponse{
			Type:  enumType,
			Items: items,
		})
	}
	
	return &models.EnumListResponse{
		Enums: enums,
	}, nil
}

// GetEnumsByType 获取指定类型的枚举值
func (s *enumService) GetEnumsByType(ctx context.Context, enumType string) (*models.EnumTypeResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)
	
	// 业务逻辑：参数验证
	if enumType == "" {
		return nil, i18n.NewI18nError("900001", lang)
	}
	
	// 从i18n管理器获取指定类型的枚举值
	enumItems := i18n.GetEnumsByType(enumType, lang)
	
	var items []models.EnumItem
	for key, display := range enumItems {
		items = append(items, models.EnumItem{
			Key:     key,
			Display: display,
		})
	}
	
	return &models.EnumTypeResponse{
		Type:  enumType,
		Items: items,
	}, nil
}