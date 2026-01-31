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

type PackageHandler struct {
	packageService service.PackageService
}

func NewPackageHandler(packageService service.PackageService) *PackageHandler {
	return &PackageHandler{
		packageService: packageService,
	}
}

// GetPackages 获取套餐列表
// @Summary 获取套餐列表
// @Description 管理员获取套餐列表，支持类型和状态筛选
// @Tags 套餐管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string false "套餐类型筛选"
// @Param status query int false "状态筛选，1-启用，0-禁用"
// @Success 200 {object} models.APIResponse{data=models.PackageListResponse} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/packages [get]
func (h *PackageHandler) GetPackages(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	var req models.PackageListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	result, err := h.packageService.GetPackageList(c.Request.Context(), &req)
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      result,
		Timestamp: getCurrentTimestamp(),
	})
}

// GetPackageDetail 获取套餐详情
// @Summary 获取套餐详情
// @Description 根据ID获取套餐详细信息
// @Tags 套餐管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "套餐ID"
// @Success 200 {object} models.APIResponse{data=models.PackageResponse} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 404 {object} models.ErrorResponse "套餐不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/packages/{id} [get]
func (h *PackageHandler) GetPackageDetail(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	id := c.Param("id")

	pkg, err := h.packageService.GetPackageDetail(c.Request.Context(), id)
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      pkg.ToResponse(),
		Timestamp: getCurrentTimestamp(),
	})
}

// CreatePackage 创建套餐
// @Summary 创建套餐
// @Description 创建新的套餐
// @Tags 套餐管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.PackageCreateRequest true "套餐创建请求"
// @Success 200 {object} models.APIResponse{data=models.PackageResponse} "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/packages [post]
func (h *PackageHandler) CreatePackage(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	var req models.PackageCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	pkg, err := h.packageService.CreatePackage(c.Request.Context(), &req)
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      pkg.ToResponse(),
		Timestamp: getCurrentTimestamp(),
	})
}

// UpdatePackage 更新套餐
// @Summary 更新套餐
// @Description 更新指定套餐的信息
// @Tags 套餐管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "套餐ID"
// @Param request body models.PackageUpdateRequest true "套餐更新请求"
// @Success 200 {object} models.APIResponse{data=models.PackageResponse} "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 404 {object} models.ErrorResponse "套餐不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/packages/{id} [put]
func (h *PackageHandler) UpdatePackage(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	id := c.Param("id")

	var req models.PackageUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message + ": " + err.Error(),
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	pkg, err := h.packageService.UpdatePackage(c.Request.Context(), id, &req)
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      pkg.ToResponse(),
		Timestamp: getCurrentTimestamp(),
	})
}

// DeletePackage 删除套餐
// @Summary 删除套餐
// @Description 删除指定套餐（软删除）
// @Tags 套餐管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "套餐ID"
// @Success 200 {object} models.APIResponse{} "成功"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 404 {object} models.ErrorResponse "套餐不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/packages/{id} [delete]
func (h *PackageHandler) DeletePackage(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	id := c.Param("id")

	if err := h.packageService.DeletePackage(c.Request.Context(), id); err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      nil,
		Timestamp: getCurrentTimestamp(),
	})
}

// UpdatePackageStatus 更新套餐状态
// @Summary 更新套餐状态
// @Description 启用或禁用套餐
// @Tags 套餐管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "套餐ID"
// @Param status query int true "状态，1-启用，0-禁用"
// @Success 200 {object} models.APIResponse{data=models.PackageResponse} "成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "权限不足"
// @Failure 404 {object} models.ErrorResponse "套餐不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/packages/{id}/status [put]
func (h *PackageHandler) UpdatePackageStatus(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	id := c.Param("id")

	statusStr := c.Query("status")
	if statusStr == "" {
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	status, err := strconv.Atoi(statusStr)
	if err != nil || (status != 0 && status != 1) {
		statusCode, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(statusCode, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	req := &models.PackageUpdateRequest{
		Status: &status,
	}

	pkg, err := h.packageService.UpdatePackage(c.Request.Context(), id, req)
	if err != nil {
		handleI18nError(c, err, lang)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:      "000000",
		Message:   i18n.GetI18nErrorMessage("000000", lang),
		Data:      pkg.ToResponse(),
		Timestamp: getCurrentTimestamp(),
	})
}
