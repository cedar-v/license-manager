package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"license-manager/internal/models"
)

type customerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository 创建客户数据访问实例
func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

// GetCustomerList 查询客户列表
func (r *customerRepository) GetCustomerList(ctx context.Context, req *models.CustomerListRequest) (*models.CustomerListResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	if req.Sort == "" {
		req.Sort = "created_at"
	}
	if req.Order == "" {
		req.Order = "desc"
	}

	// 构建查询
	query := r.db.Model(&models.Customer{})

	// 搜索关键词筛选
	if req.Search != "" {
		searchTerm := "%" + strings.TrimSpace(req.Search) + "%"
		query = query.Where(
			"customer_code LIKE ? OR customer_name LIKE ? OR contact_person LIKE ? OR email LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	// 客户类型筛选
	if req.CustomerType != "" {
		query = query.Where("customer_type = ?", req.CustomerType)
	}

	// 客户等级筛选
	if req.CustomerLevel != "" {
		query = query.Where("customer_level = ?", req.CustomerLevel)
	}

	// 状态筛选
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count customers: %w", err)
	}

	// 排序和分页
	offset := (req.Page - 1) * req.PageSize
	orderClause := fmt.Sprintf("%s %s", req.Sort, strings.ToUpper(req.Order))
	query = query.Order(orderClause).Offset(offset).Limit(req.PageSize)

	// 查询数据
	var customers []models.Customer
	if err := query.Find(&customers).Error; err != nil {
		return nil, fmt.Errorf("failed to query customers: %w", err)
	}

	// 转换为响应格式
	customerList := make([]models.CustomerListItem, 0, len(customers))
	for _, customer := range customers {
		customerList = append(customerList, models.CustomerListItem{
			ID:            customer.ID,
			CustomerCode:  customer.CustomerCode,
			CustomerName:  customer.CustomerName,
			CustomerType:  customer.CustomerType,
			ContactPerson: customer.ContactPerson,
			Email:         customer.Email,
			CustomerLevel: customer.CustomerLevel,
			Status:        customer.Status,
			CreatedAt:     customer.CreatedAt.Format(time.RFC3339),
		})
	}

	// 计算总页数
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &models.CustomerListResponse{
		List:       customerList,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetCustomerByID 根据ID获取客户信息
func (r *customerRepository) GetCustomerByID(ctx context.Context, id string) (*models.Customer, error) {
	var customer models.Customer
	
	if err := r.db.Where("id = ?", id).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrCustomerNotFound
		}
		return nil, fmt.Errorf("failed to get customer by id: %w", err)
	}
	return &customer, nil
}

// CreateCustomer 创建客户
func (r *customerRepository) CreateCustomer(ctx context.Context, customer *models.Customer) error {
	if err := r.db.Create(customer).Error; err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}
	return nil
}

// UpdateCustomer 更新客户信息
func (r *customerRepository) UpdateCustomer(ctx context.Context, customer *models.Customer) error {
	if err := r.db.Save(customer).Error; err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}
	return nil
}

// DeleteCustomer 删除客户
func (r *customerRepository) DeleteCustomer(ctx context.Context, id string) error {
	result := r.db.Where("id = ?", id).Delete(&models.Customer{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete customer: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrCustomerNotFound
	}
	return nil
}

// GetCustomerCount 获取客户总数（用于统计）
func (r *customerRepository) GetCustomerCount(ctx context.Context, filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.Model(&models.Customer{})
	
	// 应用筛选条件
	for key, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}
	
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count customers: %w", err)
	}
	
	return count, nil
}