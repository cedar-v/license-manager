package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"license-manager/internal/api/middleware"
	"license-manager/internal/models"
	"license-manager/internal/service"
	"license-manager/pkg/i18n"
)

type DashboardHandler struct {
	dashboardService service.DashboardService
}

// NewDashboardHandler 创建仪表盘处理器
func NewDashboardHandler(dashboardService service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetAuthorizationTrend 获取授权趋势数据
// @Summary 获取授权趋势数据
// @Description 根据时间类型（本周/本月/自定义）获取授权趋势统计数据
// @Tags 仪表盘
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string true "时间类型" Enums(week,month,custom)
// @Param start_date query string false "开始日期(YYYY-MM-DD格式，当type为custom时必填)"
// @Param end_date query string false "结束日期(YYYY-MM-DD格式，当type为custom时必填)"
// @Param timezone query string false "时区(如:Asia/Shanghai,UTC等，默认使用服务器本地时区)"
// @Success 200 {object} models.APIResponse{data=models.DashboardAuthorizationTrendResponse} "授权趋势数据"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/dashboard/authorization-trend [get]
func (h *DashboardHandler) GetAuthorizationTrend(c *gin.Context) {
	// 绑定查询参数
	var req models.DashboardAuthorizationTrendRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	// 设置语言到Context中
	ctx := middleware.WithLanguage(c.Request.Context(), c)

	// 调用服务层获取数据
	response, err := h.dashboardService.GetAuthorizationTrend(ctx, &req)
	if err != nil {
		// 错误已经在Service层完全包装好了，直接使用
		var i18nErr *i18n.I18nError
		if errors.As(err, &i18nErr) {
			c.JSON(i18nErr.HttpCode, models.ErrorResponse{
				Code:      i18nErr.Code,
				Message:   i18nErr.Message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		} else {
			// 兜底：理论上不应该到这里，因为Service层应该返回I18nError
			lang := middleware.GetLanguage(c)
			status, errCode, message := i18n.NewI18nErrorResponse("900004", lang)
			c.JSON(status, models.ErrorResponse{
				Code:      errCode,
				Message:   message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		}
		return
	}

	// 返回成功响应
	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    response,
	})
}

// GetRecentAuthorizations 获取最近授权列表
// @Summary 获取最近授权列表
// @Description 获取最近创建的授权码列表，支持按客户ID和状态筛选
// @Tags 仪表盘
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "返回数量限制（默认20，最大100）"
// @Param customer_id query string false "客户ID筛选"
// @Param status query string false "状态筛选" Enums(normal,locked,expired)
// @Success 200 {object} models.APIResponse{data=models.DashboardRecentAuthorizationsResponse} "最近授权列表"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/dashboard/recent-authorizations [get]
func (h *DashboardHandler) GetRecentAuthorizations(c *gin.Context) {
	// 绑定查询参数
	var req models.DashboardRecentAuthorizationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	// 设置语言到Context中
	ctx := middleware.WithLanguage(c.Request.Context(), c)

	// 调用服务层获取数据
	response, err := h.dashboardService.GetRecentAuthorizations(ctx, &req)
	if err != nil {
		// 错误已经在Service层完全包装好了，直接使用
		var i18nErr *i18n.I18nError
		if errors.As(err, &i18nErr) {
			c.JSON(i18nErr.HttpCode, models.ErrorResponse{
				Code:      i18nErr.Code,
				Message:   i18nErr.Message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		} else {
			// 兜底：理论上不应该到这里，因为Service层应该返回I18nError
			lang := middleware.GetLanguage(c)
			status, errCode, message := i18n.NewI18nErrorResponse("900004", lang)
			c.JSON(status, models.ErrorResponse{
				Code:      errCode,
				Message:   message,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		}
		return
	}

	// 返回成功响应
	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    response,
	})
}