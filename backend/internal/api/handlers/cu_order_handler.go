package handlers

import (
	"net/http"

	"license-manager/internal/api/middleware"
	"license-manager/internal/models"
	"license-manager/internal/service"
	"license-manager/pkg/i18n"
	"license-manager/pkg/utils"

	"github.com/gin-gonic/gin"
)

type CuOrderHandler struct {
	cuOrderService service.CuOrderService
	paymentService service.PaymentService
	packageService service.PackageService
}

func NewCuOrderHandler(cuOrderService service.CuOrderService, paymentService service.PaymentService, packageService service.PackageService) *CuOrderHandler {
	return &CuOrderHandler{
		cuOrderService: cuOrderService,
		paymentService: paymentService,
		packageService: packageService,
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

	// 从数据库获取套餐列表
	packages, err := h.packageService.GetCuPackageList(c.Request.Context())
	if err != nil {
		status, errCode, message := i18n.NewI18nErrorResponse("800005", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 转换为map数组格式（兼容原有结构）
	result := make([]map[string]interface{}, len(packages))
	for i, pkg := range packages {
		result[i] = map[string]interface{}{
			"id":          pkg.ID,
			"name":        pkg.Name,
			"type":        pkg.Type,
			"price":       pkg.Price,
			"max_devices": pkg.MaxDevices,
			"description": pkg.Description,
			"features":    pkg.Features,
			"details":     pkg.Details,
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      result,
		Timestamp: getCurrentTimestamp(),
	})
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

// CreateOrder 创建订单（支持支付）
// @Summary 创建订单
// @Description 创建新的产品套餐订单，支持免费和付费模式
// @Tags 客户订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body models.CuOrderCreateRequest true "订单创建信息"
// @Success 200 {object} models.APIResponse{data=object{id=string,order_no=string,total_amount=number,authorization_code=string}} "免费订单成功"
// @Success 200 {object} models.APIResponse{data=object{id=string,order_no=string,payment_no=string,total_amount=number,payment_url=string,expire_time=string}} "付费订单成功"
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

	// 根据套餐类型和支付方式决定处理逻辑
	if req.PackageID == "trial" {
		// 试用版订单：直接生成授权码（每月1-25日内）
		h.createFreeOrder(c, cuUserID.(string), customerID.(string), &req)
	} else if req.PaymentMethod == "" {
		// 免费订单：直接创建订单
		h.createFreeOrder(c, cuUserID.(string), customerID.(string), &req)
	} else {
		// 付费订单：创建订单和支付单
		h.createPaidOrder(c, cuUserID.(string), customerID.(string), &req)
	}
}

// createFreeOrder 创建免费订单
func (h *CuOrderHandler) createFreeOrder(c *gin.Context, cuUserID, customerID string, req *models.CuOrderCreateRequest) {
	order, err := h.cuOrderService.CreateOrder(c.Request.Context(), cuUserID, customerID, req)
	if err != nil {
		c.JSON(err.(*i18n.I18nError).HttpCode, models.ErrorResponse{
			Code:      err.(*i18n.I18nError).Code,
			Message:   err.(*i18n.I18nError).Message,
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

// createPaidOrder 创建付费订单
func (h *CuOrderHandler) createPaidOrder(c *gin.Context, cuUserID, customerID string, req *models.CuOrderCreateRequest) {
	// 客户用户创建订单，不需要额外的管理员权限检查
	// cuUserID 和 customerID 已经在上层方法中验证过了

	// 计算价格
	priceResult, err := h.cuOrderService.CalculatePrice(c.Request.Context(), req.PackageID, req.LicenseCount)
	if err != nil {
		c.JSON(err.(*i18n.I18nError).HttpCode, models.ErrorResponse{
			Code:      err.(*i18n.I18nError).Code,
			Message:   err.(*i18n.I18nError).Message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 创建业务订单（pending状态）
	order, err := h.cuOrderService.CreatePendingOrder(c.Request.Context(), cuUserID, customerID, req)
	if err != nil {
		c.JSON(err.(*i18n.I18nError).HttpCode, models.ErrorResponse{
			Code:      err.(*i18n.I18nError).Code,
			Message:   err.(*i18n.I18nError).Message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 创建支付单
	paymentReq := &models.PaymentCreateRequest{
		BusinessType:  models.BusinessTypePackageOrder,
		BusinessID:    &order.ID,
		CustomerID:    customerID,
		CuUserID:      cuUserID,
		Amount:        priceResult.TotalAmount,
		Currency:      "CNY",
		PaymentMethod: req.PaymentMethod,
	}

	payment, err := h.paymentService.CreatePayment(c.Request.Context(), paymentReq)
	if err != nil {
		// 如果支付单创建失败，需要回滚业务订单
		// 这里应该有事务处理，暂时先返回错误
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:      "602001",
			Message:   "支付单创建失败",
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
			"id":           order.ID,
			"order_no":     order.OrderNo,
			"payment_no":   payment.PaymentNo,
			"total_amount": priceResult.TotalAmount,
			"payment_url":  payment.PaymentURL,
			"expire_time":  payment.ExpireTime.Format("2006-01-02T15:04:05Z"),
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

// CancelOrder 删除订单
// @Summary 删除订单
// @Description 删除当前用户的订单记录
// @Tags 客户订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_id path string true "订单ID"
// @Success 200 {object} models.APIResponse{} "成功"
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Failure 409
// @Router /api/cu/orders/{order_id}/cancel [put]
func (h *CuOrderHandler) CancelOrder(c *gin.Context) {
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

	ctx := middleware.WithLanguage(c.Request.Context(), c)
	err := h.cuOrderService.DeleteOrder(ctx, orderID, cuUserID.(string))
	if err != nil {
		c.JSON(err.(*i18n.I18nError).HttpCode, models.ErrorResponse{
			Code:      err.(*i18n.I18nError).Code,
			Message:   err.(*i18n.I18nError).Message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	lang := middleware.GetLanguage(c)
	successMsg := i18n.GetI18nErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   successMsg,
		Data:      nil,
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
// @Param search query string false "订单号或授权码模糊匹配"
// @Param status query string false "订单状态筛选" Enums(pending, paid, cancelled)
// @Param time query string false "时间筛选" Enums(today, week, month, three_months)
// @Success 200 {object} models.APIResponse{data=models.CuOrderListResponse} "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Router /api/cu/orders [get]
func (h *CuOrderHandler) GetUserOrders(c *gin.Context) {
	var req models.CuOrderListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
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

	// 获取订单列表
	ctx := middleware.WithLanguage(c.Request.Context(), c)
	orders, total, err := h.cuOrderService.GetUserOrders(ctx, cuUserID.(string), &req)
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
			Page:       req.Page,
			PageSize:   req.PageSize,
		},
		Timestamp: getCurrentTimestamp(),
	})
}

// GetOrderSummary 获取订单汇总统计
// @Summary 获取订单汇总统计
// @Description 获取当前登录用户的订单汇总统计信息（总数、待支付数、已支付数）
// @Tags 客户订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.OrderSummaryResponse} "获取成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/orders/summary [get]
func (h *CuOrderHandler) GetOrderSummary(c *gin.Context) {
	// 从JWT Token中获取用户ID和客户ID
	claims := c.MustGet("cu_user").(*utils.CuClaims)

	// 调用服务层
	result, err := h.cuOrderService.GetOrderSummary(c.Request.Context(), claims.CustomerID)
	if err != nil {
		// err已经是i18n.I18nError，直接使用
		i18nErr, ok := err.(*i18n.I18nError)
		if ok {
			c.JSON(i18nErr.HttpCode, models.ErrorResponse{
				Code:      i18nErr.Code,
				Message:   i18nErr.Message,
				Timestamp: getCurrentTimestamp(),
			})
			return
		}
		// 如果不是i18n错误，使用通用错误
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900004", lang, err.Error())
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	lang := middleware.GetLanguage(c)
	successMessage := i18n.GetI18nErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMessage,
		Data:    result,
	})
}

// ContinuePay 继续支付
// @Summary 继续支付
// @Description 对指定订单继续支付，如果支付单已过期则创建新的支付单
// @Tags 客户订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_id path string true "订单ID"
// @Param payment_method query string false "支付方式" Enums(alipay, wechat)
// @Success 200 {object} models.APIResponse{data=object{order_id=string,order_no=string,payment_no=string,payment_url=string,total_amount=number,expire_time=string}} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 404 {object} models.ErrorResponse "订单不存在"
// @Failure 400 {object} models.ErrorResponse "订单状态不允许继续支付"
// @Router /api/cu/orders/{order_id}/pay [post]
func (h *CuOrderHandler) ContinuePay(c *gin.Context) {
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

	// 获取支付方式（可选，默认使用配置的默认支付方式）
	paymentMethod := c.DefaultQuery("payment_method", "")

	// 调用服务层
	ctx := middleware.WithLanguage(c.Request.Context(), c)
	result, err := h.cuOrderService.ContinuePay(ctx, orderID, cuUserID.(string), customerID.(string), paymentMethod)
	if err != nil {
		c.JSON(err.(*i18n.I18nError).HttpCode, models.ErrorResponse{
			Code:      err.(*i18n.I18nError).Code,
			Message:   err.(*i18n.I18nError).Message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 如果支付单号为空，说明需要创建新的支付单
	if result.PaymentNo == "" {
		// 获取订单信息以获取金额
		order, err := h.cuOrderService.GetOrder(ctx, orderID, cuUserID.(string))
		if err != nil {
			c.JSON(err.(*i18n.I18nError).HttpCode, models.ErrorResponse{
				Code:      err.(*i18n.I18nError).Code,
				Message:   err.(*i18n.I18nError).Message,
				Timestamp: getCurrentTimestamp(),
			})
			return
		}

		// 创建新的支付单
		paymentReq := &models.PaymentCreateRequest{
			BusinessType:  models.BusinessTypePackageOrder,
			BusinessID:    &order.ID,
			CustomerID:    customerID.(string),
			CuUserID:      cuUserID.(string),
			Amount:        order.TotalAmount,
			Currency:      "CNY",
			PaymentMethod: paymentMethod,
		}

		payment, err := h.paymentService.CreatePayment(ctx, paymentReq)
		if err != nil {
			lang := middleware.GetLanguage(c)
			status, errCode, message := i18n.NewI18nErrorResponse("602001", lang)
			c.JSON(status, models.ErrorResponse{
				Code:      errCode,
				Message:   message + ": " + err.Error(),
				Timestamp: getCurrentTimestamp(),
			})
			return
		}

		// 更新返回结果
		result.PaymentNo = payment.PaymentNo
		result.PaymentURL = payment.PaymentURL
		result.ExpireTime = payment.ExpireTime
	}

	lang := middleware.GetLanguage(c)
	successMsg := i18n.GetI18nErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    "000000",
		Message: successMsg,
		Data: gin.H{
			"order_id":     result.OrderID,
			"order_no":     result.OrderNo,
			"payment_no":   result.PaymentNo,
			"payment_url":  result.PaymentURL,
			"total_amount": result.TotalAmount,
			"expire_time":  result.ExpireTime.Format("2006-01-02T15:04:05Z"),
		},
		Timestamp: getCurrentTimestamp(),
	})
}