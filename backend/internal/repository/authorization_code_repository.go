package repository

import (
	"context"
	"license-manager/internal/models"
	"math"
	"time"

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

	// 构建查询，包含状态计算
	query := r.db.WithContext(ctx).Table("authorization_codes ac").
		Select(`ac.id, ac.code, ac.customer_id, c.customer_name, 
				ac.start_date, ac.end_date, ac.max_activations, 
				ac.deployment_type, ac.is_locked, ac.description, ac.created_at,
				CASE 
					WHEN ac.is_locked = true THEN 'locked'
					WHEN ac.end_date < NOW() THEN 'expired'
					WHEN ac.start_date <= NOW() AND ac.end_date >= NOW() THEN 'normal'
					ELSE 'expired'
				END AS status`).
		Joins("LEFT JOIN customers c ON ac.customer_id = c.id AND c.deleted_at IS NULL")

	// 添加筛选条件
	if req.CustomerID != "" {
		query = query.Where("ac.customer_id = ?", req.CustomerID)
	}

	// 状态筛选 - 在 SQL 层面处理
	if req.Status != "" {
		switch req.Status {
		case "locked":
			query = query.Where("ac.is_locked = true")
		case "expired":
			query = query.Where("ac.is_locked = false AND ac.end_date < NOW()")
		case "normal":
			query = query.Where("ac.is_locked = false AND ac.start_date <= NOW() AND ac.end_date >= NOW()")
		}
	}

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
		Status         string  `json:"status"` // 添加状态字段
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
			ID:                 result.ID,
			Code:               result.Code,
			CustomerID:         result.CustomerID,
			CustomerName:       customerName,
			Status:             result.Status, // 使用 SQL 计算的状态
			StartDate:          result.StartDate,
			EndDate:            result.EndDate,
			MaxActivations:     result.MaxActivations,
			CurrentActivations: 0, // TODO: 从licenses表统计
			DeploymentType:     result.DeploymentType,
			IsLocked:           result.IsLocked,
			Description:        result.Description,
			CreatedAt:          result.CreatedAt,
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

// GetAuthorizationChangeList 查询授权变更历史列表
func (r *authorizationCodeRepository) GetAuthorizationChangeList(ctx context.Context, authCodeID string, req *models.AuthorizationChangeListRequest) (*models.AuthorizationChangeListResponse, error) {
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
	query := r.db.WithContext(ctx).Model(&models.AuthorizationChange{}).
		Select(`authorization_changes.id, authorization_changes.change_type, 
			authorization_changes.operator_id, users.username as operator_name,
			authorization_changes.reason, authorization_changes.created_at`).
		Joins("LEFT JOIN users ON authorization_changes.operator_id = users.id").
		Where("authorization_changes.authorization_code_id = ?", authCodeID)

	// 变更类型筛选
	if req.ChangeType != "" {
		query = query.Where("authorization_changes.change_type = ?", req.ChangeType)
	}

	// 操作人筛选
	if req.OperatorID != "" {
		query = query.Where("authorization_changes.operator_id = ?", req.OperatorID)
	}

	// 时间范围筛选
	if req.StartDate != "" {
		startTime, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			query = query.Where("authorization_changes.created_at >= ?", startTime)
		}
	}
	if req.EndDate != "" {
		endTime, err := time.Parse("2006-01-02", req.EndDate)
		if err == nil {
			// 结束时间加一天，以包含当天的所有时间
			endTime = endTime.AddDate(0, 0, 1)
			query = query.Where("authorization_changes.created_at < ?", endTime)
		}
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 排序和分页
	orderBy := "authorization_changes." + req.Sort + " " + req.Order
	offset := (req.Page - 1) * req.PageSize

	var changes []struct {
		ID           string    `json:"id"`
		ChangeType   string    `json:"change_type"`
		OperatorID   string    `json:"operator_id"`
		OperatorName *string   `json:"operator_name"`
		Reason       *string   `json:"reason"`
		CreatedAt    time.Time `json:"created_at"`
	}

	if err := query.Order(orderBy).Limit(req.PageSize).Offset(offset).Scan(&changes).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]models.AuthorizationChangeListItem, len(changes))
	for i, change := range changes {
		var operatorName string
		if change.OperatorName != nil {
			operatorName = *change.OperatorName
		}

		list[i] = models.AuthorizationChangeListItem{
			ID:           change.ID,
			ChangeType:   change.ChangeType,
			OperatorID:   change.OperatorID,
			OperatorName: operatorName,
			Reason:       change.Reason,
			CreatedAt:    change.CreatedAt.Format(time.RFC3339),
		}
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &models.AuthorizationChangeListResponse{
		List:       list,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// RecordAuthorizationChange 记录授权变更历史
func (r *authorizationCodeRepository) RecordAuthorizationChange(ctx context.Context, change *models.AuthorizationChange) error {
	return r.db.WithContext(ctx).Create(change).Error
}
