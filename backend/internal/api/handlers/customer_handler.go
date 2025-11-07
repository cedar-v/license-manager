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
			Message:   message + ": " + err.Error(),
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
		Code:    "000000",
		Message: successMessage,
		Data:    data,
	})
}

// GetCustomer 获取客户详情
// @Summary 获取客户详情
// @Description 根据客户ID获取客户详细信息
// @Tags 客户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "客户ID"
// @Success 200 {object} models.APIResponse{data=models.Customer} "查询成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 404 {object} models.ErrorResponse "客户不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/customers/{id} [get]
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
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

	data, err := h.customerService.GetCustomer(ctx, id)
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

// CreateCustomer 创建客户
// @Summary 创建客户
// @Description 创建新的客户记录，自动生成客户编码
// @Tags 客户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param customer body models.CustomerCreateRequest true "客户信息"
// @Success 200 {object} models.APIResponse{data=models.Customer} "创建成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 409 {object} models.ErrorResponse "客户已存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/customers [post]
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req models.CustomerCreateRequest
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

	data, err := h.customerService.CreateCustomer(ctx, &req)
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

// UpdateCustomer 更新客户信息
// @Summary 更新客户信息
// @Description 更新客户信息，所有字段都是可选的
// @Tags 客户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "客户ID"
// @Param customer body models.CustomerUpdateRequest true "客户信息"
// @Success 200 {object} models.APIResponse{data=models.Customer} "更新成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 404 {object} models.ErrorResponse "客户不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	var req models.CustomerUpdateRequest
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
	
	data, err := h.customerService.UpdateCustomer(ctx, id, &req)
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

// DeleteCustomer 删除客户
// @Summary 删除客户
// @Description 直接删除客户记录，删除前会检查是否有关联的授权码或许可证，如有则提示先删除授权
// @Tags 客户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "客户ID"
// @Success 200 {object} models.APIResponse "删除成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 404 {object} models.ErrorResponse "客户不存在"
// @Failure 409 {object} models.ErrorResponse "客户仍有授权，无法删除"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
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
	
	err := h.customerService.DeleteCustomer(ctx, id)
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
		Data:    nil,
	})
}

// UpdateCustomerStatus 修改客户状态
// @Summary 修改客户状态
// @Description 修改客户状态（激活/停用）
// @Tags 客户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "客户ID"
// @Param status body models.CustomerStatusUpdateRequest true "状态信息"
// @Success 200 {object} models.APIResponse{data=models.Customer} "状态更新成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 404 {object} models.ErrorResponse "客户不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/customers/{id}/status [patch]
func (h *CustomerHandler) UpdateCustomerStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	var req models.CustomerStatusUpdateRequest
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
	
	data, err := h.customerService.UpdateCustomerStatus(ctx, id, &req)
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
