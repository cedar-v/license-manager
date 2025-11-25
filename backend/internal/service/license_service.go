package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"license-manager/internal/config"
	"license-manager/internal/models"
	"license-manager/internal/repository"
	pkgcontext "license-manager/pkg/context"
	"license-manager/pkg/i18n"
	"license-manager/pkg/utils"
)

type licenseService struct {
	licenseRepo   repository.LicenseRepository
	db            *gorm.DB
	logger        *logrus.Logger
	rsaPrivateKey *utils.RSAPrivateKey // RSA私钥（缓存）
}

// NewLicenseService 创建许可证服务实例
func NewLicenseService(licenseRepo repository.LicenseRepository, db *gorm.DB, logger *logrus.Logger) LicenseService {
	return &licenseService{
		licenseRepo: licenseRepo,
		db:          db,
		logger:      logger,
	}
}

// GetLicenseList 查询许可证列表
func (s *licenseService) GetLicenseList(ctx context.Context, req *models.LicenseListRequest) (*models.LicenseListResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang) // 业务错误，不覆盖多语言message
	}

	// 委托给Repository层进行数据访问
	result, err := s.licenseRepo.GetLicenseList(ctx, req)
	if err != nil {
		// 数据库相关错误，使用系统错误码，覆盖message显示详细信息
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 业务逻辑：添加多语言显示字段
	for i := range result.List {
		s.fillDisplayFields(&result.List[i], lang)
	}

	return result, nil
}

// fillDisplayFields 填充列表项多语言显示字段
func (s *licenseService) fillDisplayFields(item *models.LicenseListItem, lang string) {
	// 填充状态显示字段
	item.StatusDisplay = i18n.GetEnumMessage("license_status", item.Status, lang)

	// 填充在线状态显示字段
	if item.IsOnline {
		item.IsOnlineDisplay = i18n.GetEnumMessage("license_online_status", "online", lang)
	} else {
		item.IsOnlineDisplay = i18n.GetEnumMessage("license_online_status", "offline", lang)
	}
}

// GetLicense 获取单个许可证详情
func (s *licenseService) GetLicense(ctx context.Context, id string) (*models.LicenseDetailResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if id == "" {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 委托给Repository层进行数据访问
	license, err := s.licenseRepo.GetLicenseByID(ctx, id)
	if err != nil {
		// 根据Repository错误类型包装为完整的I18nError
		if errors.Is(err, repository.ErrLicenseNotFound) {
			return nil, i18n.NewI18nError("300006", lang) // 许可证不存在
		}

		// 数据库相关错误，使用系统错误码，覆盖message显示详细信息
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 转换为详情响应格式
	response := s.convertToDetailResponse(license, lang)

	return response, nil
}

// CreateLicense 创建许可证
func (s *licenseService) CreateLicense(ctx context.Context, req *models.LicenseCreateRequest) (*models.License, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 获取授权码信息以获取客户ID
	authCode, err := s.licenseRepo.GetAuthorizationCodeByID(ctx, req.AuthorizationCodeID)
	if err != nil {
		if errors.Is(err, repository.ErrAuthorizationCodeNotFound) {
			return nil, i18n.NewI18nError("300001", lang) // 授权码不存在
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 生成许可证密钥
	licenseKey, err := s.generateLicenseKey()
	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	now := time.Now()

	// 构建许可证实体，手动添加的离线许可证直接视为已激活
	license := &models.License{
		LicenseKey:          licenseKey,
		AuthorizationCodeID: req.AuthorizationCodeID,
		CustomerID:          authCode.CustomerID, // 从授权码获取客户ID
		HardwareFingerprint: req.HardwareFingerprint,
		ActivationIP:        req.ActivationIP,
		Status:              "active",
		ActivatedAt:         &now,
	}

	// 设置设备信息
	if req.DeviceInfo != nil {
		deviceInfoBytes, err := json.Marshal(req.DeviceInfo)
		if err == nil {
			license.DeviceInfo = models.JSON(deviceInfoBytes)
		}
	}

	// 委托给Repository层进行数据创建
	if err := s.licenseRepo.CreateLicense(ctx, license); err != nil {
		// 根据错误类型包装为完整的I18nError
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return nil, i18n.NewI18nError("300007", lang) // 许可证已存在
		}

		// 数据库相关错误
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	return license, nil
}

// convertToDetailResponse 转换为详情响应格式
func (s *licenseService) convertToDetailResponse(license *models.License, lang string) *models.LicenseDetailResponse {
	response := &models.LicenseDetailResponse{
		ID:                  license.ID,
		LicenseKey:          license.LicenseKey,
		AuthorizationCodeID: license.AuthorizationCodeID,
		CustomerID:          license.CustomerID,
		HardwareFingerprint: license.HardwareFingerprint,
		ActivationIP:        license.ActivationIP,
		Status:              license.Status,
		IsOnline:            license.IsOnline,
		LastOnlineIP:        license.LastOnlineIP,
		CreatedAt:           license.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           license.UpdatedAt.Format(time.RFC3339),
	}

	// 格式化时间字段
	if license.ActivatedAt != nil {
		activatedAt := license.ActivatedAt.Format(time.RFC3339)
		response.ActivatedAt = &activatedAt
	}
	if license.LastHeartbeat != nil {
		lastHeartbeat := license.LastHeartbeat.Format(time.RFC3339)
		response.LastHeartbeat = &lastHeartbeat
	}
	if license.ConfigUpdatedAt != nil {
		configUpdatedAt := license.ConfigUpdatedAt.Format(time.RFC3339)
		response.ConfigUpdatedAt = &configUpdatedAt
	}

	// 设置关联数据
	if license.AuthorizationCode != nil {
		response.AuthorizationCode = license.AuthorizationCode.Code
	}
	if license.Customer != nil {
		response.CustomerName = license.Customer.CustomerName
	}

	// 解析JSON字段
	if len(license.DeviceInfo) > 0 {
		var deviceInfo map[string]interface{}
		if err := json.Unmarshal(license.DeviceInfo, &deviceInfo); err == nil {
			response.DeviceInfo = deviceInfo
		}
	}
	if len(license.UsageData) > 0 {
		var usageData map[string]interface{}
		if err := json.Unmarshal(license.UsageData, &usageData); err == nil {
			response.UsageData = usageData
		}
	}

	// 填充多语言显示字段
	response.StatusDisplay = i18n.GetEnumMessage("license_status", license.Status, lang)
	if license.IsOnline {
		response.IsOnlineDisplay = i18n.GetEnumMessage("license_online_status", "online", lang)
	} else {
		response.IsOnlineDisplay = i18n.GetEnumMessage("license_online_status", "offline", lang)
	}

	return response
}

// RevokeLicense 撤销许可证
func (s *licenseService) RevokeLicense(ctx context.Context, id string, req *models.LicenseRevokeRequest) (*models.License, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 业务逻辑：参数验证
	if id == "" {
		return nil, i18n.NewI18nError("900001", lang)
	}
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 先查询现有许可证
	existingLicense, err := s.licenseRepo.GetLicenseByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrLicenseNotFound) {
			return nil, i18n.NewI18nError("300006", lang) // 许可证不存在
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 检查许可证是否已经被撤销
	if existingLicense.Status == "revoked" {
		return nil, i18n.NewI18nError("300007", lang) // 许可证已被撤销
	}

	// 更新许可证状态为撤销
	existingLicense.Status = "revoked"

	// 如果提供了撤销原因，保存到使用数据中
	if req.Reason != "" {
		revokeData := map[string]interface{}{
			"revoked_at":    time.Now().Format(time.RFC3339),
			"revoke_reason": req.Reason,
		}

		// 保留原有的使用数据
		if len(existingLicense.UsageData) > 0 {
			var existingUsageData map[string]interface{}
			if err := json.Unmarshal(existingLicense.UsageData, &existingUsageData); err == nil {
				for k, v := range existingUsageData {
					revokeData[k] = v
				}
			}
		}

		usageDataBytes, err := json.Marshal(revokeData)
		if err == nil {
			existingLicense.UsageData = models.JSON(usageDataBytes)
		}
	}

	// 委托给Repository层进行数据更新
	if err := s.licenseRepo.UpdateLicense(ctx, existingLicense); err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	return existingLicense, nil
}

// GenerateLicenseFile 生成许可证文件
func (s *licenseService) GenerateLicenseFile(ctx context.Context, id string) ([]byte, string, string, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)
	s.logger.Infof("[GenerateLicenseFile] 开始生成许可证文件，license_id: %s", id)

	// 业务逻辑：参数验证
	if id == "" {
		s.logger.Error("[GenerateLicenseFile] 错误：许可证ID为空")
		return nil, "", "", i18n.NewI18nError("900001", lang)
	}

	// 获取许可证信息
	s.logger.Infof("[GenerateLicenseFile] 查询许可证信息，license_id: %s", id)
	license, err := s.licenseRepo.GetLicenseByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[GenerateLicenseFile] 查询许可证失败，license_id: %s, error: %v", id, err)
		if errors.Is(err, repository.ErrLicenseNotFound) {
			return nil, "", "", i18n.NewI18nError("300006", lang) // 许可证不存在
		}
		return nil, "", "", i18n.NewI18nError("900004", lang, err.Error())
	}

	s.logger.Infof("[GenerateLicenseFile] 许可证信息查询成功，license_key: %s, status: %s", license.LicenseKey, license.Status)

	// 检查许可证状态
	if license.Status == "revoked" {
		s.logger.Warnf("[GenerateLicenseFile] 许可证已被撤销，license_id: %s", id)
		return nil, "", "", i18n.NewI18nError("300007", lang) // 许可证已被撤销
	}

	// 构建许可证文件内容
	licenseFileData := map[string]interface{}{
		"license_key":           license.LicenseKey,
		"authorization_code_id": license.AuthorizationCodeID,
		"hardware_fingerprint":  license.HardwareFingerprint,
		"status":                license.Status,
		"activated_at":          license.ActivatedAt,
		"config_updated_at":     license.ConfigUpdatedAt,
		"generated_at":          time.Now().Format(time.RFC3339),
	}

	// 包含授权码信息
	if license.AuthorizationCode != nil {
		licenseFileData["authorization_code"] = license.AuthorizationCode.Code
		licenseFileData["start_date"] = license.AuthorizationCode.StartDate
		licenseFileData["end_date"] = license.AuthorizationCode.EndDate
		licenseFileData["deployment_type"] = license.AuthorizationCode.DeploymentType
		licenseFileData["max_activations"] = license.AuthorizationCode.MaxActivations

		// 包含功能配置
		if len(license.AuthorizationCode.FeatureConfig) > 0 {
			var featureConfig map[string]interface{}
			if err := json.Unmarshal(license.AuthorizationCode.FeatureConfig, &featureConfig); err == nil {
				licenseFileData["feature_config"] = featureConfig
			}
		}

		// 包含使用限制
		if len(license.AuthorizationCode.UsageLimits) > 0 {
			var usageLimits map[string]interface{}
			if err := json.Unmarshal(license.AuthorizationCode.UsageLimits, &usageLimits); err == nil {
				licenseFileData["usage_limits"] = usageLimits
			}
		}

		// 包含自定义参数
		if len(license.AuthorizationCode.CustomParameters) > 0 {
			var customParameters map[string]interface{}
			if err := json.Unmarshal(license.AuthorizationCode.CustomParameters, &customParameters); err == nil {
				licenseFileData["custom_parameters"] = customParameters
			}
		}
	}

	// 序列化许可证数据
	licenseJSON, err := json.Marshal(licenseFileData)
	if err != nil {
		return nil, "", "", i18n.NewI18nError("300009", lang) // 许可证文件生成失败
	}

	// 使用RSA数字签名
	encryptedData, err := s.signLicenseFile(licenseJSON)
	if err != nil {
		return nil, "", "", i18n.NewI18nError("300009", lang) // 许可证文件生成失败
	}

	// 生成文件名
	fileName := fmt.Sprintf("license_%s.lic", license.LicenseKey)

	return encryptedData, fileName, license.LicenseKey, nil
}

// generateLicenseKey 生成许可证密钥
func (s *licenseService) generateLicenseKey() (string, error) {
	// 生成12字节随机数据
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// 转换为16进制字符串并格式化
	hex := strings.ToUpper(hex.EncodeToString(bytes))
	return fmt.Sprintf("LIC-DEVICE-%s", hex), nil
}

// validateAuthorizationCodeFormat 验证授权码格式和校验和
// 格式: LIC-{4位客户代码}-{8或12位随机}-{4位校验码}
// 兼容旧格式（8位随机，36字符集）和新格式（12位随机，62字符集）
func validateAuthorizationCodeFormat(code string) bool {
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
	expectedChecksum := generateChecksumForCode(parts[1]+parts[2], randomLen == 12)
	return expectedChecksum == parts[3]
}

// generateChecksumForCode 生成校验码（用于验证）
// 基于输入字符串的ASCII值计算
// useExtendedCharset: true使用62字符集（新格式），false使用36字符集（旧格式）
func generateChecksumForCode(input string, useExtendedCharset bool) string {
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

// ActivateLicense 激活许可证
func (s *licenseService) ActivateLicense(ctx context.Context, req *models.ActivateRequest, clientIP string) (*models.ActivateResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 参数验证
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 验证授权码格式（在查询数据库前先验证格式，避免无效查询）
	if !validateAuthorizationCodeFormat(req.AuthorizationCode) {
		return nil, i18n.NewI18nError("300001", lang) // 授权码格式无效
	}

	// 获取授权码信息
	authCode, err := s.licenseRepo.GetAuthorizationCodeByCode(ctx, req.AuthorizationCode)
	if err != nil {
		if errors.Is(err, repository.ErrAuthorizationCodeNotFound) {
			return nil, i18n.NewI18nError("300001", lang) // 授权码不存在
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 检查授权码状态和有效期
	now := time.Now()
	if authCode.IsLocked {
		return nil, i18n.NewI18nError("300003", lang) // 授权码已被锁定
	}
	if now.Before(authCode.StartDate) || now.After(authCode.EndDate) {
		return nil, i18n.NewI18nError("300001", lang) // 授权码已过期
	}

	// 使用事务确保并发安全
	var response *models.ActivateResponse
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 检查当前激活数量
		count, err := s.licenseRepo.GetActiveLicenseCount(ctx, authCode.ID)
		if err != nil {
			return err
		}

		// 检查是否已存在相同硬件指纹的许可证
		var existingLicense models.License
		err = tx.Where("authorization_code_id = ? AND hardware_fingerprint = ?",
			authCode.ID, req.HardwareFingerprint).First(&existingLicense).Error

		if err == nil {
			// 已存在，直接激活
			existingLicense.Status = "active"
			existingLicense.ActivationIP = &clientIP
			now := time.Now()
			existingLicense.ActivatedAt = &now
			existingLicense.LastHeartbeat = &now
			existingLicense.LastOnlineIP = &clientIP

			if err := tx.Save(&existingLicense).Error; err != nil {
				return err
			}

			// 生成许可证文件
			licenseFile, err := s.generateLicenseFileContent(&existingLicense, authCode)
			if err != nil {
				return err
			}

			response = &models.ActivateResponse{
				LicenseKey:        existingLicense.LicenseKey,
				LicenseFile:       licenseFile,
				HeartbeatInterval: 300, // 5分钟
			}
			return nil
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		// 检查激活数量限制
		if count >= int64(authCode.MaxActivations) {
			return repository.ErrLicenseNotFound // 使用已有错误，表示激活数量已达上限
		}

		// 生成新的许可证
		licenseKey, err := s.generateLicenseKey()
		if err != nil {
			return err
		}

		license := &models.License{
			LicenseKey:          licenseKey,
			AuthorizationCodeID: authCode.ID,
			CustomerID:          authCode.CustomerID,
			HardwareFingerprint: req.HardwareFingerprint,
			ActivationIP:        &clientIP,
			Status:              "active",
			ActivatedAt:         &now,
			LastHeartbeat:       &now,
			LastOnlineIP:        &clientIP,
		}

		// 设置设备信息
		if req.DeviceInfo != nil {
			deviceInfoBytes, err := json.Marshal(req.DeviceInfo)
			if err == nil {
				license.DeviceInfo = models.JSON(deviceInfoBytes)
			}
		}

		if err := tx.Create(license).Error; err != nil {
			return err
		}

		// 生成许可证文件
		licenseFile, err := s.generateLicenseFileContent(license, authCode)
		if err != nil {
			return err
		}

		response = &models.ActivateResponse{
			LicenseKey:        license.LicenseKey,
			LicenseFile:       licenseFile,
			HeartbeatInterval: 300,
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, repository.ErrLicenseNotFound) {
			return nil, i18n.NewI18nError("300004", lang) // 激活数量已达上限
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	return response, nil
}

// Heartbeat 心跳检测
func (s *licenseService) Heartbeat(ctx context.Context, req *models.HeartbeatRequest, clientIP string) (*models.HeartbeatResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 参数验证
	if req == nil {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 获取许可证信息
	license, err := s.licenseRepo.GetLicenseByKey(ctx, req.LicenseKey)
	if err != nil {
		if errors.Is(err, repository.ErrLicenseNotFound) {
			return nil, i18n.NewI18nError("300006", lang) // 许可证不存在
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 检查许可证状态
	if license.Status == "revoked" {
		return nil, i18n.NewI18nError("300007", lang) // 许可证已被撤销
	}

	// 更新心跳时间和使用数据
	now := time.Now()
	license.LastHeartbeat = &now
	license.LastOnlineIP = &clientIP

	// 更新使用数据
	if req.UsageData != nil {
		usageDataBytes, err := json.Marshal(req.UsageData)
		if err == nil {
			license.UsageData = models.JSON(usageDataBytes)
		}
	}

	// 检查配置是否有更新
	configUpdated := false
	var licenseFile *string

	if req.ConfigUpdatedAt != nil && license.AuthorizationCode != nil {
		clientConfigTime, err := time.Parse(time.RFC3339, *req.ConfigUpdatedAt)
		if err == nil && license.AuthorizationCode.UpdatedAt.After(clientConfigTime) {
			configUpdated = true
			// 生成新的许可证文件
			fileContent, err := s.generateLicenseFileContent(license, license.AuthorizationCode)
			if err == nil {
				licenseFile = &fileContent
				license.ConfigUpdatedAt = &now
			}
		}
	}

	// 保存更新
	if err := s.licenseRepo.UpdateLicense(ctx, license); err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	response := &models.HeartbeatResponse{
		Status:            license.Status,
		ConfigUpdated:     configUpdated,
		LicenseFile:       licenseFile,
		HeartbeatInterval: 300,
	}

	return response, nil
}

// generateLicenseFileContent 生成许可证文件内容
func (s *licenseService) generateLicenseFileContent(license *models.License, authCode *models.AuthorizationCode) (string, error) {
	// 构建许可证文件内容
	licenseFileData := map[string]interface{}{
		"license_key":           license.LicenseKey,
		"authorization_code_id": license.AuthorizationCodeID,
		"hardware_fingerprint":  license.HardwareFingerprint,
		"status":                license.Status,
		"activated_at":          license.ActivatedAt,
		"config_updated_at":     license.ConfigUpdatedAt,
		"generated_at":          time.Now().Format(time.RFC3339),
	}

	if authCode != nil {
		licenseFileData["authorization_code"] = authCode.Code
		licenseFileData["start_date"] = authCode.StartDate
		licenseFileData["end_date"] = authCode.EndDate
		licenseFileData["deployment_type"] = authCode.DeploymentType
		licenseFileData["max_activations"] = authCode.MaxActivations

		// 包含功能配置等
		if len(authCode.FeatureConfig) > 0 {
			var featureConfig map[string]interface{}
			if err := json.Unmarshal(authCode.FeatureConfig, &featureConfig); err == nil {
				licenseFileData["feature_config"] = featureConfig
			}
		}

		if len(authCode.UsageLimits) > 0 {
			var usageLimits map[string]interface{}
			if err := json.Unmarshal(authCode.UsageLimits, &usageLimits); err == nil {
				licenseFileData["usage_limits"] = usageLimits
			}
		}

		if len(authCode.CustomParameters) > 0 {
			var customParameters map[string]interface{}
			if err := json.Unmarshal(authCode.CustomParameters, &customParameters); err == nil {
				licenseFileData["custom_parameters"] = customParameters
			}
		}
	}

	// 序列化
	licenseJSON, err := json.Marshal(licenseFileData)
	if err != nil {
		return "", err
	}

	// 使用RSA数字签名
	encryptedData, err := s.signLicenseFile(licenseJSON)
	if err != nil {
		return "", err
	}

	return string(encryptedData), nil
}

// loadRSAPrivateKey 加载RSA私钥（懒加载）
func (s *licenseService) loadRSAPrivateKey() (*utils.RSAPrivateKey, error) {
	if s.rsaPrivateKey != nil {
		return s.rsaPrivateKey, nil
	}

	cfg := config.GetConfig()
	if cfg == nil {
		return nil, fmt.Errorf("configuration not initialized")
	}

	privateKeyPath := cfg.License.RSA.PrivateKeyPath
	if privateKeyPath == "" {
		return nil, fmt.Errorf("RSA private key path not configured")
	}

	key, err := utils.LoadRSAPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("加载RSA私钥失败: %w", err)
	}

	s.rsaPrivateKey = key
	return key, nil
}

// signLicenseFile 使用RSA数字签名许可证文件
func (s *licenseService) signLicenseFile(data []byte) ([]byte, error) {
	// 加载RSA私钥
	privateKey, err := s.loadRSAPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("加载RSA私钥失败: %w", err)
	}

	// 对数据进行数字签名
	signature, err := privateKey.SignData(data)
	if err != nil {
		return nil, fmt.Errorf("签名失败: %w", err)
	}

	// 构建签名后的许可证文件结构
	licenseWithSignature := map[string]interface{}{
		"data":      string(data), // 原始数据（JSON字符串）
		"signature": signature,    // 数字签名
		"algorithm": "RSA-PSS-SHA256",
	}

	// 序列化为JSON
	licenseJSON, err := json.Marshal(licenseWithSignature)
	if err != nil {
		return nil, fmt.Errorf("序列化失败: %w", err)
	}

	// 使用base64编码
	encoded := base64.StdEncoding.EncodeToString(licenseJSON)

	return []byte(encoded), nil
}

// GetStatsOverview 获取授权概览统计
func (s *licenseService) GetStatsOverview(ctx context.Context) (*models.StatsOverviewResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 获取配置信息
	cfg := config.GetConfig()
	offlineTimeout := time.Duration(cfg.License.HeartbeatTimeout) * time.Second

	// 计算时间点
	now := time.Now()
	onlineThreshold := now.Add(-offlineTimeout)
	lastMonth := now.AddDate(0, -1, 0)

	var stats models.StatsOverviewResponse

	// 执行统计查询
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 总授权码数量
		if err := tx.Model(&models.AuthorizationCode{}).Count(&stats.TotalAuthCodes).Error; err != nil {
			return err
		}

		// 2. 活跃许可证数量（状态为active）
		if err := tx.Model(&models.License{}).Where("status = ?", "active").Count(&stats.ActiveLicenses).Error; err != nil {
			return err
		}

		// 3. 即将过期数量（授权码在30天内过期，且状态为active）
		expireThreshold := now.AddDate(0, 0, 30)
		if err := tx.Model(&models.AuthorizationCode{}).
			Where("end_date <= ? AND end_date > ? AND is_locked = false", expireThreshold, now).
			Count(&stats.ExpiringSoon).Error; err != nil {
			return err
		}

		// 4. 异常告警数量（许可证超时未心跳）
		if err := tx.Model(&models.License{}).
			Where("status = 'active' AND (last_heartbeat < ? OR last_heartbeat IS NULL)", onlineThreshold).
			Count(&stats.AbnormalAlerts).Error; err != nil {
			return err
		}

		// 5. 计算同比上月的授权码和许可证数量（用于计算增长率）
		var lastMonthAuthCodes, lastMonthActiveLicenses int64

		// 上月同期授权码总数（一个月前同一时刻的累计总数）
		if err := tx.Model(&models.AuthorizationCode{}).
			Where("created_at <= ?", lastMonth).
			Count(&lastMonthAuthCodes).Error; err != nil {
			return err
		}

		// 上月同期活跃许可证总数（一个月前同一时刻的状态为active的许可证数）
		if err := tx.Model(&models.License{}).
			Where("status = ? AND created_at <= ?", "active", lastMonth).
			Count(&lastMonthActiveLicenses).Error; err != nil {
			return err
		}

		// 计算增长率（同比上月：当前总数相比一个月前总数的增长率）
		if lastMonthAuthCodes > 0 {
			stats.GrowthRate.AuthCodes = float64(stats.TotalAuthCodes-lastMonthAuthCodes) / float64(lastMonthAuthCodes) * 100
		} else {
			stats.GrowthRate.AuthCodes = 0
		}

		if lastMonthActiveLicenses > 0 {
			stats.GrowthRate.Licenses = float64(stats.ActiveLicenses-lastMonthActiveLicenses) / float64(lastMonthActiveLicenses) * 100
		} else {
			stats.GrowthRate.Licenses = 0
		}

		return nil
	})

	if err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	return &stats, nil
}
