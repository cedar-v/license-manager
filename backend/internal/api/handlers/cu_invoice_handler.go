package handlers

import (
	"net/http"

	"license-manager/internal/api/middleware"
	"license-manager/internal/models"
	"license-manager/internal/service"
	"license-manager/pkg/i18n"

	"github.com/gin-gonic/gin"
)

type CuInvoiceHandler struct {
	invoiceService service.CuInvoiceService
}

func NewCuInvoiceHandler(invoiceService service.CuInvoiceService) *CuInvoiceHandler {
	return &CuInvoiceHandler{
		invoiceService: invoiceService,
	}
}

// CreateInvoice 客户申请发票
// @Summary 客户申请发票
// @Description 客户为已支付的订单申请发票
// @Tags 发票管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.InvoiceCreateRequest true "发票申请请求"
// @Success 200 {object} models.APIResponse{data=object{id=string,invoice_no=string,status=string}} "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/invoices [post]
func (h *CuInvoiceHandler) CreateInvoice(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	var req models.InvoiceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:      "900001",
			Message:   i18n.GetI18nErrorMessage("900001", lang),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	invoice, err := h.invoiceService.CreateInvoice(c.Request.Context(), getCuUserID(c), &req)
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	response := invoice.ToResponse()
	// 只返回关键信息
	result := map[string]interface{}{
		"id":         response.ID,
		"invoice_no": response.InvoiceNo,
		"status":     response.Status,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      result,
		Timestamp: getCurrentTimestamp(),
	})
}

// GetUserInvoices 客户发票列表
// @Summary 客户发票列表
// @Description 获取当前客户的发票列表
// @Tags 发票管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认1" minimum(1)
// @Param page_size query int false "每页条数，默认10，最大100" minimum(1) maximum(100)
// @Param status query string false "状态筛选，pending-待处理/issued-已开票/rejected-已驳回"
// @Param search query string false "搜索关键词，支持发票号或订单号模糊匹配"
// @Param apply_date query string false "申请日期筛选，格式YYYY-MM-DD"
// @Success 200 {object} models.APIResponse{data=models.InvoiceListResponse} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/invoices [get]
func (h *CuInvoiceHandler) GetUserInvoices(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	var req models.InvoiceListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:      "900001",
			Message:   i18n.GetI18nErrorMessage("900001", lang),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	invoices, err := h.invoiceService.GetUserInvoices(c.Request.Context(), getCuUserID(c), &req)
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      invoices,
		Timestamp: getCurrentTimestamp(),
	})
}

// GetUserInvoiceDetail 客户发票详情
// @Summary 客户发票详情
// @Description 获取指定发票的详细信息
// @Tags 发票管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "发票ID"
// @Success 200 {object} models.APIResponse{data=models.InvoiceResponse} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 404 {object} models.ErrorResponse "发票不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/invoices/{id} [get]
func (h *CuInvoiceHandler) GetUserInvoiceDetail(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	invoiceID := c.Param("id")

	invoice, err := h.invoiceService.GetInvoice(c.Request.Context(), invoiceID, getCuUserID(c))
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      invoice.ToResponse(),
		Timestamp: getCurrentTimestamp(),
	})
}

// GetUserInvoiceSummary 客户发票汇总信息
// @Summary 客户发票汇总信息
// @Description 获取当前客户的发票统计汇总信息
// @Tags 发票管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.InvoiceSummaryResponse} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/invoices/summary [get]
func (h *CuInvoiceHandler) GetUserInvoiceSummary(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	summary, err := h.invoiceService.GetInvoiceSummary(c.Request.Context(), getCuUserID(c))
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      summary,
		Timestamp: getCurrentTimestamp(),
	})
}

// DownloadInvoice 客户发票下载
// @Summary 客户发票下载
// @Description 下载指定发票的文件
// @Tags 发票管理
// @Accept json
// @Produce application/pdf
// @Security BearerAuth
// @Param id path string true "发票ID"
// @Success 200 {file} binary "发票文件"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 404 {object} models.ErrorResponse "发票不存在或文件不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/invoices/{id}/download [get]
func (h *CuInvoiceHandler) DownloadInvoice(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	invoiceID := c.Param("id")

	fileURL, err := h.invoiceService.DownloadInvoice(c.Request.Context(), invoiceID, getCuUserID(c))
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	// 重定向到文件URL
	c.Redirect(http.StatusFound, fileURL)
}
