package repository

import (
	"context"
	"math"
	"time"

	"gorm.io/gorm"

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
	if req.IsOnline != nil {
		// 在线状态：最后心跳时间在5分钟内
		fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
		if *req.IsOnline {
			query = query.Where("licenses.last_heartbeat >= ?", fiveMinutesAgo)
		} else {
			query = query.Where("licenses.last_heartbeat < ? OR licenses.last_heartbeat IS NULL", fiveMinutesAgo)
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
	for i, license := range licenses {
		// 计算在线状态
		isOnline := false
		if license.LastHeartbeat != nil {
			fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
			isOnline = license.LastHeartbeat.After(fiveMinutesAgo)
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
		fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
		license.IsOnline = license.LastHeartbeat.After(fiveMinutesAgo)
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
		fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
		license.IsOnline = license.LastHeartbeat.After(fiveMinutesAgo)
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