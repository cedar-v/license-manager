package service

import (
	"context"
	"errors"
	"strings"

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

	// 业务逻辑：添加多语言显示字段
	for i := range result.List {
		s.fillDisplayFields(&result.List[i], lang)
	}

	return result, nil
}

// GetCustomer 获取单个客户详情
func (s *customerService) GetCustomer(ctx context.Context, id string) (*models.Customer, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if id == "" {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 委托给Repository层进行数据访问
	customer, err := s.customerRepo.GetCustomerByID(ctx, id)
	if err != nil {
		// 根据Repository错误类型包装为完整的I18nError
		if errors.Is(err, repository.ErrCustomerNotFound) {
			return nil, i18n.NewI18nError("200001", lang)
		}

		// 数据库相关错误
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 获取授权统计信息
	stats, err := s.customerRepo.GetCustomerAuthorizationStats(ctx, id)
	if err != nil {
		// 统计信息获取失败不影响主流程，记录错误但继续执行
		// 如果获取失败，stats 为 nil，客户详情仍然可以返回，只是没有统计信息
		// 根据需求文档，空数据时应该返回0，所以这里创建一个空的统计对象
		stats = &models.AuthorizationStats{}
	}
	customer.AuthorizationStats = stats

	// 填充多语言显示字段
	s.fillCustomerDisplayFields(customer, lang)

	return customer, nil
}

// fillDisplayFields 填充列表项多语言显示字段
func (s *customerService) fillDisplayFields(item *models.CustomerListItem, lang string) {
	item.CustomerTypeDisplay = i18n.GetEnumMessage("customer_type", item.CustomerType, lang)
	item.CustomerLevelDisplay = i18n.GetEnumMessage("customer_level", item.CustomerLevel, lang)
	item.StatusDisplay = i18n.GetEnumMessage("customer_status", item.Status, lang)
}

// CreateCustomer 创建客户
func (s *customerService) CreateCustomer(ctx context.Context, req *models.CustomerCreateRequest) (*models.Customer, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 获取当前用户ID（从Context中获取）
	// TODO: 实现用户上下文获取，这里暂时使用硬编码
	currentUserID := "admin_uuid" // 实际应该从JWT token中获取

	// 构建客户实体
	customer := &models.Customer{
		CustomerName:  req.CustomerName,
		CustomerType:  req.CustomerType,
		ContactPerson: req.ContactPerson,
		ContactTitle:  req.ContactTitle,
		Email:         req.Email,
		Phone:         req.Phone,
		Address:       req.Address,
		CompanySize:   req.CompanySize,
		CustomerLevel: req.CustomerLevel,
		Status:        req.Status,
		Description:   req.Description,
		CreatedBy:     currentUserID,
		// CustomerCode 将在Repository层自动生成
		// CreatedAt, UpdatedAt 将在模型的BeforeCreate中自动设置
	}

	// 委托给Repository层进行数据创建
	if err := s.customerRepo.CreateCustomer(ctx, customer); err != nil {
		// 根据错误类型包装为完整的I18nError
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return nil, i18n.NewI18nError("200002", lang) // 客户已存在
		}

		// 数据库相关错误
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 填充多语言显示字段
	s.fillCustomerDisplayFields(customer, lang)

	return customer, nil
}

// UpdateCustomer 更新客户信息
func (s *customerService) UpdateCustomer(ctx context.Context, id string, req *models.CustomerUpdateRequest) (*models.Customer, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if id == "" {
		return nil, i18n.NewI18nError("900001", lang)
	}
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 先查询现有客户
	existingCustomer, err := s.customerRepo.GetCustomerByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerNotFound) {
			return nil, i18n.NewI18nError("200001", lang)
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 获取当前用户ID
	currentUserID := "admin_uuid" // TODO: 从JWT token中获取

	// 只更新提供的字段
	if req.CustomerName != nil {
		existingCustomer.CustomerName = *req.CustomerName
	}
	if req.CustomerType != nil {
		existingCustomer.CustomerType = *req.CustomerType
	}
	if req.ContactPerson != nil {
		existingCustomer.ContactPerson = *req.ContactPerson
	}
	if req.ContactTitle != nil {
		existingCustomer.ContactTitle = req.ContactTitle
	}
	if req.Email != nil {
		existingCustomer.Email = req.Email
	}
	if req.Phone != nil {
		existingCustomer.Phone = req.Phone
	}
	if req.Address != nil {
		existingCustomer.Address = req.Address
	}
	if req.CompanySize != nil {
		existingCustomer.CompanySize = req.CompanySize
	}
	if req.CustomerLevel != nil {
		existingCustomer.CustomerLevel = *req.CustomerLevel
	}
	if req.Status != nil {
		existingCustomer.Status = *req.Status
	}
	if req.Description != nil {
		existingCustomer.Description = req.Description
	}

	// 设置更新者
	existingCustomer.UpdatedBy = &currentUserID

	// 委托给Repository层进行数据更新
	if err := s.customerRepo.UpdateCustomer(ctx, existingCustomer); err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 填充多语言显示字段
	s.fillCustomerDisplayFields(existingCustomer, lang)

	return existingCustomer, nil
}

// DeleteCustomer 删除客户（物理删除，删除前检查是否有授权）
func (s *customerService) DeleteCustomer(ctx context.Context, id string) error {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if id == "" {
		return i18n.NewI18nError("900001", lang)
	}

	// 检查客户是否存在
	_, err := s.customerRepo.GetCustomerByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerNotFound) {
			return i18n.NewI18nError("200001", lang)
		}
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	// 检查是否有关联的授权码
	hasAuthCodes, err := s.customerRepo.CheckCustomerHasAuthorizationCodes(ctx, id)
	if err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}
	if hasAuthCodes {
		return i18n.NewI18nError("200006", lang)
	}

	// 检查是否有关联的许可证
	hasLicenses, err := s.customerRepo.CheckCustomerHasLicenses(ctx, id)
	if err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}
	if hasLicenses {
		return i18n.NewI18nError("200006", lang)
	}

	// 委托给Repository层进行物理删除
	if err := s.customerRepo.DeleteCustomer(ctx, id); err != nil {
		if errors.Is(err, repository.ErrCustomerNotFound) {
			return i18n.NewI18nError("200001", lang)
		}
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	return nil
}

// UpdateCustomerStatus 更新客户状态
func (s *customerService) UpdateCustomerStatus(ctx context.Context, id string, req *models.CustomerStatusUpdateRequest) (*models.Customer, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if id == "" {
		return nil, i18n.NewI18nError("900001", lang)
	}
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 先查询现有客户
	existingCustomer, err := s.customerRepo.GetCustomerByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerNotFound) {
			return nil, i18n.NewI18nError("200001", lang)
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 获取当前用户ID
	currentUserID := "admin_uuid" // TODO: 从JWT token中获取

	// 更新状态
	existingCustomer.Status = req.Status
	existingCustomer.UpdatedBy = &currentUserID

	// 委托给Repository层进行数据更新
	if err := s.customerRepo.UpdateCustomer(ctx, existingCustomer); err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 填充多语言显示字段
	s.fillCustomerDisplayFields(existingCustomer, lang)

	return existingCustomer, nil
}

// fillCustomerDisplayFields 填充完整客户模型的多语言显示字段
func (s *customerService) fillCustomerDisplayFields(customer *models.Customer, lang string) {
	customer.CustomerTypeDisplay = i18n.GetEnumMessage("customer_type", customer.CustomerType, lang)
	customer.CustomerLevelDisplay = i18n.GetEnumMessage("customer_level", customer.CustomerLevel, lang)
	customer.StatusDisplay = i18n.GetEnumMessage("customer_status", customer.Status, lang)
	if customer.CompanySize != nil {
		customer.CompanySizeDisplay = i18n.GetEnumMessage("company_size", *customer.CompanySize, lang)
	}
}
