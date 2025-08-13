package service

import (
	"context"
	"errors"

	"license-manager/internal/models"
	"license-manager/internal/repository"
	pkgcontext "license-manager/pkg/context"
	"license-manager/pkg/i18n"
)

type customerService struct {
	customerRepo repository.CustomerRepository
}

// NewCustomerService 创建客户服务实例
func NewCustomerService(customerRepo repository.CustomerRepository) CustomerService {
	return &customerService{
		customerRepo: customerRepo,
	}
}

// GetCustomerList 查询客户列表
func (s *customerService) GetCustomerList(ctx context.Context, req *models.CustomerListRequest) (*models.CustomerListResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)
	
	// 业务逻辑：参数验证
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang) // 业务错误，不覆盖多语言message
	}

	// 委托给Repository层进行数据访问
	result, err := s.customerRepo.GetCustomerList(ctx, req)
	if err != nil {
		// 根据Repository错误类型包装为完整的I18nError
		if errors.Is(err, repository.ErrCustomerNotFound) {
			return nil, i18n.NewI18nError("200001", lang) // 业务错误，保持多语言
		}
		
		// 数据库相关错误，使用系统错误码，覆盖message显示详细信息
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 业务逻辑：可以在这里添加额外的业务规则
	// 比如：权限检查、数据过滤、统计计算等

	return result, nil
}
