package repository

import (
	"context"
	"math"
	"time"

	"gorm.io/gorm"

	"license-manager/internal/config"
	"license-manager/internal/models"
)

type licenseRepository struct {
	db *gorm.DB
}

// NewLicenseRepository 创建许可证数据访问实例
func NewLicenseRepository(db *gorm.DB) LicenseRepository {
	return &licenseRepository{
		db: db,
	}
}

// GetLicenseList 查询许可证列表
func (r *licenseRepository) GetLicenseList(ctx context.Context, req *models.LicenseListRequest) (*models.LicenseListResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	if req.Sort == "" {
		req.Sort = "created_at"
	}
	if req.Order == "" {
		req.Order = "desc"
	}

	// 构建查询
	query := r.db.Model(&models.License{}).
		Select(`licenses.id, licenses.license_key, licenses.authorization_code_id, 
				authorization_codes.code as authorization_code, customers.customer_name,
				licenses.hardware_fingerprint, licenses.status, licenses.activation_ip,
				licenses.last_online_ip, licenses.activated_at, licenses.last_heartbeat`).
		Joins("LEFT JOIN authorization_codes ON licenses.authorization_code_id = authorization_codes.id").
		Joins("LEFT JOIN customers ON licenses.customer_id = customers.id")

	// 授权码ID筛选
	if req.AuthorizationCodeID != "" {
		query = query.Where("licenses.authorization_code_id = ?", req.AuthorizationCodeID)
	}

	// 客户ID筛选
	if req.CustomerID != "" {
		query = query.Where("licenses.customer_id = ?", req.CustomerID)
	}

	// 状态筛选
	if req.Status != "" {
		query = query.Where("licenses.status = ?", req.Status)
	}

	// 在线状态筛选
	if req.IsOnline != nil && *req.IsOnline != "" {
		cfg := config.GetConfig()
		heartbeatTimeoutSeconds := cfg.License.HeartbeatTimeout
		if heartbeatTimeoutSeconds <= 0 {
			heartbeatTimeoutSeconds = 300 // 默认5分钟
		}
		onlineThreshold := time.Now().Add(-time.Duration(heartbeatTimeoutSeconds) * time.Second)

		if *req.IsOnline == "true" {
			query = query.Where("licenses.last_heartbeat >= ?", onlineThreshold)
		} else if *req.IsOnline == "false" {
			query = query.Where("licenses.last_heartbeat < ? OR licenses.last_heartbeat IS NULL", onlineThreshold)
		}
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 排序和分页
	orderBy := "licenses." + req.Sort + " " + req.Order
	offset := (req.Page - 1) * req.PageSize

	var licenses []struct {
		ID                  string     `json:"id"`
		LicenseKey          string     `json:"license_key"`
		AuthorizationCodeID string     `json:"authorization_code_id"`
		AuthorizationCode   string     `json:"authorization_code"`
		CustomerName        string     `json:"customer_name"`
		HardwareFingerprint string     `json:"hardware_fingerprint"`
		Status              string     `json:"status"`
		ActivationIP        *string    `json:"activation_ip"`
		LastOnlineIP        *string    `json:"last_online_ip"`
		ActivatedAt         *time.Time `json:"activated_at"`
		LastHeartbeat       *time.Time `json:"last_heartbeat"`
	}

	if err := query.Order(orderBy).Limit(req.PageSize).Offset(offset).Scan(&licenses).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]models.LicenseListItem, len(licenses))

	cfg := config.GetConfig()
	heartbeatTimeoutSeconds := cfg.License.HeartbeatTimeout
	if heartbeatTimeoutSeconds <= 0 {
		heartbeatTimeoutSeconds = 300 // 默认5分钟
	}
	onlineThreshold := time.Now().Add(-time.Duration(heartbeatTimeoutSeconds) * time.Second)

	for i, license := range licenses {
		// 计算在线状态
		isOnline := false
		if license.LastHeartbeat != nil {
			isOnline = license.LastHeartbeat.After(onlineThreshold)
		}

		// 格式化时间
		var activatedAtStr, lastHeartbeatStr *string
		if license.ActivatedAt != nil {
			str := license.ActivatedAt.Format(time.RFC3339)
			activatedAtStr = &str
		}
		if license.LastHeartbeat != nil {
			str := license.LastHeartbeat.Format(time.RFC3339)
			lastHeartbeatStr = &str
		}

		list[i] = models.LicenseListItem{
			ID:                  license.ID,
			LicenseKey:          license.LicenseKey,
			AuthorizationCodeID: license.AuthorizationCodeID,
			AuthorizationCode:   license.AuthorizationCode,
			CustomerName:        license.CustomerName,
			HardwareFingerprint: license.HardwareFingerprint,
			Status:              license.Status,
			IsOnline:            isOnline,
			ActivationIP:        license.ActivationIP,
			LastOnlineIP:        license.LastOnlineIP,
			ActivatedAt:         activatedAtStr,
			LastHeartbeat:       lastHeartbeatStr,
		}
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &models.LicenseListResponse{
		List:       list,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetLicenseByID 根据ID获取许可证信息
func (r *licenseRepository) GetLicenseByID(ctx context.Context, id string) (*models.License, error) {
	var license models.License

	err := r.db.Preload("AuthorizationCode").Preload("Customer").
		Where("id = ?", id).First(&license).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrLicenseNotFound
		}
		return nil, err
	}

	// 计算在线状态
	if license.LastHeartbeat != nil {
		cfg := config.GetConfig()
		heartbeatTimeoutSeconds := cfg.License.HeartbeatTimeout
		if heartbeatTimeoutSeconds <= 0 {
			heartbeatTimeoutSeconds = 300 // 默认5分钟
		}
		onlineThreshold := time.Now().Add(-time.Duration(heartbeatTimeoutSeconds) * time.Second)
		license.IsOnline = license.LastHeartbeat.After(onlineThreshold)
	}

	return &license, nil
}

// CreateLicense 创建许可证
func (r *licenseRepository) CreateLicense(ctx context.Context, license *models.License) error {
	return r.db.Create(license).Error
}

// UpdateLicense 更新许可证信息
func (r *licenseRepository) UpdateLicense(ctx context.Context, license *models.License) error {
	return r.db.Save(license).Error
}

// CheckAuthorizationCodeExists 检查授权码是否存在
func (r *licenseRepository) CheckAuthorizationCodeExists(ctx context.Context, authCodeID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.AuthorizationCode{}).Where("id = ?", authCodeID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAuthorizationCodeByID 根据ID获取授权码信息
func (r *licenseRepository) GetAuthorizationCodeByID(ctx context.Context, authCodeID string) (*models.AuthorizationCode, error) {
	var authCode models.AuthorizationCode
	err := r.db.Where("id = ?", authCodeID).First(&authCode).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrAuthorizationCodeNotFound
		}
		return nil, err
	}
	return &authCode, nil
}

// GetAuthorizationCodeByCode 根据授权码获取授权码信息
func (r *licenseRepository) GetAuthorizationCodeByCode(ctx context.Context, code string) (*models.AuthorizationCode, error) {
	var authCode models.AuthorizationCode
	err := r.db.Where("code = ?", code).First(&authCode).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrAuthorizationCodeNotFound
		}
		return nil, err
	}
	return &authCode, nil
}

// GetLicenseByKey 根据许可证密钥获取许可证信息
func (r *licenseRepository) GetLicenseByKey(ctx context.Context, licenseKey string) (*models.License, error) {
	var license models.License

	err := r.db.Preload("AuthorizationCode").Preload("Customer").
		Where("license_key = ?", licenseKey).First(&license).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrLicenseNotFound
		}
		return nil, err
	}

	// 计算在线状态
	if license.LastHeartbeat != nil {
		cfg := config.GetConfig()
		heartbeatTimeoutSeconds := cfg.License.HeartbeatTimeout
		if heartbeatTimeoutSeconds <= 0 {
			heartbeatTimeoutSeconds = 300 // 默认5分钟
		}
		onlineThreshold := time.Now().Add(-time.Duration(heartbeatTimeoutSeconds) * time.Second)
		license.IsOnline = license.LastHeartbeat.After(onlineThreshold)
	}

	return &license, nil
}

// GetActiveLicenseCount 获取指定授权码的激活许可证数量
func (r *licenseRepository) GetActiveLicenseCount(ctx context.Context, authCodeID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.License{}).
		Where("authorization_code_id = ? AND status = ?", authCodeID, "active").
		Count(&count).Error
	return count, err
}

// GetCustomerDeviceList 查询客户设备列表（关联授权码信息）
func (r *licenseRepository) GetCustomerDeviceList(ctx context.Context, customerID string, req *models.DeviceListRequest) (*models.DeviceListResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 构建查询
	query := r.db.Model(&models.License{}).
		Select(`licenses.id, licenses.device_info, licenses.last_online_ip,
				licenses.activated_at, licenses.last_heartbeat,
				authorization_codes.code as authorization_code,
				authorization_codes.id as authorization_code_id,
				authorization_codes.end_date, authorization_codes.description`).
		Joins("LEFT JOIN authorization_codes ON licenses.authorization_code_id = authorization_codes.id").
		Where("licenses.customer_id = ? AND licenses.status = ? AND licenses.deleted_at IS NULL",
			customerID, "active")

	// 获取心跳超时配置（秒）
	cfg := config.GetConfig()
	heartbeatTimeoutSeconds := cfg.License.HeartbeatTimeout
	if heartbeatTimeoutSeconds <= 0 {
		heartbeatTimeoutSeconds = 300 // 默认5分钟
	}

	now := time.Now()
	onlineThreshold := now.Add(-time.Duration(heartbeatTimeoutSeconds) * time.Second)

	// 按授权码ID筛选
	if req.AuthorizationCodeID != "" {
		query = query.Where("licenses.authorization_code_id = ?", req.AuthorizationCodeID)
	}

	// 设备名称模糊搜索
	if req.DeviceName != "" {
		query = query.Where("JSON_EXTRACT(licenses.device_info, '$.name') LIKE ?", "%"+req.DeviceName+"%")
	}

	// 在线状态筛选
	if req.IsOnline != nil && *req.IsOnline != "" {
		if *req.IsOnline == "true" {
			query = query.Where("licenses.last_heartbeat > ?", onlineThreshold)
		} else if *req.IsOnline == "false" {
			query = query.Where("(licenses.last_heartbeat <= ? OR licenses.last_heartbeat IS NULL)", onlineThreshold)
		}
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 定义扫描结构体
	type deviceScanResult struct {
		ID                  string     `json:"id"`
		DeviceInfo          string     `json:"device_info"`
		LastOnlineIP        *string    `json:"last_online_ip"`
		ActivatedAt         *time.Time `json:"activated_at"`
		LastHeartbeat       *time.Time `json:"last_heartbeat"`
		AuthorizationCode   string     `json:"authorization_code"`
		AuthorizationCodeID string     `json:"authorization_code_id"`
		EndDate             time.Time  `json:"end_date"`
		Description         *string    `json:"description"`
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	var scanResults []deviceScanResult
	err := query.Order("licenses.activated_at DESC").
		Offset(offset).Limit(req.PageSize).
		Scan(&scanResults).Error
	if err != nil {
		return nil, err
	}

	devices := make([]models.DeviceListItem, len(scanResults))

	for i, result := range scanResults {
		// 解析device_info JSON
		var deviceInfo map[string]interface{}
		if result.DeviceInfo != "" {
			// 这里需要JSON解析，暂时简化处理
			deviceInfo = map[string]interface{}{"raw": result.DeviceInfo}
		}

		// 计算在线状态（使用心跳超时配置）
		isOnline := false
		if result.LastHeartbeat != nil && result.LastHeartbeat.After(onlineThreshold) {
			isOnline = true
		}

		// 格式化时间字段
		var activatedAtStr, lastHeartbeatStr *string
		if result.ActivatedAt != nil {
			str := result.ActivatedAt.Format(time.RFC3339)
			activatedAtStr = &str
		}
		if result.LastHeartbeat != nil {
			str := result.LastHeartbeat.Format(time.RFC3339)
			lastHeartbeatStr = &str
		}

		// 构建授权信息
		description := ""
		if result.Description != nil {
			description = *result.Description
		}

		devices[i] = models.DeviceListItem{
			ID:            result.ID,
			DeviceInfo:    deviceInfo,
			IsOnline:      isOnline,
			LastOnlineIP:  result.LastOnlineIP,
			LastHeartbeat: lastHeartbeatStr,
			ActivatedAt:   activatedAtStr,
			AuthorizationInfo: models.AuthorizationInfo{
				AuthorizationCode:   result.AuthorizationCode,
				AuthorizationCodeID: result.AuthorizationCodeID,
				EndDate:             result.EndDate.Format(time.RFC3339),
				Description:         description,
			},
		}
	}

	return &models.DeviceListResponse{
		List:     devices,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// DeleteLicenseByID 根据ID删除许可证（物理删除，用于设备解绑）
func (r *licenseRepository) DeleteLicenseByID(ctx context.Context, id string) error {
	return r.db.Unscoped().Delete(&models.License{}, "id = ?", id).Error
}

// GetCustomerDeviceSummary 获取客户设备汇总统计
func (r *licenseRepository) GetCustomerDeviceSummary(ctx context.Context, customerID string) (*models.DeviceSummaryResponse, error) {
	// 获取心跳超时配置（秒）
	cfg := config.GetConfig()
	heartbeatTimeoutSeconds := cfg.License.HeartbeatTimeout
	if heartbeatTimeoutSeconds <= 0 {
		heartbeatTimeoutSeconds = 300 // 默认5分钟
	}

	// 查询设备总数
	var totalDevices int64
	err := r.db.Model(&models.License{}).
		Where("customer_id = ? AND status = ? AND deleted_at IS NULL", customerID, "active").
		Count(&totalDevices).Error
	if err != nil {
		return nil, err
	}

	// 查询在线设备数
	var onlineDevices int64
	err = r.db.Model(&models.License{}).
		Where("customer_id = ? AND status = ? AND deleted_at IS NULL AND last_heartbeat > ?",
			customerID, "active", time.Now().Add(-time.Duration(heartbeatTimeoutSeconds)*time.Second)).
		Count(&onlineDevices).Error
	if err != nil {
		return nil, err
	}

	// 计算离线设备数
	offlineDevices := totalDevices - onlineDevices

	return &models.DeviceSummaryResponse{
		TotalDevices:   totalDevices,
		OnlineDevices:  onlineDevices,
		OfflineDevices: offlineDevices,
	}, nil
}

// CheckLicenseBelongsToCustomer 检查许可证是否属于指定客户
func (r *licenseRepository) CheckLicenseBelongsToCustomer(ctx context.Context, licenseID, customerID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.License{}).
		Where("id = ? AND customer_id = ? AND deleted_at IS NULL", licenseID, customerID).
		Count(&count).Error
	return count > 0, err
}
