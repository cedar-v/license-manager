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

type CuAuthorizationHandler struct {
	authCodeService service.AuthorizationCodeService
}

func NewCuAuthorizationHandler(authCodeService service.AuthorizationCodeService) *CuAuthorizationHandler {
	return &CuAuthorizationHandler{
		authCodeService: authCodeService,
	}
}

// ShareAuthorizationCode 用户分享授权码
// @Summary 用户分享授权码
// @Description 用户可以将自己的授权码分享给其他用户
// @Tags 用户端授权码管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param codeId path string true "授权码ID"
// @Param request body models.AuthorizationCodeShareRequest true "分享请求"
// @Success 200 {object} models.APIResponse{data=models.AuthorizationCodeShareResponse} "分享成功"
// @Failure 400 {object} models.ErrorResponse "请求参数无效"
// @Failure 401 {object} models.ErrorResponse "未认证"
// @Failure 403 {object} models.ErrorResponse "授权码已被锁定"
// @Failure 404 {object} models.ErrorResponse "授权码不存在或目标用户不存在"
// @Failure 409 {object} models.ErrorResponse "不能分享给自己或分享数量超过可用激活数"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/cu/authorization-codes/{codeId}/share [post]
func (h *CuAuthorizationHandler) ShareAuthorizationCode(c *gin.Context) {
	// 从JWT Token中获取用户ID
	claims := c.MustGet("cu_user").(*utils.CuClaims)
	userID := claims.UserID

	// 获取路径参数
	codeID := c.Param("codeId")
	if codeID == "" {
		lang := middleware.GetLanguage(c)
		status, errCode, message := i18n.NewI18nErrorResponse("900001", lang)
		c.JSON(status, models.ErrorResponse{
			Code:      errCode,
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
		return
	}

	var req models.AuthorizationCodeShareRequest
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

	// 调用服务层处理分享逻辑
	data, err := h.authCodeService.ShareAuthorizationCode(c.Request.Context(), codeID, userID, &req)
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
			Message:   message,
			Timestamp: getCurrentTimestamp(),
		})
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
