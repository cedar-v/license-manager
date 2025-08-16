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

type EnumHandler struct {
	enumService service.EnumService
}

// NewEnumHandler 创建枚举处理器
func NewEnumHandler(enumService service.EnumService) *EnumHandler {
	return &EnumHandler{
		enumService: enumService,
	}
}

// GetAllEnums 获取所有枚举类型及其值
// @Summary 获取所有枚举值
// @Description 获取系统中所有枚举类型及其对应的多语言显示值
// @Tags 枚举管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.EnumListResponse} "查询成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/enums [get]
func (h *EnumHandler) GetAllEnums(c *gin.Context) {
	// 设置语言到Context中
	ctx := middleware.WithLanguage(c.Request.Context(), c)
	
	data, err := h.enumService.GetAllEnums(ctx)
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
		Code:    0,
		Message: successMessage,
		Data:    data,
	})
}

// GetEnumsByType 获取指定类型的枚举值
// @Summary 获取指定类型的枚举值
// @Description 根据枚举类型获取对应的多语言显示值
// @Tags 枚举管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type path string true "枚举类型" Enums(customer_type, customer_level, customer_status, company_size)
// @Success 200 {object} models.APIResponse{data=models.EnumTypeResponse} "查询成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/enums/{type} [get]
func (h *EnumHandler) GetEnumsByType(c *gin.Context) {
	enumType := c.Param("type")
	if enumType == "" {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	// 设置语言到Context中
	ctx := middleware.WithLanguage(c.Request.Context(), c)
	
	data, err := h.enumService.GetEnumsByType(ctx, enumType)
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
		Code:    0,
		Message: successMessage,
		Data:    data,
	})
}