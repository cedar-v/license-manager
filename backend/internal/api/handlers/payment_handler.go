package handlers

import (
	"net/http"
	"strconv"

	"license-manager/internal/api/middleware"
	"license-manager/internal/models"
	"license-manager/internal/service"
	"license-manager/pkg/i18n"
	"license-manager/pkg/utils"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// GetPaymentStatus 获取支付状态
// @Summary 获取支付状态
// @Description 根据支付单号获取支付状态
// @Tags 支付管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payment_no path string true "支付单号"
// @Success 200 {object} models.APIResponse{data=models.PaymentStatusResponse} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 404 {object} models.ErrorResponse "支付单不存在"
// @Router /api/payment/{payment_no}/status [get]
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	paymentNo := c.Param("payment_no")
	if paymentNo == "" {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	status, err := h.paymentService.GetPaymentStatus(c.Request.Context(), paymentNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:      "602001",
			Message:   "获取支付状态失败",
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	lang := middleware.GetLanguage(c)
	successMsg := i18n.GetI18nErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   successMsg,
		Data:      status,
		Timestamp: getCurrentTimestamp(),
	})
}

// GetUserPayments 获取用户支付历史
// @Summary 获取用户支付历史
// @Description 获取当前用户的支付记录列表
// @Tags 支付管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query integer false "页码" default(1)
// @Param page_size query integer false "每页数量" default(10)
// @Success 200 {object} models.APIResponse{data=models.PaymentListResponse} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Router /api/payment/history [get]
func (h *PaymentHandler) GetUserPayments(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 获取客户用户ID和客户ID
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

	payments, err := h.paymentService.GetUserPayments(c.Request.Context(), customerID.(string), cuUserID.(string), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:      "602001",
			Message:   "获取支付历史失败",
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	lang := middleware.GetLanguage(c)
	successMsg := i18n.GetI18nErrorMessage("000000", lang)
	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   successMsg,
		Data:      payments,
		Timestamp: getCurrentTimestamp(),
	})
}

// AlipayCallback 支付宝回调处理
// @Summary 支付宝回调
// @Description 处理支付宝支付回调
// @Tags 支付管理
// @Accept application/x-www-form-urlencoded
// @Produce text/plain
// @Param notification body string false "支付宝回调参数"
// @Success 200 {string} string "success"
// @Failure 400 {string} string "failure"
// @Router /api/payment/alipay/callback [post]
func (h *PaymentHandler) AlipayCallback(c *gin.Context) {
	// 这里我们需要一个alipay client实例来解析回调
	// 由于我们在service层已经有了client，这里暂时简化处理
	// 实际应该从service中获取client

	// 暂时使用简化的处理方式
	// 实际项目中应该使用service中的client来验证和处理
	utils.SendNotificationResponse(c.Writer, true)
}
