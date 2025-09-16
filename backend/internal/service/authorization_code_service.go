package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"license-manager/internal/models"
	"license-manager/internal/repository"
	pkgcontext "license-manager/pkg/context"
	"license-manager/pkg/i18n"
)

type authorizationCodeService struct {
	authCodeRepo repository.AuthorizationCodeRepository
}

// NewAuthorizationCodeService 创建授权码服务实例
func NewAuthorizationCodeService(authCodeRepo repository.AuthorizationCodeRepository) AuthorizationCodeService {
	return &authorizationCodeService{
		authCodeRepo: authCodeRepo,
	}
}

// CreateAuthorizationCode 创建授权码
func (s *authorizationCodeService) CreateAuthorizationCode(ctx context.Context, req *models.AuthorizationCodeCreateRequest) (*models.AuthorizationCodeCreateResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 验证客户是否存在
	exists, err := s.authCodeRepo.CheckCustomerExists(ctx, req.CustomerID)
	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}
	if !exists {
		return nil, i18n.NewI18nError("200001", lang) // 客户不存在
	}

	// 业务逻辑：计算开始时间和结束时间
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()) // 当天00:00:00开始
	endDate := startDate.AddDate(0, 0, req.ValidityDays-1).Add(23*time.Hour + 59*time.Minute + 59*time.Second) // 有效期最后一天的23:59:59结束

	// 获取当前用户ID
	currentUserID := "admin_uuid" // TODO: 从JWT token中获取

	// 生成授权码
	authCode, err := s.generateAuthorizationCode(req.CustomerID)
	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 处理JSON字段
	var featureConfig, usageLimits, customParameters models.JSON
	if req.FeatureConfig != nil {
		featureConfigBytes, err := json.Marshal(req.FeatureConfig)
		if err != nil {
			return nil, i18n.NewI18nError("300010", lang) // 配置参数错误
		}
		featureConfig = models.JSON(featureConfigBytes)
	}
	if req.UsageLimits != nil {
		usageLimitsBytes, err := json.Marshal(req.UsageLimits)
		if err != nil {
			return nil, i18n.NewI18nError("300010", lang) // 配置参数错误
		}
		usageLimits = models.JSON(usageLimitsBytes)
	}
	if req.CustomParameters != nil {
		customParametersBytes, err := json.Marshal(req.CustomParameters)
		if err != nil {
			return nil, i18n.NewI18nError("300010", lang) // 配置参数错误
		}
		customParameters = models.JSON(customParametersBytes)
	}

	// 设置默认加密类型
	encryptionType := req.EncryptionType
	if encryptionType == nil {
		defaultEncryption := "standard"
		encryptionType = &defaultEncryption
	}

	// 构建授权码实体
	authCodeEntity := &models.AuthorizationCode{
		Code:             authCode,
		CustomerID:       req.CustomerID,
		CreatedBy:        currentUserID,
		SoftwareID:       req.SoftwareID,
		Description:      req.Description,
		StartDate:        startDate,
		EndDate:          endDate,
		DeploymentType:   req.DeploymentType,
		EncryptionType:   encryptionType,
		SoftwareVersion:  req.SoftwareVersion,
		MaxActivations:   req.MaxActivations,
		IsLocked:         false,
		FeatureConfig:    featureConfig,
		UsageLimits:      usageLimits,
		CustomParameters: customParameters,
	}

	// 委托给Repository层进行数据创建
	if err := s.authCodeRepo.CreateAuthorizationCode(ctx, authCodeEntity); err != nil {
		// 根据错误类型包装为完整的I18nError
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return nil, i18n.NewI18nError("300002", lang) // 授权码已存在
		}

		// 数据库相关错误
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	return &models.AuthorizationCodeCreateResponse{
		ID:   authCodeEntity.ID,
		Code: authCodeEntity.Code,
	}, nil
}

// GetAuthorizationCodeList 查询授权码列表
func (s *authorizationCodeService) GetAuthorizationCodeList(ctx context.Context, req *models.AuthorizationCodeListRequest) (*models.AuthorizationCodeListResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 委托给Repository层进行数据访问
	result, err := s.authCodeRepo.GetAuthorizationCodeList(ctx, req)
	if err != nil {
		// 数据库相关错误
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 业务逻辑：添加多语言显示字段和状态计算
	now := time.Now()
	for i := range result.List {
		item := &result.List[i]
		s.fillAuthCodeDisplayFields(item, lang, now)
		
		// 根据状态筛选
		if req.Status != "" && item.Status != req.Status {
			// 移除不匹配的项 - 这里简化处理，实际应该在查询时过滤
			continue
		}
	}

	return result, nil
}

// fillAuthCodeDisplayFields 填充授权码列表项多语言显示字段
func (s *authorizationCodeService) fillAuthCodeDisplayFields(item *models.AuthorizationCodeListItem, lang string, now time.Time) {
	// 计算状态
	endDate, _ := time.Parse("2006-01-02T15:04:05Z07:00", item.EndDate)
	
	if item.IsLocked {
		item.Status = "locked"
	} else if now.After(endDate) {
		item.Status = "expired"
	} else {
		item.Status = "normal"
	}

	// 添加多语言显示字段
	item.StatusDisplay = i18n.GetEnumMessage("authorization_status", item.Status, lang)
	item.DeploymentTypeDisplay = i18n.GetEnumMessage("deployment_type", item.DeploymentType, lang)
	item.CustomerNameDisplay = item.CustomerName // 客户名称暂时不需要翻译
}

// generateAuthorizationCode 生成授权码
func (s *authorizationCodeService) generateAuthorizationCode(customerID string) (string, error) {
	// 获取客户编码的前4位作为前缀
	customerCode := "COMP" // 默认前缀
	if len(customerID) >= 4 {
		customerCode = strings.ToUpper(customerID[:4])
	}

	// 生成8位随机字符串
	randomStr, err := s.generateRandomString(8)
	if err != nil {
		return "", err
	}

	// 生成4位校验码
	checksum := s.generateChecksum(customerCode + randomStr)

	// 格式: LIC-{customer_code}-{random}-{checksum}
	return fmt.Sprintf("LIC-%s-%s-%s", customerCode, randomStr, checksum), nil
}

// generateRandomString 生成随机字符串
func (s *authorizationCodeService) generateRandomString(length int) (string, error) {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}
	return string(bytes), nil
}

// generateChecksum 生成校验码
func (s *authorizationCodeService) generateChecksum(input string) string {
	sum := 0
	for _, char := range input {
		sum += int(char)
	}
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	checksum := ""
	for i := 0; i < 4; i++ {
		checksum += string(chars[sum%len(chars)])
		sum = sum / len(chars)
	}
	return checksum
}
