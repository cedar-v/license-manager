package service

import (
	"context"
	"fmt"
	"time"

	"license-manager/internal/models"
	"license-manager/internal/repository"
	pkgcontext "license-manager/pkg/context"
	"license-manager/pkg/i18n"
)

type dashboardService struct {
	dashboardRepo repository.DashboardRepository
}

// NewDashboardService 创建仪表盘服务实例
func NewDashboardService(dashboardRepo repository.DashboardRepository) DashboardService {
	return &dashboardService{
		dashboardRepo: dashboardRepo,
	}
}

// GetAuthorizationTrend 获取授权趋势数据
func (s *dashboardService) GetAuthorizationTrend(ctx context.Context, req *models.DashboardAuthorizationTrendRequest) (*models.DashboardAuthorizationTrendResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	var startDate, endDate time.Time
	var descriptionDisplay string
	
	// 获取时区，默认使用本地时区
	loc := time.Local
	if req.Timezone != "" {
		var err error
		loc, err = time.LoadLocation(req.Timezone)
		if err != nil {
			// 时区解析失败，使用本地时区
			loc = time.Local
		}
	}
	
	now := time.Now().In(loc)

	// 根据类型计算时间范围
	switch req.Type {
	case "week":
		// 本周：从周一到周日
		weekday := int(now.Weekday())
		if weekday == 0 { // Sunday
			weekday = 7
		}
		startDate = now.AddDate(0, 0, -(weekday-1))
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
		endDate = startDate.AddDate(0, 0, 6)
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, loc)
		descriptionDisplay = i18n.GetEnumMessage("dashboard_period", "week", lang)

	case "month":
		// 本月：从1号到月末
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		endDate = startDate.AddDate(0, 1, -1)
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, loc)
		descriptionDisplay = i18n.GetEnumMessage("dashboard_period", "month", lang)

	case "custom":
		// 自定义时间范围
		if req.StartDate == "" || req.EndDate == "" {
			return nil, i18n.NewI18nError("400002", lang) // 开始日期和结束日期都是必填的
		}

		var err error
		startDate, err = time.ParseInLocation("2006-01-02", req.StartDate, loc)
		if err != nil {
			return nil, i18n.NewI18nError("400002", lang) // 开始日期格式错误
		}

		endDate, err = time.ParseInLocation("2006-01-02", req.EndDate, loc)
		if err != nil {
			return nil, i18n.NewI18nError("400003", lang) // 结束日期格式错误
		}

		// 验证时间范围
		if startDate.After(endDate) {
			return nil, i18n.NewI18nError("400005", lang) // 开始日期不能晚于结束日期
		}

		// 验证时间跨度不超过一年
		if endDate.Sub(startDate) > 365*24*time.Hour {
			return nil, i18n.NewI18nError("400004", lang) // 时间范围超过限制
		}

		// 设置为全天
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
		descriptionDisplay = fmt.Sprintf("%s - %s", req.StartDate, req.EndDate)

	default:
		return nil, i18n.NewI18nError("400001", lang) // 时间类型参数错误
	}

	// 获取趋势数据
	trendData, err := s.dashboardRepo.GetAuthorizationTrendData(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 计算汇总信息
	var totalCount, newCount, expiredCount int64
	for _, data := range trendData {
		newCount += data.NewAuthorizations
		expiredCount += data.ExpiredAuthorizations
	}

	// 获取期间结束时的总授权数
	if len(trendData) > 0 {
		totalCount = trendData[len(trendData)-1].TotalAuthorizations
	}

	// 计算增长率（本期新增相对于期初的百分比）
	var growthRate float64
	if totalCount > 0 && newCount > 0 {
		periodStartCount := totalCount - newCount
		if periodStartCount > 0 {
			growthRate = float64(newCount) / float64(periodStartCount) * 100
		} else {
			growthRate = 100.0 // 如果期初为0，则增长率为100%
		}
	}

	return &models.DashboardAuthorizationTrendResponse{
		Period: models.TrendPeriod{
			Type:               req.Type,
			StartDate:          startDate.Format("2006-01-02"),
			EndDate:            endDate.Format("2006-01-02"),
			DescriptionDisplay: descriptionDisplay,
		},
		TrendData: trendData,
		Summary: models.TrendSummary{
			TotalCount:   totalCount,
			NewCount:     newCount,
			ExpiredCount: expiredCount,
			GrowthRate:   growthRate,
		},
	}, nil
}

// GetRecentAuthorizations 获取最近授权列表
func (s *dashboardService) GetRecentAuthorizations(ctx context.Context, req *models.DashboardRecentAuthorizationsRequest) (*models.DashboardRecentAuthorizationsResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 验证limit参数
	if req.Limit < 0 || req.Limit > 100 {
		return nil, i18n.NewI18nError("400006", lang) // 数量限制参数错误
	}

	// 验证status参数
	if req.Status != "" && req.Status != "normal" && req.Status != "locked" && req.Status != "expired" {
		return nil, i18n.NewI18nError("400001", lang) // 状态参数错误
	}

	// 调用repository获取数据
	return s.dashboardRepo.GetRecentAuthorizations(ctx, req)
}