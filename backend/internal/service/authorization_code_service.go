package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"license-manager/internal/models"
	"license-manager/internal/repository"
	pkgcontext "license-manager/pkg/context"
	"license-manager/pkg/i18n"
)

type authorizationCodeService struct {
	authCodeRepo repository.AuthorizationCodeRepository
	customerRepo repository.CustomerRepository
	licenseRepo  repository.LicenseRepository
}

// NewAuthorizationCodeService 创建授权码服务实例
func NewAuthorizationCodeService(
	authCodeRepo repository.AuthorizationCodeRepository,
	customerRepo repository.CustomerRepository,
	licenseRepo repository.LicenseRepository,
) AuthorizationCodeService {
	return &authorizationCodeService{
		authCodeRepo: authCodeRepo,
		customerRepo: customerRepo,
		licenseRepo:  licenseRepo,
	}
}

// CreateAuthorizationCode 创建授权码
func (s *authorizationCodeService) CreateAuthorizationCode(ctx context.Context, req *models.AuthorizationCodeCreateRequest) (*models.AuthorizationCodeCreateResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 验证客户是否存在并检查状态
	customer, err := s.customerRepo.GetCustomerByID(ctx, req.CustomerID)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerNotFound) {
			return nil, i18n.NewI18nError("200001", lang) // 客户不存在
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 检查客户状态：如果客户状态为停用（disabled），不允许创建授权
	if customer.Status == "disabled" {
		return nil, i18n.NewI18nError("200007", lang) // 客户已停用，无法创建授权
	}

	// 业务逻辑：计算开始时间和结束时间
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())                     // 当天00:00:00开始
	endDate := startDate.AddDate(0, 0, req.ValidityDays-1).Add(23*time.Hour + 59*time.Minute + 59*time.Second) // 有效期最后一天的23:59:59结束

	// 获取当前用户ID
	currentUserID := pkgcontext.GetUserIDFromContext(ctx)
	if currentUserID == "" {
		return nil, i18n.NewI18nError("100004", lang) // 缺少认证信息
	}

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

// GetAuthorizationCode 获取单个授权码详情
func (s *authorizationCodeService) GetAuthorizationCode(ctx context.Context, id string) (*models.AuthorizationCode, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if id == "" {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 委托给Repository层进行数据访问
	authCode, err := s.authCodeRepo.GetAuthorizationCodeByID(ctx, id)
	if err != nil {
		// 根据Repository错误类型包装为完整的I18nError
		if errors.Is(err, repository.ErrAuthorizationCodeNotFound) {
			return nil, i18n.NewI18nError("300001", lang)
		}

		// 数据库相关错误
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 查询客户信息
	customer, err := s.customerRepo.GetCustomerByID(ctx, authCode.CustomerID)
	if err != nil {
		// 如果客户不存在，不影响授权码查询，只是不填充客户信息
		// 记录错误日志但不返回错误
	} else {
		// 填充客户信息
		authCode.CustomerInfo = &models.CustomerInfoForAuthCode{
			ID:                  customer.ID,
			CustomerCode:        customer.CustomerCode,
			CustomerName:        customer.CustomerName,
			CustomerType:        customer.CustomerType,
			CustomerTypeDisplay: i18n.GetEnumMessage("customer_type", customer.CustomerType, lang),
			Status:              customer.Status,
			StatusDisplay:       i18n.GetEnumMessage("customer_status", customer.Status, lang),
			CreatedAt:           customer.CreatedAt.Format(time.RFC3339),
		}
	}

	// 查询已激活的许可证数量
	activatedCount, err := s.licenseRepo.GetActiveLicenseCount(ctx, id)
	if err != nil {
		// 查询失败时，数量设为0，记录错误日志但不返回错误
		activatedCount = 0
	}
	authCode.ActivatedLicensesCount = activatedCount

	// 填充多语言显示字段和计算状态
	s.fillAuthorizationCodeDisplayFields(authCode, lang)

	return authCode, nil
}

// GenerateAuthorizationFile 生成授权码文件内容
func (s *authorizationCodeService) GenerateAuthorizationFile(ctx context.Context, id string) ([]byte, string, string, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	if id == "" {
		return nil, "", "", i18n.NewI18nError("900001", lang)
	}

	authCode, err := s.authCodeRepo.GetAuthorizationCodeByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAuthorizationCodeNotFound) {
			return nil, "", "", i18n.NewI18nError("300001", lang) // 授权码不存在
		}
		return nil, "", "", i18n.NewI18nError("900004", lang, err.Error())
	}

	now := time.Now()
	if authCode.IsLocked {
		return nil, "", "", i18n.NewI18nError("300003", lang) // 授权码已锁定
	}
	if now.Before(authCode.StartDate) || now.After(authCode.EndDate) {
		return nil, "", "", i18n.NewI18nError("300001", lang) // 授权码未生效或已过期
	}

	fileName := fmt.Sprintf("authorization_%s.txt", authCode.Code)
	return []byte(authCode.Code), fileName, authCode.Code, nil
}

// fillAuthorizationCodeDisplayFields 填充完整授权码模型的多语言显示字段
func (s *authorizationCodeService) fillAuthorizationCodeDisplayFields(authCode *models.AuthorizationCode, lang string) {
	// 计算状态
	now := time.Now()
	if authCode.IsLocked {
		authCode.Status = "locked"
	} else if now.After(authCode.EndDate) {
		authCode.Status = "expired"
	} else {
		authCode.Status = "normal"
	}

	// 添加多语言显示字段
	authCode.StatusDisplay = i18n.GetEnumMessage("authorization_status", authCode.Status, lang)
	authCode.DeploymentTypeDisplay = i18n.GetEnumMessage("deployment_type", authCode.DeploymentType, lang)
	if authCode.EncryptionType != nil {
		authCode.EncryptionTypeDisplay = i18n.GetEnumMessage("encryption_type", *authCode.EncryptionType, lang)
	}

	// TODO: 统计当前激活数量
	authCode.CurrentActivations = 0

	// TODO: 获取客户名称
	// 可以考虑在Repository层通过JOIN获取，或在这里单独查询
}

// UpdateAuthorizationCode 更新授权码
func (s *authorizationCodeService) UpdateAuthorizationCode(ctx context.Context, id string, req *models.AuthorizationCodeUpdateRequest) (*models.AuthorizationCode, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if id == "" {
		return nil, i18n.NewI18nError("900001", lang)
	}
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 先查询现有授权码
	existingAuthCode, err := s.authCodeRepo.GetAuthorizationCodeByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAuthorizationCodeNotFound) {
			return nil, i18n.NewI18nError("300001", lang)
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 获取当前用户ID
	currentUserID := pkgcontext.GetUserIDFromContext(ctx)
	if currentUserID == "" {
		return nil, i18n.NewI18nError("100004", lang)
	}

	// 记录变更前的配置
	oldConfig := s.buildConfigSnapshot(existingAuthCode)

	// 只更新提供的字段
	if req.SoftwareID != nil {
		existingAuthCode.SoftwareID = req.SoftwareID
	}
	if req.Description != nil {
		existingAuthCode.Description = req.Description
	}
	// 支持通过起止时间直接设置（优先），格式 YYYY-MM-DD；若未提供，则继续支持 validity_days（向后兼容）
	if req.StartDate != nil || req.EndDate != nil {
		if req.StartDate == nil || req.EndDate == nil {
			return nil, i18n.NewI18nError("900001", lang)
		}
		startParsed, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, i18n.NewI18nError("900001", lang)
		}
		endParsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, i18n.NewI18nError("900001", lang)
		}
		// 规范化到本地时区的 00:00:00 / 23:59:59
		startDate := time.Date(startParsed.Year(), startParsed.Month(), startParsed.Day(), 0, 0, 0, 0, time.Local)
		endDate := time.Date(endParsed.Year(), endParsed.Month(), endParsed.Day(), 23, 59, 59, 0, time.Local)
		existingAuthCode.StartDate = startDate
		existingAuthCode.EndDate = endDate
	} else if req.ValidityDays != nil {
		// 重新计算开始时间和结束时间（兼容旧的 validity_days 字段）
		now := time.Now()
		startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate := startDate.AddDate(0, 0, *req.ValidityDays-1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		existingAuthCode.StartDate = startDate
		existingAuthCode.EndDate = endDate
	}
	if req.DeploymentType != nil {
		existingAuthCode.DeploymentType = *req.DeploymentType
	}
	if req.EncryptionType != nil {
		existingAuthCode.EncryptionType = req.EncryptionType
	}
	if req.SoftwareVersion != nil {
		existingAuthCode.SoftwareVersion = req.SoftwareVersion
	}
	if req.MaxActivations != nil {
		existingAuthCode.MaxActivations = *req.MaxActivations
	}

	// 处理JSON字段
	if req.FeatureConfig != nil {
		featureConfigBytes, err := json.Marshal(req.FeatureConfig)
		if err != nil {
			return nil, i18n.NewI18nError("300010", lang)
		}
		existingAuthCode.FeatureConfig = models.JSON(featureConfigBytes)
	}
	if req.UsageLimits != nil {
		usageLimitsBytes, err := json.Marshal(req.UsageLimits)
		if err != nil {
			return nil, i18n.NewI18nError("300010", lang)
		}
		existingAuthCode.UsageLimits = models.JSON(usageLimitsBytes)
	}
	if req.CustomParameters != nil {
		customParametersBytes, err := json.Marshal(req.CustomParameters)
		if err != nil {
			return nil, i18n.NewI18nError("300010", lang)
		}
		existingAuthCode.CustomParameters = models.JSON(customParametersBytes)
	}

	// 委托给Repository层进行数据更新
	if err := s.authCodeRepo.UpdateAuthorizationCode(ctx, existingAuthCode); err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 填充多语言显示字段和计算状态
	s.fillAuthorizationCodeDisplayFields(existingAuthCode, lang)

	// 记录变更历史到 authorization_changes 表
	newConfig := s.buildConfigSnapshot(existingAuthCode)
	if err := s.recordAuthorizationChange(ctx, id, req.ChangeType, req.Reason, currentUserID, oldConfig, newConfig); err != nil {
		log.Printf("记录授权变更历史失败: %v", err)
	}

	return existingAuthCode, nil
}

// LockUnlockAuthorizationCode 锁定/解锁授权码
func (s *authorizationCodeService) LockUnlockAuthorizationCode(ctx context.Context, id string, req *models.AuthorizationCodeLockRequest) (*models.AuthorizationCode, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if id == "" {
		return nil, i18n.NewI18nError("900001", lang)
	}
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 先查询现有授权码
	existingAuthCode, err := s.authCodeRepo.GetAuthorizationCodeByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAuthorizationCodeNotFound) {
			return nil, i18n.NewI18nError("300001", lang)
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 获取当前用户ID
	currentUserID := pkgcontext.GetUserIDFromContext(ctx)
	if currentUserID == "" {
		return nil, i18n.NewI18nError("100004", lang)
	}

	// 记录变更前的配置
	oldConfig := s.buildConfigSnapshot(existingAuthCode)

	// 更新锁定状态
	now := time.Now()
	existingAuthCode.IsLocked = req.IsLocked

	if req.IsLocked {
		// 锁定
		existingAuthCode.LockReason = req.LockReason
		existingAuthCode.LockedAt = &now
		existingAuthCode.LockedBy = &currentUserID
	} else {
		// 解锁
		existingAuthCode.LockReason = nil
		existingAuthCode.LockedAt = nil
		existingAuthCode.LockedBy = nil
	}

	// 委托给Repository层进行数据更新
	if err := s.authCodeRepo.UpdateAuthorizationCode(ctx, existingAuthCode); err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 填充多语言显示字段和计算状态
	s.fillAuthorizationCodeDisplayFields(existingAuthCode, lang)

	// 记录变更历史到 authorization_changes 表
	changeType := "lock"
	if !req.IsLocked {
		changeType = "unlock"
	}
	newConfig := s.buildConfigSnapshot(existingAuthCode)
	if err := s.recordAuthorizationChange(ctx, id, changeType, req.Reason, currentUserID, oldConfig, newConfig); err != nil {
		log.Printf("记录授权变更历史失败: %v", err)
	}

	return existingAuthCode, nil
}

// DeleteAuthorizationCode 删除授权码（软删除）
func (s *authorizationCodeService) DeleteAuthorizationCode(ctx context.Context, id string) error {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if id == "" {
		return i18n.NewI18nError("900001", lang)
	}

	// 获取当前用户ID
	currentUserID := pkgcontext.GetUserIDFromContext(ctx)
	if currentUserID == "" {
		return i18n.NewI18nError("100004", lang)
	}

	// 先查询现有授权码，确保存在
	existingAuthCode, err := s.authCodeRepo.GetAuthorizationCodeByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAuthorizationCodeNotFound) {
			return i18n.NewI18nError("300001", lang)
		}
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	// 记录变更前的配置
	oldConfig := s.buildConfigSnapshot(existingAuthCode)

	// 委托给Repository层进行软删除
	if err := s.authCodeRepo.DeleteAuthorizationCode(ctx, id); err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	// 记录变更历史到 authorization_changes 表
	emptyConfig := make(map[string]interface{})
	if err := s.recordAuthorizationChange(ctx, id, "delete", nil, currentUserID, oldConfig, emptyConfig); err != nil {
		log.Printf("记录授权变更历史失败: %v", err)
	}

	return nil
}

// recordAuthorizationChange 记录授权变更历史
func (s *authorizationCodeService) recordAuthorizationChange(ctx context.Context, authCodeID string, changeType string, reason *string, operatorID string, oldConfig, newConfig map[string]interface{}) error {
	// 构建变更历史记录
	change := &models.AuthorizationChange{
		AuthorizationCodeID: authCodeID,
		ChangeType:          changeType,
		OperatorID:          operatorID,
		Reason:              reason,
	}

	// 序列化配置为JSON
	if oldConfig != nil {
		oldConfigBytes, err := json.Marshal(oldConfig)
		if err != nil {
			return err
		}
		change.OldConfig = models.JSON(oldConfigBytes)
	}

	if newConfig != nil {
		newConfigBytes, err := json.Marshal(newConfig)
		if err != nil {
			return err
		}
		change.NewConfig = models.JSON(newConfigBytes)
	}

	// 委托给Repository层记录变更历史
	return s.authCodeRepo.RecordAuthorizationChange(ctx, change)
}

// buildConfigSnapshot 构建配置快照，用于记录变更历史
func (s *authorizationCodeService) buildConfigSnapshot(authCode *models.AuthorizationCode) map[string]interface{} {
	config := make(map[string]interface{})

	// 基础配置
	config["code"] = authCode.Code
	config["software_id"] = authCode.SoftwareID
	config["description"] = authCode.Description
	config["start_date"] = authCode.StartDate.Format(time.RFC3339)
	config["end_date"] = authCode.EndDate.Format(time.RFC3339)
	config["deployment_type"] = authCode.DeploymentType
	config["encryption_type"] = authCode.EncryptionType
	config["software_version"] = authCode.SoftwareVersion
	config["max_activations"] = authCode.MaxActivations
	config["is_locked"] = authCode.IsLocked
	config["lock_reason"] = authCode.LockReason

	// JSON配置字段
	if len(authCode.FeatureConfig) > 0 {
		var featureConfig map[string]interface{}
		if err := json.Unmarshal(authCode.FeatureConfig, &featureConfig); err == nil {
			config["feature_config"] = featureConfig
		}
	}

	if len(authCode.UsageLimits) > 0 {
		var usageLimits map[string]interface{}
		if err := json.Unmarshal(authCode.UsageLimits, &usageLimits); err == nil {
			config["usage_limits"] = usageLimits
		}
	}

	if len(authCode.CustomParameters) > 0 {
		var customParameters map[string]interface{}
		if err := json.Unmarshal(authCode.CustomParameters, &customParameters); err == nil {
			config["custom_parameters"] = customParameters
		}
	}

	return config
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
// 格式: LIC-{4位客户代码}-{12位随机}-{4位校验码}
// 安全性：12位随机字符（62字符集）提供约 3.23 × 10^21 种可能组合
func (s *authorizationCodeService) generateAuthorizationCode(customerID string) (string, error) {
	// 获取客户编码的前4位作为前缀
	customerCode := "COMP" // 默认前缀
	if len(customerID) >= 4 {
		customerCode = strings.ToUpper(customerID[:4])
	}

	// 生成12位随机字符串（从8位增加到12位以提升安全性）
	randomStr, err := s.generateRandomString(12)
	if err != nil {
		return "", err
	}

	// 生成4位校验码
	checksum := s.generateChecksum(customerCode + randomStr)

	// 格式: LIC-{customer_code}-{random}-{checksum}
	return fmt.Sprintf("LIC-%s-%s-%s", customerCode, randomStr, checksum), nil
}

// generateRandomString 生成随机字符串
// 使用62字符集（A-Z, a-z, 0-9）提供更高的熵值
func (s *authorizationCodeService) generateRandomString(length int) (string, error) {
	// 扩展字符集到62字符（A-Z, a-z, 0-9）以提升安全性
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
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
// 基于输入字符串的ASCII值计算，用于格式验证
func (s *authorizationCodeService) generateChecksum(input string) string {
	sum := 0
	for _, char := range input {
		sum += int(char)
	}
	// 使用与随机字符串相同的62字符集以保持一致性
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	checksum := ""
	for i := 0; i < 4; i++ {
		checksum += string(chars[sum%len(chars)])
		sum = sum / len(chars)
	}
	return checksum
}

// ValidateAuthorizationCodeFormat 验证授权码格式和校验和
// 格式: LIC-{4位客户代码}-{8或12位随机}-{4位校验码}
// 兼容旧格式（8位随机，36字符集）和新格式（12位随机，62字符集）
func (s *authorizationCodeService) ValidateAuthorizationCodeFormat(code string) bool {
	// 解析格式
	parts := strings.Split(code, "-")
	if len(parts) != 4 || parts[0] != "LIC" {
		return false
	}

	// 验证各部分长度（兼容旧格式8位和新格式12位）
	customerCodeLen := len(parts[1])
	randomLen := len(parts[2])
	checksumLen := len(parts[3])

	// 客户代码必须为4位，校验和必须为4位
	if customerCodeLen != 4 || checksumLen != 4 {
		return false
	}

	// 随机部分：兼容8位（旧格式）和12位（新格式）
	if randomLen != 8 && randomLen != 12 {
		return false
	}

	// 根据随机部分长度选择字符集验证校验和
	// 8位：使用36字符集（旧格式）
	// 12位：使用62字符集（新格式）
	expectedChecksum := s.generateChecksumWithCharset(parts[1]+parts[2], randomLen == 12)
	return expectedChecksum == parts[3]
}

// generateChecksumWithCharset 生成校验码（支持指定字符集）
// useExtendedCharset: true使用62字符集（新格式），false使用36字符集（旧格式）
func (s *authorizationCodeService) generateChecksumWithCharset(input string, useExtendedCharset bool) string {
	sum := 0
	for _, char := range input {
		sum += int(char)
	}

	var chars string
	if useExtendedCharset {
		// 使用62字符集（A-Z, a-z, 0-9）用于新格式
		chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	} else {
		// 使用36字符集（A-Z, 0-9）用于旧格式
		chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	checksum := ""
	for i := 0; i < 4; i++ {
		checksum += string(chars[sum%len(chars)])
		sum = sum / len(chars)
	}
	return checksum
}

// GetAuthorizationChangeList 查询授权变更历史列表
func (s *authorizationCodeService) GetAuthorizationChangeList(ctx context.Context, authCodeID string, req *models.AuthorizationChangeListRequest) (*models.AuthorizationChangeListResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if authCodeID == "" {
		return nil, i18n.NewI18nError("900001", lang)
	}
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 验证授权码是否存在
	_, err := s.authCodeRepo.GetAuthorizationCodeByID(ctx, authCodeID)
	if err != nil {
		if errors.Is(err, repository.ErrAuthorizationCodeNotFound) {
			return nil, i18n.NewI18nError("300001", lang) // 授权码不存在
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 委托给Repository层进行数据访问
	result, err := s.authCodeRepo.GetAuthorizationChangeList(ctx, authCodeID, req)
	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 业务逻辑：添加多语言显示字段
	for i := range result.List {
		s.fillChangeDisplayFields(&result.List[i], lang)
	}

	return result, nil
}

// fillChangeDisplayFields 填充变更历史列表项多语言显示字段
func (s *authorizationCodeService) fillChangeDisplayFields(item *models.AuthorizationChangeListItem, lang string) {
	// 填充变更类型显示字段
	item.ChangeTypeDisplay = i18n.GetEnumMessage("authorization_change_type", item.ChangeType, lang)
}
