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

type CustomerHandler struct {
	customerService service.CustomerService
}

// NewCustomerHandler 创建客户处理器
func NewCustomerHandler(customerService service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

// GetCustomerList 查询客户列表
// @Summary 查询客户列表
// @Description 分页查询客户列表，支持搜索和筛选
// @Tags 客户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认1" minimum(1)
// @Param page_size query int false "每页条数，默认20，最大100" minimum(1) maximum(100)
// @Param search query string false "搜索关键词(支持客户编码、名称、联系人、邮箱)"
// @Param customer_type query string false "客户类型筛选" Enums(individual, enterprise, government, education)
// @Param customer_level query string false "客户等级筛选" Enums(normal, vip, enterprise, strategic)
// @Param status query string false "状态筛选" Enums(active, disabled)
// @Param sort query string false "排序字段，默认created_at" Enums(created_at, updated_at, customer_name, customer_code)
// @Param order query string false "排序方向，默认desc" Enums(asc, desc)
// @Success 200 {object} models.APIResponse{data=models.CustomerListResponse} "查询成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/customers [get]
func (h *CustomerHandler) GetCustomerList(c *gin.Context) {
	var req models.CustomerListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
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
	
	data, err := h.customerService.GetCustomerList(ctx, &req)
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