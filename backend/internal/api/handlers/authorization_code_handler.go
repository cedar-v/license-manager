package handlers

import (
	"errors"
	"net/http"
	"time"

	"license-manager/internal/api/middleware"
	"license-manager/internal/models"
	"license-manager/internal/service"
	"license-manager/pkg/i18n"

	"github.com/gin-gonic/gin"
)

type AuthorizationCodeHandler struct {
	authCodeService service.AuthorizationCodeService
}

// NewAuthorizationCodeHandler 创建授权码处理器
func NewAuthorizationCodeHandler(authCodeService service.AuthorizationCodeService) *AuthorizationCodeHandler {
	return &AuthorizationCodeHandler{
		authCodeService: authCodeService,
	}
}

// CreateAuthorizationCode 创建授权码
// @Summary 创建授权码
// @Description 创建新的授权码，自动生成授权码字符串。有效期从当天00:00:00开始，到指定天数后的23:59:59结束
// @Tags 授权码管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param authorization_code body models.AuthorizationCodeCreateRequest true "授权码信息"
// @Success 200 {object} models.APIResponse{data=models.AuthorizationCodeCreateResponse} "创建成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 404 {object} models.ErrorResponse "客户不存在"
// @Failure 409 {object} models.ErrorResponse "授权码已存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/authorization-codes [post]
func (h *AuthorizationCodeHandler) CreateAuthorizationCode(c *gin.Context) {
	var req models.AuthorizationCodeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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

	data, err := h.authCodeService.CreateAuthorizationCode(ctx, &req)
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

	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    data,
	})
}

// GetAuthorizationCodeList 查询授权码列表
// @Summary 查询授权码列表
// @Description 分页查询授权码列表，支持筛选和排序
// @Tags 授权码管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认1" minimum(1)
// @Param page_size query int false "每页条数，默认20，最大100" minimum(1) maximum(100)
// @Param customer_id query string false "客户ID筛选"
// @Param status query string false "状态筛选" Enums(normal, locked, expired)
// @Param start_date query string false "创建开始时间"
// @Param end_date query string false "创建结束时间"
// @Param sort query string false "排序字段，默认created_at" Enums(created_at, updated_at, code)
// @Param order query string false "排序方向，默认desc" Enums(asc, desc)
// @Success 200 {object} models.APIResponse{data=models.AuthorizationCodeListResponse} "查询成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/authorization-codes [get]
func (h *AuthorizationCodeHandler) GetAuthorizationCodeList(c *gin.Context) {
	var req models.AuthorizationCodeListRequest
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

	data, err := h.authCodeService.GetAuthorizationCodeList(ctx, &req)
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

	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    data,
	})
}