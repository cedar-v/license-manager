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

	// 构建查询，包含状态计算与当前激活数量
	query := r.db.WithContext(ctx).Table("authorization_codes ac").
		Select(`ac.id, ac.code, ac.customer_id, c.customer_name, 
				ac.start_date, ac.end_date, ac.max_activations, 
				COALESCE(l.active_count, 0) AS current_activations,
				ac.deployment_type, ac.is_locked, ac.description, ac.created_at,
				CASE 
					WHEN ac.is_locked = true THEN 'locked'
					WHEN ac.end_date < NOW() THEN 'expired'
					WHEN ac.start_date <= NOW() AND ac.end_date >= NOW() THEN 'normal'
					ELSE 'expired'
				END AS status`).
		Joins("LEFT JOIN customers c ON ac.customer_id = c.id AND c.deleted_at IS NULL").
		Joins(`LEFT JOIN (
			SELECT authorization_code_id, COUNT(*) AS active_count
			FROM licenses
			WHERE status = 'active' AND deleted_at IS NULL
			GROUP BY authorization_code_id
		) l ON ac.id = l.authorization_code_id`)

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
		ID                 string  `json:"id"`
		Code               string  `json:"code"`
		CustomerID         string  `json:"customer_id"`
		CustomerName       *string `json:"customer_name"`
		StartDate          string  `json:"start_date"`
		EndDate            string  `json:"end_date"`
		MaxActivations     int     `json:"max_activations"`
		CurrentActivations int     `json:"current_activations" gorm:"column:current_activations"`
		DeploymentType     string  `json:"deployment_type"`
		IsLocked           bool    `json:"is_locked"`
		Description        *string `json:"description"`
		CreatedAt          string  `json:"created_at"`
		Status             string  `json:"status"` // 添加状态字段
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
			CurrentActivations: result.CurrentActivations,
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

// GetCuAuthorizationCodeList 用户端：查询用户授权码列表
func (r *authorizationCodeRepository) GetCuAuthorizationCodeList(ctx context.Context, cuUserID string, req *models.CuAuthorizationCodeListRequest) (*models.CuAuthorizationCodeListResponse, error) {
	page := 1
	pageSize := 10
	status := ""
	search := ""
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
		status = req.Status
		search = req.Search
	}

	query := r.db.WithContext(ctx).Table("authorization_codes ac").
		Select(`ac.id, ac.code, ac.created_at, ac.end_date, ac.max_activations,
				COALESCE(l.active_count, 0) AS current_activations,
				CASE
					WHEN ac.end_date < NOW() THEN 'expired'
					WHEN ac.is_locked = true THEN 'locked'
					WHEN ac.start_date <= NOW() AND ac.end_date >= NOW() THEN 'normal'
					ELSE 'expired'
				END AS status`).
		Joins(`LEFT JOIN (
			SELECT authorization_code_id, COUNT(*) AS active_count
			FROM licenses
			WHERE status = 'active' AND deleted_at IS NULL
			GROUP BY authorization_code_id
		) l ON ac.id = l.authorization_code_id`).
		Where("ac.created_by = ?", cuUserID)

	if status != "" {
		switch status {
		case "locked":
			query = query.Where("ac.end_date >= NOW() AND ac.is_locked = true")
		case "expired":
			query = query.Where("ac.end_date < NOW()")
		case "normal":
			query = query.Where("ac.end_date >= NOW() AND ac.is_locked = false")
		}
	}

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("ac.code LIKE ?", like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	var results []struct {
		ID                 string    `json:"id"`
		Code               string    `json:"code"`
		CreatedAt          time.Time `json:"created_at"`
		EndDate            time.Time `json:"end_date"`
		MaxActivations     int       `json:"max_activations"`
		CurrentActivations int       `json:"current_activations" gorm:"column:current_activations"`
		Status             string    `json:"status"`
	}

	if err := query.Order("ac.created_at DESC").Limit(pageSize).Offset(offset).Scan(&results).Error; err != nil {
		return nil, err
	}

	list := make([]models.CuAuthorizationCodeListItem, len(results))
	for i, item := range results {
		remaining := item.MaxActivations - item.CurrentActivations
		if remaining < 0 {
			remaining = 0
		}
		list[i] = models.CuAuthorizationCodeListItem{
			ID:                   item.ID,
			Code:                 item.Code,
			Status:               item.Status,
			MaxActivations:       item.MaxActivations,
			CurrentActivations:   item.CurrentActivations,
			RemainingActivations: remaining,
			CreatedAt:            item.CreatedAt.Format(time.RFC3339),
			EndDate:              item.EndDate.Format(time.RFC3339),
		}
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return &models.CuAuthorizationCodeListResponse{
		List:       list,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetCuAuthorizationCodeSummary 用户端：授权信息统计
func (r *authorizationCodeRepository) GetCuAuthorizationCodeSummary(ctx context.Context, cuUserID string) (*models.CuAuthorizationCodeSummaryResponse, error) {
	var result struct {
		TotalCount             int64 `gorm:"column:total_count"`
		ExpiredCount           int64 `gorm:"column:expired_count"`
		ValidCount             int64 `gorm:"column:valid_count"`
		ValidMaxActivationsSum int64 `gorm:"column:valid_max_activations_sum"`
	}

	err := r.db.WithContext(ctx).Table("authorization_codes ac").
		Select(`
			COUNT(*) AS total_count,
			SUM(CASE WHEN ac.end_date < NOW() THEN 1 ELSE 0 END) AS expired_count,
			SUM(CASE WHEN ac.end_date >= NOW() THEN 1 ELSE 0 END) AS valid_count,
			SUM(CASE WHEN ac.end_date >= NOW() THEN ac.max_activations ELSE 0 END) AS valid_max_activations_sum
		`).
		Where("ac.created_by = ?", cuUserID).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return &models.CuAuthorizationCodeSummaryResponse{
		TotalCount:             result.TotalCount,
		ExpiredCount:           result.ExpiredCount,
		ValidCount:             result.ValidCount,
		ValidMaxActivationsSum: result.ValidMaxActivationsSum,
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

// BeginTransaction 开始事务
func (r *authorizationCodeRepository) BeginTransaction(ctx context.Context) interface{} {
	return r.db.WithContext(ctx).Begin()
}

// CreateAuthorizationCodeWithTx 在事务中创建授权码
func (r *authorizationCodeRepository) CreateAuthorizationCodeWithTx(ctx context.Context, tx interface{}, authCode *models.AuthorizationCode) error {
	if gormTx, ok := tx.(*gorm.DB); ok {
		return gormTx.WithContext(ctx).Create(authCode).Error
	}
	return ErrInvalidTransaction
}

// UpdateMaxActivationsWithTx 在事务中更新授权码的最大激活次数
func (r *authorizationCodeRepository) UpdateMaxActivationsWithTx(ctx context.Context, tx interface{}, authCodeID string, newMaxActivations int) error {
	if gormTx, ok := tx.(*gorm.DB); ok {
		return gormTx.WithContext(ctx).Model(&models.AuthorizationCode{}).Where("id = ?", authCodeID).Update("max_activations", newMaxActivations).Error
	}
	return ErrInvalidTransaction
}
