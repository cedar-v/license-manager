package repository

import (
	"context"
	"time"

	"license-manager/internal/models"
	pkgcontext "license-manager/pkg/context"
	"license-manager/pkg/i18n"

	"gorm.io/gorm"
)

type dashboardRepository struct {
	db *gorm.DB
}

// NewDashboardRepository 创建仪表盘数据访问层实例
func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{
		db: db,
	}
}

// GetAuthorizationTrendData 获取授权趋势数据
func (r *dashboardRepository) GetAuthorizationTrendData(ctx context.Context, startDate, endDate time.Time) ([]models.TrendData, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)
	var trendData []models.TrendData

	// 生成日期范围
	dates := make([]time.Time, 0)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d)
	}

	// 查询每日数据
	for _, date := range dates {
		dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		dayEnd := dayStart.Add(24*time.Hour - time.Nanosecond)

		// 统计当日的授权数据
		var totalCount, newCount, expiredCount int64

		// 当日总授权数（截至当日24:00的累计有效授权数）
		err := r.db.WithContext(ctx).Model(&models.AuthorizationCode{}).
			Where("created_at <= ?", dayEnd).
			Count(&totalCount).Error
		if err != nil {
			return nil, i18n.NewI18nError("900004", lang, err.Error())
		}

		// 当日新增授权数
		err = r.db.WithContext(ctx).Model(&models.AuthorizationCode{}).
			Where("created_at >= ? AND created_at <= ?", dayStart, dayEnd).
			Count(&newCount).Error
		if err != nil {
			return nil, i18n.NewI18nError("900004", lang, err.Error())
		}

		// 当日过期授权数
		err = r.db.WithContext(ctx).Model(&models.AuthorizationCode{}).
			Where("end_date >= ? AND end_date < ?", dayStart, dayEnd).
			Count(&expiredCount).Error
		if err != nil {
			return nil, i18n.NewI18nError("900004", lang, err.Error())
		}

		trendData = append(trendData, models.TrendData{
			Date:                  date.Format("2006-01-02"),
			TotalAuthorizations:   totalCount,
			NewAuthorizations:     newCount,
			ExpiredAuthorizations: expiredCount,
		})
	}

	return trendData, nil
}

// GetRecentAuthorizations 获取最近授权列表
func (r *dashboardRepository) GetRecentAuthorizations(ctx context.Context, req *models.DashboardRecentAuthorizationsRequest) (*models.DashboardRecentAuthorizationsResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 设置默认值
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// 构建查询
	query := r.db.WithContext(ctx).Table("authorization_codes ac").
		Select(`ac.id, ac.code, ac.customer_id, c.customer_name as customer_name, 
		        ac.description, ac.start_date, ac.end_date, ac.max_activations,
		        COALESCE(l.active_count, 0) as current_activations, 
		        ac.created_at, ac.updated_at, ac.is_locked`).
		Joins("LEFT JOIN customers c ON ac.customer_id = c.id").
		Joins(`LEFT JOIN (
			SELECT authorization_code_id, COUNT(*) as active_count 
			FROM licenses 
			WHERE status = 'active' AND deleted_at IS NULL 
			GROUP BY authorization_code_id
		) l ON ac.id = l.authorization_code_id`).
		Where("1=1").
		Order("ac.created_at DESC")

	// 添加筛选条件
	if req.CustomerID != "" {
		query = query.Where("ac.customer_id = ?", req.CustomerID)
	}

	// 状态筛选（虚字段）
	now := time.Now()
	switch req.Status {
	case "normal":
		query = query.Where("ac.is_locked = ? AND ac.end_date > ?", false, now)
	case "locked":
		query = query.Where("ac.is_locked = ?", true)
	case "expired":
		query = query.Where("ac.end_date <= ?", now)
	}

	// 查询总数 - 使用简化查询避免JOIN复杂性
	var total int64
	countQuery := r.db.WithContext(ctx).Model(&models.AuthorizationCode{}).Where("1=1")
	
	// 添加相同的筛选条件
	if req.CustomerID != "" {
		countQuery = countQuery.Where("customer_id = ?", req.CustomerID)
	}
	
	// 状态筛选
	switch req.Status {
	case "normal":
		countQuery = countQuery.Where("is_locked = ? AND end_date > ?", false, now)
	case "locked":
		countQuery = countQuery.Where("is_locked = ?", true)
	case "expired":
		countQuery = countQuery.Where("end_date <= ?", now)
	}
	
	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 查询列表数据
	type queryResult struct {
		ID                 string    `gorm:"column:id"`
		Code               string    `gorm:"column:code"`
		CustomerID         string    `gorm:"column:customer_id"`
		CustomerName       string    `gorm:"column:customer_name"`
		Description        string    `gorm:"column:description"`
		StartDate          time.Time `gorm:"column:start_date"`
		EndDate            time.Time `gorm:"column:end_date"`
		MaxActivations     int       `gorm:"column:max_activations"`
		CurrentActivations int       `gorm:"column:current_activations"`
		CreatedAt          time.Time `gorm:"column:created_at"`
		UpdatedAt          time.Time `gorm:"column:updated_at"`
		IsLocked           bool      `gorm:"column:is_locked"`
	}

	var results []queryResult
	err = query.Limit(limit).Find(&results).Error
	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 转换为响应格式
	authorizations := make([]models.RecentAuthorization, 0, len(results))
	for _, result := range results {
		// 计算状态
		status := "normal"
		statusDisplay := i18n.GetEnumMessage("authorization_code_status", "normal", lang)

		if result.IsLocked {
			status = "locked"
			statusDisplay = i18n.GetEnumMessage("authorization_code_status", "locked", lang)
		} else if result.EndDate.Before(now) {
			status = "expired"
			statusDisplay = i18n.GetEnumMessage("authorization_code_status", "expired", lang)
		}

		authorizations = append(authorizations, models.RecentAuthorization{
			ID:                 result.ID,
			Code:               result.Code,
			CustomerID:         result.CustomerID,
			CustomerName:       result.CustomerName,
			Description:        result.Description,
			Status:             status,
			StatusDisplay:      statusDisplay,
			StartDate:          result.StartDate,
			EndDate:            result.EndDate,
			MaxActivations:     result.MaxActivations,
			CurrentActivations: result.CurrentActivations,
			CreatedAt:          result.CreatedAt,
			UpdatedAt:          result.UpdatedAt,
		})
	}

	return &models.DashboardRecentAuthorizationsResponse{
		List:  authorizations,
		Total: total,
	}, nil
}
