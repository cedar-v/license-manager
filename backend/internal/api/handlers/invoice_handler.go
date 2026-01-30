package handlers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"license-manager/internal/api/middleware"
	"license-manager/internal/models"
	"license-manager/internal/service"
	"license-manager/pkg/i18n"

	"github.com/gin-gonic/gin"
)

type AdminInvoiceHandler struct {
	invoiceService service.AdminInvoiceService
}

func NewAdminInvoiceHandler(invoiceService service.AdminInvoiceService) *AdminInvoiceHandler {
	return &AdminInvoiceHandler{
		invoiceService: invoiceService,
	}
}

// GetAdminInvoices 管理员发票列表
// @Summary 管理员发票列表
// @Description 管理员查看所有发票列表，支持按客户筛选
// @Tags 发票管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认1" minimum(1)
// @Param page_size query int false "每页条数，默认10，最大100" minimum(1) maximum(100)
// @Param status query string false "状态筛选，pending-待处理/issued-已开票/rejected-已驳回"
// @Param search query string false "搜索关键词，支持发票号或订单号模糊匹配"
// @Param apply_date query string false "申请日期筛选，格式YYYY-MM-DD"
// @Param customer_id query string false "客户ID筛选"
// @Success 200 {object} models.APIResponse{data=models.InvoiceListResponse} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/invoices [get]
func (h *AdminInvoiceHandler) GetAdminInvoices(c *gin.Context) {
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

	invoices, err := h.invoiceService.GetAdminInvoices(c.Request.Context(), &req)
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

// GetAdminInvoiceDetail 管理员发票详情
// @Summary 管理员发票详情
// @Description 管理员查看发票详细信息，包含关联信息
// @Tags 发票管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "发票ID"
// @Success 200 {object} models.APIResponse{data=models.InvoiceDetailResponse} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 404 {object} models.ErrorResponse "发票不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/invoices/{id} [get]
func (h *AdminInvoiceHandler) GetAdminInvoiceDetail(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	invoiceID := c.Param("id")
	adminID := getUserID(c)

	detail, err := h.invoiceService.GetInvoiceDetail(c.Request.Context(), invoiceID, true, adminID)
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      detail,
		Timestamp: getCurrentTimestamp(),
	})
}

// GetAdminInvoiceSummary 管理员发票汇总信息
// @Summary 管理员发票汇总信息
// @Description 管理员查看全平台的发票统计汇总信息
// @Tags 发票管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.InvoiceSummaryResponse} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/invoices/summary [get]
func (h *AdminInvoiceHandler) GetAdminInvoiceSummary(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	summary, err := h.invoiceService.GetAdminInvoiceSummary(c.Request.Context())
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

// RejectInvoice 管理员发票驳回
// @Summary 管理员发票驳回
// @Description 管理员驳回发票申请，需要填写驳回原因和修改建议
// @Tags 发票管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "发票ID"
// @Param request body models.InvoiceRejectRequest true "驳回请求"
// @Success 200 {object} models.APIResponse{data=models.InvoiceResponse} "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 404 {object} models.ErrorResponse "发票不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/invoices/{id}/reject [post]
func (h *AdminInvoiceHandler) RejectInvoice(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	invoiceID := c.Param("id")
	adminID := getUserID(c)

	var req models.InvoiceRejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:      "900001",
			Message:   i18n.GetI18nErrorMessage("900001", lang),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	invoice, err := h.invoiceService.RejectInvoice(c.Request.Context(), invoiceID, adminID, &req)
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

// IssueInvoice 管理员发票开票
// @Summary 管理员发票开票
// @Description 管理员为发票开票，上传发票文件并设置开票时间
// @Tags 发票管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "发票ID"
// @Param request body models.InvoiceIssueRequest true "开票请求"
// @Success 200 {object} models.APIResponse{data=models.InvoiceResponse} "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 404 {object} models.ErrorResponse "发票不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/invoices/{id}/issue [post]
func (h *AdminInvoiceHandler) IssueInvoice(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	invoiceID := c.Param("id")
	adminID := getUserID(c)

	var req models.InvoiceIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:      "900001",
			Message:   i18n.GetI18nErrorMessage("900001", lang),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	invoice, err := h.invoiceService.IssueInvoice(c.Request.Context(), invoiceID, adminID, &req)
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

// UploadInvoiceFile 管理端发票文件上传
// @Summary 管理端发票文件上传
// @Description 管理员上传发票PDF文件
// @Tags 发票管理
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param invoice_no formData string true "发票申请号"
// @Param file formData file true "发票PDF文件"
// @Success 200 {object} models.APIResponse{data=models.InvoiceUploadResponse} "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/invoices/upload [post]
func (h *AdminInvoiceHandler) UploadInvoiceFile(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	// 获取表单参数
	invoiceNo := c.PostForm("invoice_no")
	if invoiceNo == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:      "900001",
			Message:   i18n.GetI18nErrorMessage("900001", lang),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:      "900001",
			Message:   i18n.GetI18nErrorMessage("900001", lang),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}
	defer file.Close()

	// 验证文件类型
	if filepath.Ext(header.Filename) != ".pdf" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:      "900001",
			Message:   i18n.GetI18nErrorMessage("900001", lang),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 构建文件保存路径
	uploadDir := filepath.Join("..", "files", "invoices")
	filename := invoiceNo + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".pdf"
	filePath := filepath.Join(uploadDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(header, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:      "900004",
			Message:   i18n.GetI18nErrorMessage("900004", lang),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 返回相对URL
	fileURL := "/files/invoices/" + filename
	response := models.InvoiceUploadResponse{
		FileURL: fileURL,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      response,
		Timestamp: getCurrentTimestamp(),
	})
}

// DownloadByToken 公共发票下载（通过token）
// @Summary 公共发票下载
// @Description 通过下载token下载发票文件，无需登录验证
// @Tags 发票管理
// @Accept json
// @Produce application/pdf
// @Param token query string true "下载token"
// @Success 200 {file} binary "发票文件"
// @Failure 400 {object} models.ErrorResponse "token无效"
// @Failure 404 {object} models.ErrorResponse "文件不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/public/invoices/download [get]
func (h *AdminInvoiceHandler) DownloadByToken(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	token := c.Query("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:      "900001",
			Message:   i18n.GetI18nErrorMessage("900001", lang),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	fileURL, err := h.invoiceService.DownloadByToken(c.Request.Context(), token)
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	// 重定向到文件URL
	c.Redirect(http.StatusFound, fileURL)
}

// 辅助函数

// getCuUserID 获取客户用户ID
func getCuUserID(c *gin.Context) string {
	if userID, exists := c.Get("cu_user_id"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// getUserID 获取管理员用户ID
func getUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// handleI18nError 处理国际化错误
func handleI18nError(c *gin.Context, err error, lang string) {
	if i18nErr, ok := err.(*i18n.I18nError); ok {
		c.JSON(i18nErr.HttpCode, models.ErrorResponse{
			Code:      i18nErr.Code,
			Message:   i18nErr.Message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	// 默认错误处理
	c.JSON(http.StatusInternalServerError, models.ErrorResponse{
		Code:      "900004",
		Message:   i18n.GetI18nErrorMessage("900004", lang),
		Timestamp: getCurrentTimestamp(),
	})
}
