package repository

import (
	"context"
	"license-manager/internal/models"

	"gorm.io/gorm"
)

type authorizationCodeRepository struct {
	db *gorm.DB
}

// NewAuthorizationCodeRepository 创建授权码仓储实例
func NewAuthorizationCodeRepository(db *gorm.DB) AuthorizationCodeRepository {
	return &authorizationCodeRepository{
		db: db,
	}
}

// CreateAuthorizationCode 创建授权码
func (r *authorizationCodeRepository) CreateAuthorizationCode(ctx context.Context, authCode *models.AuthorizationCode) error {
	return r.db.WithContext(ctx).Create(authCode).Error
}

// GetAuthorizationCodeByID 根据ID获取授权码
func (r *authorizationCodeRepository) GetAuthorizationCodeByID(ctx context.Context, id string) (*models.AuthorizationCode, error) {
	var authCode models.AuthorizationCode
	err := r.db.WithContext(ctx).First(&authCode, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrAuthorizationCodeNotFound
		}
		return nil, err
	}
	return &authCode, nil
}

// GetAuthorizationCodeList 查询授权码列表
func (r *authorizationCodeRepository) GetAuthorizationCodeList(ctx context.Context, req *models.AuthorizationCodeListRequest) (*models.AuthorizationCodeListResponse, error) {
	// 设置默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	sort := req.Sort
	if sort == "" {
		sort = "created_at"
	}
	order := req.Order
	if order == "" {
		order = "desc"
	}

	// 构建查询
	query := r.db.WithContext(ctx).Table("authorization_codes ac").
		Select(`ac.id, ac.code, ac.customer_id, c.customer_name, 
				ac.start_date, ac.end_date, ac.max_activations, 
				ac.deployment_type, ac.is_locked, ac.description, ac.created_at`).
		Joins("LEFT JOIN customers c ON ac.customer_id = c.id AND c.deleted_at IS NULL").
		Where("ac.deleted_at IS NULL")

	// 添加筛选条件
	if req.CustomerID != "" {
		query = query.Where("ac.customer_id = ?", req.CustomerID)
	}

	// 状态筛选需要在应用层处理，因为status是计算字段
	if req.StartDate != "" {
		query = query.Where("ac.created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("ac.created_at <= ?", req.EndDate)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	var results []struct {
		ID             string  `json:"id"`
		Code           string  `json:"code"`
		CustomerID     string  `json:"customer_id"`
		CustomerName   *string `json:"customer_name"`
		StartDate      string  `json:"start_date"`
		EndDate        string  `json:"end_date"`
		MaxActivations int     `json:"max_activations"`
		DeploymentType string  `json:"deployment_type"`
		IsLocked       bool    `json:"is_locked"`
		Description    *string `json:"description"`
		CreatedAt      string  `json:"created_at"`
	}

	offset := (page - 1) * pageSize
	orderClause := sort + " " + order
	if err := query.Order(orderClause).Limit(pageSize).Offset(offset).Scan(&results).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]models.AuthorizationCodeListItem, len(results))
	for i, result := range results {
		customerName := ""
		if result.CustomerName != nil {
			customerName = *result.CustomerName
		}
		
		list[i] = models.AuthorizationCodeListItem{
			ID:                     result.ID,
			Code:                   result.Code,
			CustomerID:             result.CustomerID,
			CustomerName:           customerName,
			StartDate:              result.StartDate,
			EndDate:                result.EndDate,
			MaxActivations:         result.MaxActivations,
			CurrentActivations:     0, // TODO: 从licenses表统计
			DeploymentType:         result.DeploymentType,
			IsLocked:               result.IsLocked,
			Description:            result.Description,
			CreatedAt:              result.CreatedAt,
		}
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &models.AuthorizationCodeListResponse{
		List:       list,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateAuthorizationCode 更新授权码
func (r *authorizationCodeRepository) UpdateAuthorizationCode(ctx context.Context, authCode *models.AuthorizationCode) error {
	return r.db.WithContext(ctx).Save(authCode).Error
}

// DeleteAuthorizationCode 删除授权码（软删除）
func (r *authorizationCodeRepository) DeleteAuthorizationCode(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.AuthorizationCode{}, "id = ?", id).Error
}

// CheckCustomerExists 检查客户是否存在
func (r *authorizationCodeRepository) CheckCustomerExists(ctx context.Context, customerID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Customer{}).Where("id = ? AND deleted_at IS NULL", customerID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}