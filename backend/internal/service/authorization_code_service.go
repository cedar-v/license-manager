package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"license-manager/internal/config"
	"license-manager/internal/models"
	"license-manager/internal/repository"
	pkgcontext "license-manager/pkg/context"
	"license-manager/pkg/i18n"
	"license-manager/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authorizationCodeService struct {
	authCodeRepo repository.AuthorizationCodeRepository
	customerRepo repository.CustomerRepository
	cuUserRepo   repository.CuUserRepository
	licenseRepo  repository.LicenseRepository
}

// NewAuthorizationCodeService 创建授权码服务实例
func NewAuthorizationCodeService(
	authCodeRepo repository.AuthorizationCodeRepository,
	customerRepo repository.CustomerRepository,
	cuUserRepo repository.CuUserRepository,
	licenseRepo repository.LicenseRepository,
) AuthorizationCodeService {
	return &authorizationCodeService{
		authCodeRepo: authCodeRepo,
		customerRepo: customerRepo,
		cuUserRepo:   cuUserRepo,
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

	// 授权码生成规则回退/统一为旧规则（不再自包含配置）
	authCode, err := s.generateAuthorizationCode(req.CustomerID)
	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 处理JSON字段（用于数据库存储）
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
		ID:               uuid.New().String(), // 生成新的UUID作为主键
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

// GetProductActivationCode 获取产品激活码：{授权码}&{payload}
func (s *authorizationCodeService) GetProductActivationCode(ctx context.Context, customerID string, req *models.ProductActivationCodeRequest) (*models.ProductActivationCodeResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	if customerID == "" || req == nil || strings.TrimSpace(req.AuthorizationCode) == "" {
		return nil, i18n.NewI18nError("900001", lang)
	}

	code := strings.TrimSpace(req.AuthorizationCode)
	if idx := strings.Index(code, "&"); idx > 0 {
		code = strings.TrimSpace(code[:idx])
	}
	if code == "" {
		return nil, i18n.NewI18nError("900001", lang)
	}

	authCode, err := s.licenseRepo.GetAuthorizationCodeByCode(ctx, code)
	if err != nil {
		if errors.Is(err, repository.ErrAuthorizationCodeNotFound) {
			return nil, i18n.NewI18nError("300001", lang) // 授权码不存在
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 校验归属
	if authCode.CustomerID != customerID {
		return nil, i18n.NewI18nError("100005", lang) // 权限不足
	}

	// 校验状态与有效期
	now := time.Now()
	if authCode.IsLocked {
		return nil, i18n.NewI18nError("300003", lang) // 授权码已被锁定
	}
	if now.Before(authCode.StartDate) || now.After(authCode.EndDate) {
		return nil, i18n.NewI18nError("300001", lang) // 授权码未生效或已过期
	}

	// 构造配置 JSON（客户端离线解析/自校验）
	payloadData := map[string]interface{}{
		"ver":                1,
		"authorization_code": authCode.Code,
		"start_date":         authCode.StartDate,
		"end_date":           authCode.EndDate,
		"deployment_type":    authCode.DeploymentType,
		"max_activations":    authCode.MaxActivations,
		"generated_at":       time.Now().Format(time.RFC3339),
	}

	if featureConfig := parseJSONField(authCode.FeatureConfig); len(featureConfig) > 0 {
		payloadData["feature_config"] = featureConfig
	}
	if usageLimits := parseJSONField(authCode.UsageLimits); len(usageLimits) > 0 {
		payloadData["usage_limits"] = usageLimits
	}
	if customParameters := parseJSONField(authCode.CustomParameters); len(customParameters) > 0 {
		payloadData["custom_parameters"] = customParameters
	}

	payloadDataBytes, err := json.Marshal(payloadData)
	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// RSA 私钥签名封装（与 license_service.signLicenseFile 相同结构）
	cfg := config.GetConfig()
	if cfg == nil || cfg.License.RSA.PrivateKeyPath == "" {
		return nil, i18n.NewI18nError("900004", lang, "RSA private key path not configured")
	}

	privateKey, err := utils.LoadRSAPrivateKeyFromFile(cfg.License.RSA.PrivateKeyPath)
	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	signature, err := privateKey.SignData(payloadDataBytes)
	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	signedPayload := map[string]interface{}{
		"data":      string(payloadDataBytes),
		"signature": signature,
		"algorithm": "RSA-PSS-SHA256",
	}

	signedPayloadBytes, err := json.Marshal(signedPayload)
	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	payload := base64.StdEncoding.EncodeToString(signedPayloadBytes)
	productActivationCode := fmt.Sprintf("%s&%s", authCode.Code, payload)
	return &models.ProductActivationCodeResponse{ProductActivationCode: productActivationCode}, nil
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
	return utils.GenerateLegacyAuthorizationCode(customerID)
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

// ShareAuthorizationCode 用户分享授权码
func (s *authorizationCodeService) ShareAuthorizationCode(ctx context.Context, authCodeID, userID string, req *models.AuthorizationCodeShareRequest) (*models.AuthorizationCodeShareResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if authCodeID == "" || userID == "" || req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 根据联系方式查找目标用户
	targetUser, err := s.findUserByContact(req.TargetContact)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, i18n.NewI18nError("300104", lang) // 目标用户不存在
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 验证不能分享给自己
	if targetUser.ID == userID {
		return nil, i18n.NewI18nError("300105", lang) // 不能分享给自己
	}

	// 验证目标用户状态
	if targetUser.Status != "active" {
		return nil, i18n.NewI18nError("300104", lang) // 目标用户不存在或状态异常
	}

	// 获取原授权码
	authCode, err := s.authCodeRepo.GetAuthorizationCodeByID(ctx, authCodeID)
	if err != nil {
		if errors.Is(err, repository.ErrAuthorizationCodeNotFound) {
			return nil, i18n.NewI18nError("300101", lang) // 授权码不存在
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 验证授权码所有权（只能分享自己的授权码）
	if authCode.CreatedBy != userID {
		return nil, i18n.NewI18nError("300101", lang) // 授权码不存在（无权限访问）
	}

	// 验证授权码状态
	if authCode.IsLocked {
		return nil, i18n.NewI18nError("300102", lang) // 授权码已被锁定
	}

	// 计算当前已激活次数
	currentActivations, err := s.licenseRepo.GetActiveLicenseCount(ctx, authCodeID)
	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 计算可分享次数
	availableActivations := authCode.MaxActivations - int(currentActivations)
	if req.ShareCount > availableActivations {
		return nil, i18n.NewI18nError("300103", lang) // 分享数量超过可用激活数
	}

	var newAuthCode *models.AuthorizationCode

	// 开启事务
	tx := s.authCodeRepo.BeginTransaction(ctx)
	if tx == nil {
		return nil, i18n.NewI18nError("300106", lang) // 数据库事务失败
	}
	defer func() {
		if r := recover(); r != nil {
			if gormTx, ok := tx.(*gorm.DB); ok {
				gormTx.Rollback()
			}
			panic(r)
		}
	}()

	// 在事务中执行操作
	txErr := func() error {
		// 1. 减少原授权码的max_activations
		newMaxActivations := authCode.MaxActivations - req.ShareCount
		err := s.authCodeRepo.UpdateMaxActivationsWithTx(ctx, tx, authCodeID, newMaxActivations)
		if err != nil {
			return err
		}

		// 2. 为目标用户创建新的授权码
		now := time.Now()
		code, err := s.generateAuthorizationCode(targetUser.CustomerID)
		if err != nil {
			return i18n.NewI18nError("900004", lang, err.Error())
		}
		newAuthCode = &models.AuthorizationCode{
			Code:             code,
			CustomerID:       targetUser.CustomerID, // 使用目标用户的客户ID
			CreatedBy:        targetUser.ID,         // 记录为目标用户创建的
			SoftwareID:       authCode.SoftwareID,
			Description:      authCode.Description,
			StartDate:        now,              // 从分享时刻开始
			EndDate:          authCode.EndDate, // 到原授权码结束时间
			DeploymentType:   authCode.DeploymentType,
			EncryptionType:   authCode.EncryptionType,
			SoftwareVersion:  authCode.SoftwareVersion,
			MaxActivations:   req.ShareCount,
			IsLocked:         false,
			FeatureConfig:    authCode.FeatureConfig,
			UsageLimits:      authCode.UsageLimits,
			CustomParameters: authCode.CustomParameters,
		}

		err = s.authCodeRepo.CreateAuthorizationCodeWithTx(ctx, tx, newAuthCode)
		if err != nil {
			return err
		}

		return nil
	}()

	if txErr != nil {
		if gormTx, ok := tx.(*gorm.DB); ok {
			gormTx.Rollback()
		}
		return nil, i18n.NewI18nError("300106", lang, txErr.Error()) // 数据库事务失败
	}

	// 提交事务
	if gormTx, ok := tx.(*gorm.DB); ok {
		if err := gormTx.Commit().Error; err != nil {
			return nil, i18n.NewI18nError("300106", lang, err.Error()) // 数据库事务失败
		}
	}

	// 返回新创建的授权码信息
	response := &models.AuthorizationCodeShareResponse{
		NewAuthorizationCode: models.AuthorizationCodeShareResponseItem{
			ID:             newAuthCode.ID,
			Code:           newAuthCode.Code,
			StartDate:      newAuthCode.StartDate.Format(time.RFC3339),
			EndDate:        newAuthCode.EndDate.Format(time.RFC3339),
			MaxActivations: newAuthCode.MaxActivations,
		},
	}

	return response, nil
}

// findUserByContact 根据联系方式查找用户（手机号或邮箱）
func (s *authorizationCodeService) findUserByContact(contact string) (*models.CuUser, error) {
	// 判断是手机号还是邮箱
	if strings.Contains(contact, "@") {
		// 邮箱
		return s.cuUserRepo.GetByEmail(contact)
	} else {
		// 手机号：如果没有国家代码，默认使用 +68
		phone := contact
		countryCode := "+68" // 默认国家代码

		// 检查是否已经包含国家代码
		if strings.HasPrefix(contact, "+") {
			// 如果以 + 开头，分离国家代码和手机号
			parts := strings.SplitN(contact, " ", 2)
			if len(parts) == 2 {
				countryCode = parts[0]
				phone = parts[1]
			} else {
				// 可能是直接的格式如 +8613800000000，需要分离
				// 简单处理：假设国家代码是 + 开头的1-4位数字
				for i := 1; i <= 4 && i < len(contact); i++ {
					if contact[i] >= '0' && contact[i] <= '9' {
						continue
					}
					countryCode = contact[:i]
					phone = contact[i:]
					break
				}
			}
		}

		return s.cuUserRepo.GetByPhone(phone, countryCode)
	}
}

// GetCuAuthorizationCodeList 用户端：获取当前用户授权码列表
func (s *authorizationCodeService) GetCuAuthorizationCodeList(ctx context.Context, cuUserID string, req *models.CuAuthorizationCodeListRequest) (*models.CuAuthorizationCodeListResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	if cuUserID == "" {
		return nil, i18n.NewI18nError("100004", lang)
	}

	result, err := s.authCodeRepo.GetCuAuthorizationCodeList(ctx, cuUserID, req)
	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	for i := range result.List {
		result.List[i].StatusDisplay = i18n.GetEnumMessage("authorization_code_status", result.List[i].Status, lang)
	}

	return result, nil
}

// GetCuAuthorizationCodeSummary 用户端：授权信息统计
func (s *authorizationCodeService) GetCuAuthorizationCodeSummary(ctx context.Context, cuUserID string) (*models.CuAuthorizationCodeSummaryResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	if cuUserID == "" {
		return nil, i18n.NewI18nError("100004", lang)
	}

	result, err := s.authCodeRepo.GetCuAuthorizationCodeSummary(ctx, cuUserID)
	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	return result, nil
}
