package handlers

import (
	"net/http"
	"strconv"

	"license-manager/internal/api/middleware"
	"license-manager/internal/models"
	"license-manager/internal/service"
	"license-manager/pkg/i18n"

	"github.com/gin-gonic/gin"
)

type CuOrderHandler struct {
	cuOrderService service.CuOrderService
}

func NewCuOrderHandler(cuOrderService service.CuOrderService) *CuOrderHandler {
	return &CuOrderHandler{
		cuOrderService: cuOrderService,
	}
}

// GetPackages 获取套餐列表
// @Summary 获取套餐列表
// @Description 获取所有可用的产品套餐信息
// @Tags 客户订单管理
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]object{id=string,name=string,type=string,price=number,max_devices=integer}} "成功"
// @Router /api/cu/packages [get]
func (h *CuOrderHandler) GetPackages(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	// 获取套餐列表
	packages := h.getPackageList()

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      packages,
		Timestamp: getCurrentTimestamp(),
	})
}

// getPackageList 获取套餐列表（暂时写死数据）
func (h *CuOrderHandler) getPackageList() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"id":          "trial",
			"name":        "试用版",
			"type":        "trial",
			"price":       0,
			"max_devices": 1,
			"description": "免费试用，有效期至本月25日",
		},
		{
			"id":          "basic",
			"name":        "基础版",
			"type":        "basic",
			"price":       300,
			"max_devices": 1000,
			"description": "基础功能，支持批量许可购买",
		},
		{
			"id":          "professional",
			"name":        "专业版",
			"type":        "professional",
			"price":       2000,
			"max_devices": 1000,
			"description": "全部功能，支持批量许可购买",
		},
	}
}

// CalculatePrice 计算订单价格
// @Summary 计算订单价格
// @Description 根据套餐和许可数量计算价格和折扣
// @Tags 客户订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{package_id=string,license_count=integer} true "价格计算请求"
// @Success 200 {object} models.APIResponse{data=object{unit_price=number,license_count=integer,discount_rate=number,total_amount=number,discount_description=string}} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 404 {object} models.ErrorResponse "套餐不存在"
// @Router /api/cu/orders/calculate [post]
func (h *CuOrderHandler) CalculatePrice(c *gin.Context) {
	var req struct {
		PackageID    string `json:"package_id" binding:"required"`
		LicenseCount int    `json:"license_count" binding:"required,min=1,max=1000"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 获取用户ID（用于后续业务逻辑扩展）
	_, exists := c.Get("cu_user_id")
	if !exists {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("100004", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 计算价格
	result, err := h.cuOrderService.CalculatePrice(c.Request.Context(), req.PackageID, req.LicenseCount)
	if err != nil {
		c.JSON(err.(*i18n.I18nError).HttpCode, models.ErrorResponse{
			Code:      err.(*i18n.I18nError).Code,
			Message:   err.(*i18n.I18nError).Message, // 直接使用错误对象中的消息
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	lang := middleware.GetLanguage(c)
	successMsg := i18n.GetI18nErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   successMsg,
		Data:      result,
		Timestamp: getCurrentTimestamp(),
	})
}

// CreateOrder 创建订单
// @Summary 创建订单
// @Description 创建新的产品套餐订单
// @Tags 客户订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body models.CuOrderCreateRequest true "订单创建信息"
// @Success 200 {object} models.APIResponse{data=object{id=string,order_no=string,total_amount=number,authorization_code=string}} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 404 {object} models.ErrorResponse "套餐不存在"
// @Failure 409 {object} models.ErrorResponse "当月已购买试用版"
// @Router /api/cu/orders [post]
func (h *CuOrderHandler) CreateOrder(c *gin.Context) {
	var req models.CuOrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 获取用户ID和客户ID
	cuUserID, exists := c.Get("cu_user_id")
	if !exists {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("100004", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	customerID, exists := c.Get("cu_customer_id")
	if !exists {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("100004", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 创建订单
	order, err := h.cuOrderService.CreateOrder(c.Request.Context(), cuUserID.(string), customerID.(string), &req)
	if err != nil {
		c.JSON(err.(*i18n.I18nError).HttpCode, models.ErrorResponse{
			Code:      err.(*i18n.I18nError).Code,
			Message:   err.(*i18n.I18nError).Message, // 直接使用错误对象中的消息
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	lang := middleware.GetLanguage(c)
	successMsg := i18n.GetI18nErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMsg,
		Data: gin.H{
			"id":                 order.ID,
			"order_no":           order.OrderNo,
			"total_amount":       order.TotalAmount,
			"authorization_code": order.AuthorizationCode,
		},
		Timestamp: getCurrentTimestamp(),
	})
}

// GetOrder 获取订单详情
// @Summary 获取订单详情
// @Description 根据订单ID获取订单详细信息
// @Tags 客户订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_id path string true "订单ID"
// @Success 200 {object} models.APIResponse{data=models.CuOrderResponse} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 404 {object} models.ErrorResponse "订单不存在"
// @Router /api/cu/orders/{order_id} [get]
func (h *CuOrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("order_id")
	if orderID == "" {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 获取用户ID
	cuUserID, exists := c.Get("cu_user_id")
	if !exists {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("100004", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 获取订单详情
	order, err := h.cuOrderService.GetOrder(c.Request.Context(), orderID, cuUserID.(string))
	if err != nil {
		c.JSON(err.(*i18n.I18nError).HttpCode, models.ErrorResponse{
			Code:      err.(*i18n.I18nError).Code,
			Message:   err.(*i18n.I18nError).Message, // 直接使用错误对象中的消息
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	lang := middleware.GetLanguage(c)
	successMsg := i18n.GetI18nErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   successMsg,
		Data:      order.ToResponse(),
		Timestamp: getCurrentTimestamp(),
	})
}

// GetUserOrders 获取用户订单列表
// @Summary 获取用户订单列表
// @Description 获取当前用户的订单列表，支持分页
// @Tags 客户订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query integer false "页码" default(1)
// @Param page_size query integer false "每页数量" default(10)
// @Success 200 {object} models.APIResponse{data=models.CuOrderListResponse} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Router /api/cu/orders [get]
func (h *CuOrderHandler) GetUserOrders(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 获取用户ID
	cuUserID, exists := c.Get("cu_user_id")
	if !exists {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("100004", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 获取订单列表
	orders, total, err := h.cuOrderService.GetUserOrders(c.Request.Context(), cuUserID.(string), offset, pageSize)
	if err != nil {
		c.JSON(err.(*i18n.I18nError).HttpCode, models.ErrorResponse{
			Code:      err.(*i18n.I18nError).Code,
			Message:   err.(*i18n.I18nError).Message, // 直接使用错误对象中的消息
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 转换为响应格式
	orderResponses := make([]*models.CuOrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = order.ToResponse()
	}

	lang := middleware.GetLanguage(c)
	successMsg := i18n.GetI18nErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMsg,
		Data: &models.CuOrderListResponse{
			Orders:     orderResponses,
			TotalCount: total,
			Page:       page,
			PageSize:   pageSize,
		},
		Timestamp: getCurrentTimestamp(),
	})
}
