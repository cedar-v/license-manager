package service

import (
	"context"
	"fmt"
	"license-manager/internal/config"
	"license-manager/internal/models"
	"license-manager/internal/repository"
	pkgcontext "license-manager/pkg/context"
	"license-manager/pkg/i18n"
	"license-manager/pkg/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CuOrderService interface {
	CreateOrder(ctx context.Context, cuUserID, customerID string, req *models.CuOrderCreateRequest) (*models.CuOrder, error)
	CreatePendingOrder(ctx context.Context, cuUserID, customerID string, req *models.CuOrderCreateRequest) (*models.CuOrder, error)
	GetOrder(ctx context.Context, orderID, cuUserID string) (*models.CuOrder, error)
	GetUserOrders(ctx context.Context, cuUserID string, offset, limit int) ([]*models.CuOrder, int64, error)
	GetOrderSummary(ctx context.Context, customerID string) (*models.OrderSummaryResponse, error)
	CalculatePrice(ctx context.Context, packageID string, licenseCount int) (*PriceCalculationResult, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status string) error
}

type PriceCalculationResult struct {
	UnitPrice    float64 `json:"unit_price"`
	LicenseCount int     `json:"license_count"`
	DiscountRate float64 `json:"discount_rate"`
	TotalAmount  float64 `json:"total_amount"`
	DiscountDesc string  `json:"discount_description"`
}

// 套餐配置（临时写死，后续可改为配置化）
var packageConfigs = map[string]PackageConfig{
	"trial": {
		Name:       "试用版",
		Price:      0,
		MaxDevices: 1, // 试用版限制1个许可
	},
	"basic": {
		Name:       "基础版",
		Price:      300,
		MaxDevices: 1000, // 基础版最多1000个许可
	},
	"professional": {
		Name:       "专业版",
		Price:      2000,
		MaxDevices: 1000, // 专业版最多1000个许可
	},
}

type PackageConfig struct {
	Name       string
	Price      float64
	MaxDevices int
}

// 折扣规则配置
var discountRules = []DiscountRule{
	{MinQuantity: 50, MaxQuantity: 99, Rate: 0.9, Description: "50-99许可9折优惠"},
	{MinQuantity: 100, MaxQuantity: 499, Rate: 0.8, Description: "100-499许可8折优惠"},
	{MinQuantity: 500, MaxQuantity: 0, Rate: 0.7, Description: "500+许可7折优惠"},
}

type DiscountRule struct {
	MinQuantity int
	MaxQuantity int // 0表示无上限
	Rate        float64
	Description string
}

type cuOrderService struct {
	repo         repository.CuOrderRepository
	cuUserRepo   repository.CuUserRepository
	authCodeRepo repository.AuthorizationCodeRepository
	db           *gorm.DB
}

func NewCuOrderService(repo repository.CuOrderRepository, cuUserRepo repository.CuUserRepository, authCodeRepo repository.AuthorizationCodeRepository, db *gorm.DB) CuOrderService {
	return &cuOrderService{
		repo:         repo,
		cuUserRepo:   cuUserRepo,
		authCodeRepo: authCodeRepo,
		db:           db,
	}
}

func (s *cuOrderService) CreatePendingOrder(ctx context.Context, cuUserID, customerID string, req *models.CuOrderCreateRequest) (*models.CuOrder, error) {
	// 计算价格
	priceResult, err := s.CalculatePrice(ctx, req.PackageID, req.LicenseCount)
	if err != nil {
		return nil, err
	}

	// 生成订单号
	orderNo := s.generateOrderNo()

	// 获取套餐配置
	packageConfig, exists := packageConfigs[req.PackageID]
	if !exists {
		return nil, i18n.NewI18nError("600001", "套餐不存在")
	}

	// 创建订单（pending状态）
	order := &models.CuOrder{
		OrderNo:      orderNo,
		CustomerID:   customerID,
		CuUserID:     cuUserID,
		PackageID:    req.PackageID,
		PackageName:  packageConfig.Name,
		LicenseCount: req.LicenseCount,
		UnitPrice:    priceResult.UnitPrice,
		DiscountRate: priceResult.DiscountRate,
		TotalAmount:  priceResult.TotalAmount,
		Status:       "pending", // 待支付状态
	}

	// 对于试用版，设置过期时间
	if req.PackageID == "trial" {
		expiredAt := s.getTrialExpiryTime()
		order.ExpiredAt = &expiredAt
	}

	// 保存订单
	if err := s.repo.Create(order); err != nil {
		return nil, i18n.NewI18nError("601001", "订单创建失败: "+err.Error())
	}

	return order, nil
}

func (s *cuOrderService) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	// 这里需要实现更新订单状态的逻辑
	// 暂时先返回nil，后面再实现
	return nil
}

// getTrialExpiryTime 获取试用版过期时间（当月25日23:59:59）
func (s *cuOrderService) getTrialExpiryTime() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 25, 23, 59, 59, 0, now.Location())
}

func (s *cuOrderService) CreateOrder(ctx context.Context, cuUserID, customerID string, req *models.CuOrderCreateRequest) (*models.CuOrder, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 参数验证
	if req == nil || req.PackageID == "" || req.LicenseCount <= 0 {
		return nil, i18n.NewI18nError("900001", lang)
	}

	// 检查套餐是否存在
	pkg, exists := packageConfigs[req.PackageID]
	if !exists {
		return nil, i18n.NewI18nError("600001", lang) // 套餐不存在
	}

	// 验证许可数量
	if req.LicenseCount > pkg.MaxDevices {
		return nil, i18n.NewI18nError("601005", lang) // 许可数量超出限制
	}

	// 试用版特殊检查
	if req.PackageID == "trial" {
		// 试用版不支持批量购买
		if req.LicenseCount != 1 {
			return nil, i18n.NewI18nError("600004", lang) // 试用版不支持批量购买
		}

		// 检查购买时间限制
		now := time.Now()
		if now.Day() > 25 {
			return nil, i18n.NewI18nError("600003", lang) // 试用版购买时间限制
		}

		// 每个用户每月只能有一个试用版
		currentMonth := now.Format("2006-01")
		exists, err := s.repo.CheckTrialOrderExists(cuUserID, currentMonth)
		if err != nil {
			return nil, i18n.NewI18nError("900004", lang, err.Error())
		}
		if exists {
			return nil, i18n.NewI18nError("600005", lang) // 当月已购买试用版
		}
	}

	// 计算价格
	priceResult, err := s.CalculatePrice(ctx, req.PackageID, req.LicenseCount)
	if err != nil {
		return nil, err
	}

	// 生成订单号
	orderNo := s.generateOrderNo()

	// 创建订单
	order := &models.CuOrder{
		OrderNo:      orderNo,
		CustomerID:   customerID,
		CuUserID:     cuUserID,
		PackageID:    req.PackageID,
		PackageName:  pkg.Name,
		LicenseCount: req.LicenseCount,
		UnitPrice:    priceResult.UnitPrice,
		DiscountRate: priceResult.DiscountRate,
		TotalAmount:  priceResult.TotalAmount,
		Status:       "paid", // 直接设置为已支付状态
	}

	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建订单
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 支付成功，创建授权码
	// 实际项目中这里应该先调用支付接口，支付成功后再执行后续逻辑
	now := time.Now()
	order.UpdatedAt = now

	// 计算授权时间
	var startDate, endDate time.Time
	var expiredAt *time.Time
	if req.PackageID == "trial" {
		// 试用版：从今天到当月25日
		trialExpiry := time.Date(now.Year(), now.Month(), 25, 23, 59, 59, 0, now.Location())
		if now.Day() > 25 {
			// 如果今天已经过了25日，则无法购买
			tx.Rollback()
			return nil, i18n.NewI18nError("600003", lang) // 试用版购买时间限制
		}
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = trialExpiry
		expiredAt = &trialExpiry
	} else {
		// 永久版：从今天开始，结束时间设为很远的未来
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(1000, 0, 0) // 1000年后
		expiredAt = nil
	}

	// 生成自包含授权码
	var featureConfigMap, usageLimitsMap map[string]interface{}

	// 根据套餐类型设置配置信息
	if req.PackageID == "basic" {
		usageLimitsMap = map[string]interface{}{
			"type": "standard",
		}
	}

	licenseConfig := &utils.LicenseConfig{
		EndDate:       endDate,
		FeatureConfig: featureConfigMap,
		UsageLimits:   usageLimitsMap,
		CustomParameters: map[string]interface{}{
			"package_id":    req.PackageID,
			"package_name":  pkg.Name,
			"license_count": req.LicenseCount,
		},
	}

	// 从配置获取RSA私钥路径
	privateKeyPath := config.GetConfig().License.RSA.PrivateKeyPath
	authCode, err := utils.EncodeLicenseData(licenseConfig, privateKeyPath)
	if err != nil {
		tx.Rollback()
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 构建订单描述
	description := fmt.Sprintf("%s - %d个许可", pkg.Name, req.LicenseCount)

	// 构建授权码实体（使用生成的授权码）
	authCodeEntity := &models.AuthorizationCode{
		ID:             uuid.New().String(), // 生成新的UUID作为主键
		Code:           authCode,
		CustomerID:     customerID,
		CreatedBy:      cuUserID,
		Description:    &description,
		StartDate:      startDate,
		EndDate:        endDate,
		DeploymentType: "cloud",
		EncryptionType: &[]string{"standard"}[0],
		MaxActivations: req.LicenseCount,
		IsLocked:       false,
		// FeatureConfig 和 UsageLimits 已在授权码字符串中，无需重复存储
	}

	// 创建授权码
	if err := s.authCodeRepo.CreateAuthorizationCode(ctx, authCodeEntity); err != nil {
		tx.Rollback()
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 更新订单的授权码信息
	order.AuthorizationCode = &authCode
	order.ExpiredAt = expiredAt

	// 更新订单
	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	return order, nil
}

func (s *cuOrderService) GetOrder(ctx context.Context, orderID, cuUserID string) (*models.CuOrder, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	order, err := s.repo.GetByID(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, i18n.NewI18nError("601001", lang) // 订单不存在
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 检查权限：只能查看自己的订单
	if order.CuUserID != cuUserID {
		return nil, i18n.NewI18nError("100005", lang) // 权限不足
	}

	return order, nil
}

func (s *cuOrderService) GetUserOrders(ctx context.Context, cuUserID string, offset, limit int) ([]*models.CuOrder, int64, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 参数验证
	if cuUserID == "" {
		return nil, 0, i18n.NewI18nError("900001", lang)
	}

	orders, total, err := s.repo.GetByCuUserID(cuUserID, offset, limit)
	if err != nil {
		return nil, 0, i18n.NewI18nError("900004", lang, err.Error())
	}

	return orders, total, nil
}

// GetOrderSummary 获取订单汇总统计
func (s *cuOrderService) GetOrderSummary(ctx context.Context, customerID string) (*models.OrderSummaryResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 参数验证
	if customerID == "" {
		return nil, i18n.NewI18nError("900001", lang) // 业务错误，不覆盖多语言message
	}

	// 委托给Repository层进行数据访问
	return s.repo.GetCustomerOrderSummary(ctx, customerID)
}

func (s *cuOrderService) CalculatePrice(ctx context.Context, packageID string, licenseCount int) (*PriceCalculationResult, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 检查套餐是否存在
	pkg, exists := packageConfigs[packageID]
	if !exists {
		return nil, i18n.NewI18nError("600001", lang) // 套餐不存在
	}

	// 检查套餐是否已下架（虽然当前是写死配置，但保留这个检查）
	// TODO: 如果将来添加status字段，可以在这里检查
	// if pkg.Status != "active" {
	//     return nil, i18n.NewI18nError("600002", lang) // 套餐已下架
	// }

	// 验证许可数量
	if licenseCount > pkg.MaxDevices {
		return nil, i18n.NewI18nError("601005", lang) // 许可数量超出限制
	}

	// 计算折扣
	discountRate := 1.0
	discountDesc := "不享受折扣"
	for _, rule := range discountRules {
		if licenseCount >= rule.MinQuantity && (rule.MaxQuantity == 0 || licenseCount <= rule.MaxQuantity) {
			discountRate = rule.Rate
			discountDesc = rule.Description
			break
		}
	}

	// 计算总价
	totalAmount := pkg.Price * float64(licenseCount) * discountRate

	return &PriceCalculationResult{
		UnitPrice:    pkg.Price,
		LicenseCount: licenseCount,
		DiscountRate: discountRate,
		TotalAmount:  totalAmount,
		DiscountDesc: discountDesc,
	}, nil
}

// 生成订单号
func (s *cuOrderService) generateOrderNo() string {
	now := time.Now()
	dateStr := now.Format("20060102")

	// 使用纳秒时间戳的最后9位作为序列号，确保全局唯一性
	// 纳秒时间戳在同一纳秒内不会重复，重启应用也不会影响
	nano := now.UnixNano() % 1000000000 // 取最后9位

	return fmt.Sprintf("ORD%s%09d", dateStr, nano)
}
