package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	GetUserOrders(ctx context.Context, cuUserID string, req *models.CuOrderListRequest) ([]*models.CuOrder, int64, error)
	GetOrderSummary(ctx context.Context, customerID string) (*models.OrderSummaryResponse, error)
	CalculatePrice(ctx context.Context, packageID string, licenseCount int) (*PriceCalculationResult, error)
	CancelOrder(ctx context.Context, orderID, cuUserID string) (*models.CuOrder, error)
	DeleteOrder(ctx context.Context, orderID, cuUserID string) error
	UpdateOrderStatus(ctx context.Context, orderID string, status string) error
	ContinuePay(ctx context.Context, orderID, cuUserID, customerID, paymentMethod string) (*ContinuePayResponse, error)
}

// ContinuePayResponse 继续支付响应
type ContinuePayResponse struct {
	OrderID     string    `json:"order_id"`     // 订单ID
	OrderNo     string    `json:"order_no"`     // 订单号
	PaymentNo   string    `json:"payment_no"`   // 支付单号
	PaymentURL  *string   `json:"payment_url"`  // 支付链接
	TotalAmount float64   `json:"total_amount"` // 订单总金额
	ExpireTime  time.Time `json:"expire_time"`  // 支付过期时间
}

type PriceCalculationResult struct {
	UnitPrice    float64 `json:"unit_price"`
	LicenseCount int     `json:"license_count"`
	DiscountRate float64 `json:"discount_rate"`
	TotalAmount  float64 `json:"total_amount"`
	DiscountDesc string  `json:"discount_description"`
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
	packageRepo  repository.PackageRepository
	paymentRepo  repository.PaymentRepository
	invoiceRepo  repository.CuInvoiceRepository
	db           *gorm.DB
}

func NewCuOrderService(
	repo repository.CuOrderRepository,
	cuUserRepo repository.CuUserRepository,
	authCodeRepo repository.AuthorizationCodeRepository,
	packageRepo repository.PackageRepository,
	paymentRepo repository.PaymentRepository,
	invoiceRepo repository.CuInvoiceRepository,
	db *gorm.DB,
) CuOrderService {
	return &cuOrderService{
		repo:         repo,
		cuUserRepo:   cuUserRepo,
		authCodeRepo: authCodeRepo,
		packageRepo:  packageRepo,
		paymentRepo:  paymentRepo,
		invoiceRepo:  invoiceRepo,
		db:           db,
	}
}

func (s *cuOrderService) CreatePendingOrder(ctx context.Context, cuUserID, customerID string, req *models.CuOrderCreateRequest) (*models.CuOrder, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 计算价格
	priceResult, err := s.CalculatePrice(ctx, req.PackageID, req.LicenseCount)
	if err != nil {
		return nil, err
	}

	// 从数据库获取套餐配置
	pkgEntity, maxDevices, err := s.getPackageForOrder(ctx, req.PackageID)
	if err != nil {
		return nil, err
	}

	// 再次校验数量（防御性校验）
	if req.LicenseCount > maxDevices {
		return nil, i18n.NewI18nError("601005", lang)
	}

	// 生成订单号
	orderNo := s.generateOrderNo()

	// 创建订单（pending状态）
	order := &models.CuOrder{
		OrderNo:      orderNo,
		CustomerID:   customerID,
		CuUserID:     cuUserID,
		PackageID:    req.PackageID,
		PackageName:  pkgEntity.Name,
		LicenseCount: req.LicenseCount,
		UnitPrice:    priceResult.UnitPrice,
		DiscountRate: priceResult.DiscountRate,
		TotalAmount:  priceResult.TotalAmount,
		Status:       "pending", // 待支付状态
	}

	// 对于试用版，设置过期时间（根据套餐类型判断）
	if pkgEntity.Type == string(models.PackageTypeTrial) {
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

	// 检查套餐是否存在（从数据库）
	pkgEntity, maxDevices, err := s.getPackageForOrder(ctx, req.PackageID)
	if err != nil {
		return nil, err
	}

	// 验证许可数量
	if req.LicenseCount > maxDevices {
		return nil, i18n.NewI18nError("601005", lang) // 许可数量超出限制
	}

	// 试用版特殊检查（根据套餐类型判断）
	if pkgEntity.Type == string(models.PackageTypeTrial) {
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
		PackageName:  pkgEntity.Name,
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
	if pkgEntity.Type == string(models.PackageTypeTrial) {
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

	// 构建授权配置（用于后续生成产品激活码 payload）
	var featureConfigMap, usageLimitsMap map[string]interface{}

	// 根据套餐类型设置配置信息
	if pkgEntity.Type == string(models.PackageTypeBasic) {
		usageLimitsMap = map[string]interface{}{
			"type": "standard",
		}
	}

	customParametersMap := map[string]interface{}{
		"package_id":    req.PackageID,
		"package_name":  pkgEntity.Name,
		"license_count": req.LicenseCount,
	}

	// 授权码生成规则回退/统一为旧规则（不再自包含配置）
	authCode, err := utils.GenerateLegacyAuthorizationCode(customerID)
	if err != nil {
		tx.Rollback()
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 构建订单描述
	description := fmt.Sprintf("%s - %d个许可", pkgEntity.Name, req.LicenseCount)

	// 处理JSON字段（用于数据库存储与后续生成产品激活码 payload）
	var featureConfig, usageLimits, customParameters models.JSON

	b, err := json.Marshal(featureConfigMap)
	if err != nil {
		tx.Rollback()
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}
	featureConfig = models.JSON(b)

	if usageLimitsMap != nil {
		d, err := json.Marshal(usageLimitsMap)
		if err != nil {
			tx.Rollback()
			return nil, i18n.NewI18nError("900004", lang, err.Error())
		}
		usageLimits = models.JSON(d)
	}

	c, err := json.Marshal(customParametersMap)
	if err != nil {
		tx.Rollback()
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}
	customParameters = models.JSON(c)

	// 构建授权码实体（使用生成的授权码）
	authCodeEntity := &models.AuthorizationCode{
		ID:               uuid.New().String(), // 生成新的UUID作为主键
		Code:             authCode,
		CustomerID:       customerID,
		CreatedBy:        cuUserID,
		Description:      &description,
		StartDate:        startDate,
		EndDate:          endDate,
		DeploymentType:   "cloud",
		EncryptionType:   &[]string{"standard"}[0],
		MaxActivations:   req.LicenseCount,
		IsLocked:         false,
		FeatureConfig:    featureConfig,
		UsageLimits:      usageLimits,
		CustomParameters: customParameters,
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

func (s *cuOrderService) CancelOrder(ctx context.Context, orderID, cuUserID string) (*models.CuOrder, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	if orderID == "" || cuUserID == "" {
		return nil, i18n.NewI18nError("900001", lang)
	}

	order, err := s.repo.GetByID(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, i18n.NewI18nError("601001", lang) // 订单不存在
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 检查权限：只能取消自己的订单
	if order.CuUserID != cuUserID {
		return nil, i18n.NewI18nError("100005", lang) // 权限不足
	}

	switch order.Status {
	case "cancelled":
		// 幂等：已取消直接返回成功
		return order, nil
	case "paid":
		return nil, i18n.NewI18nError("601006", lang) // 订单已支付，无法取消
	case "pending":
		order.Status = "cancelled"
		if err := s.repo.Update(order); err != nil {
			return nil, i18n.NewI18nError("900004", lang, err.Error())
		}
		return order, nil
	default:
		// 兜底：未知状态按不可取消处理
		return nil, i18n.NewI18nError("601006", lang)
	}
}

func (s *cuOrderService) DeleteOrder(ctx context.Context, orderID, cuUserID string) error {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	if orderID == "" || cuUserID == "" {
		return i18n.NewI18nError("900001", lang)
	}

	order, err := s.repo.GetByID(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return i18n.NewI18nError("601001", lang) // 订单不存在
		}
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	// 检查权限：只能删除自己的订单
	if order.CuUserID != cuUserID {
		return i18n.NewI18nError("100005", lang) // 权限不足
	}

	// 开启事务，确保订单和相关支付记录一起删除
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除相关的支付记录
	if err := tx.Unscoped().Where("business_type = ? AND business_id = ?", models.BusinessTypePackageOrder, orderID).Delete(&models.Payment{}).Error; err != nil {
		tx.Rollback()
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	// 删除订单记录（物理删除）
	if err := tx.Unscoped().Where("id = ?", orderID).Delete(&models.CuOrder{}).Error; err != nil {
		tx.Rollback()
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return i18n.NewI18nError("900004", lang, err.Error())
	}

	return nil
}

func (s *cuOrderService) GetUserOrders(ctx context.Context, cuUserID string, req *models.CuOrderListRequest) ([]*models.CuOrder, int64, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 参数验证
	if cuUserID == "" {
		return nil, 0, i18n.NewI18nError("900001", lang)
	}

	if req == nil {
		req = &models.CuOrderListRequest{}
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	createdAtStart, createdAtEnd, err := parseCuOrderTimeRange(req.Time, time.Now())
	if err != nil {
		return nil, 0, i18n.NewI18nError("900001", lang, err.Error())
	}

	orders, total, err := s.repo.GetByCuUserID(ctx, cuUserID, req, createdAtStart, createdAtEnd)
	if err != nil {
		return nil, 0, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 为订单填充“是否已申请发票”标记
	if s.invoiceRepo != nil && len(orders) > 0 {
		for _, order := range orders {
			exists, err := s.invoiceRepo.CheckOrderInvoiceExists(order.ID)
			if err != nil {
				return nil, 0, i18n.NewI18nError("900004", lang, err.Error())
			}
			order.HasInvoiceApplied = exists
		}
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

	// 从数据库获取套餐配置
	pkgEntity, maxDevices, err := s.getPackageForOrder(ctx, packageID)
	if err != nil {
		return nil, err
	}

	// 验证许可数量
	if licenseCount > maxDevices {
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
	totalAmount := pkgEntity.Price * float64(licenseCount) * discountRate

	return &PriceCalculationResult{
		UnitPrice:    pkgEntity.Price,
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

// getPackageForOrder 从数据库获取套餐配置，并解析最大许可数量
func (s *cuOrderService) getPackageForOrder(ctx context.Context, packageID string) (*models.Package, int, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	if packageID == "" {
		return nil, 0, i18n.NewI18nError("600001", lang)
	}

	pkgEntity, err := s.packageRepo.GetByID(packageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, i18n.NewI18nError("600001", lang) // 套餐不存在
		}
		return nil, 0, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 检查套餐是否启用
	if pkgEntity.Status != 1 {
		return nil, 0, i18n.NewI18nError("600002", lang) // 套餐已下架
	}

	maxDevices := parseMaxDevicesFromFeatures(pkgEntity.Features)
	return pkgEntity, maxDevices, nil
}

// parseMaxDevicesFromFeatures 从 features JSON 中解析 "X个许可"
func parseMaxDevicesFromFeatures(features string) int {
	// 默认最大许可数，与原始写死配置保持一致
	maxDevices := 1000
	if features == "" {
		return maxDevices
	}

	var featureList []string
	if err := json.Unmarshal([]byte(features), &featureList); err != nil {
		return maxDevices
	}

	for _, f := range featureList {
		var n int
		if _, err := fmt.Sscanf(f, "%d个许可", &n); err == nil && n > 0 {
			maxDevices = n
			break
		}
	}

	return maxDevices
}

// ContinuePay 继续支付：获取订单的支付信息，如果支付单已过期则返回订单信息供重新创建支付单
func (s *cuOrderService) ContinuePay(ctx context.Context, orderID, cuUserID, customerID, paymentMethod string) (*ContinuePayResponse, error) {
	lang := pkgcontext.GetLanguageFromContext(ctx)

	// 1. 验证订单
	order, err := s.repo.GetByID(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, i18n.NewI18nError("601001", lang) // 订单不存在
		}
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 2. 检查权限：只能支付自己的订单
	if order.CuUserID != cuUserID {
		return nil, i18n.NewI18nError("100005", lang) // 权限不足
	}

	// 3. 检查订单状态：只有待支付订单才能继续支付
	if order.Status != "pending" {
		return nil, i18n.NewI18nError("601007", lang) // 订单状态不允许继续支付
	}

	// 4. 查询订单的最新支付单
	payment, err := s.paymentRepo.GetByBusinessID(ctx, models.BusinessTypePackageOrder, orderID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, i18n.NewI18nError("900004", lang, err.Error())
	}

	// 5. 如果支付单存在且可用（pending 且未过期），直接返回
	if payment != nil && err == nil {
		now := time.Now()
		if payment.Status == models.PaymentStatusPending && payment.ExpireTime.After(now) {
			return &ContinuePayResponse{
				OrderID:     order.ID,
				OrderNo:     order.OrderNo,
				PaymentNo:   payment.PaymentNo,
				PaymentURL:  payment.PaymentURL,
				TotalAmount: order.TotalAmount,
				ExpireTime:  payment.ExpireTime,
			}, nil
		}
	}

	// 6. 支付单不存在或已过期，返回订单信息（让 handler 层创建新的支付单）
	// 这里返回一个特殊标记，表示需要创建新的支付单
	return &ContinuePayResponse{
		OrderID:     order.ID,
		OrderNo:     order.OrderNo,
		PaymentNo:   "", // 空字符串表示需要创建新的支付单
		PaymentURL:  nil,
		TotalAmount: order.TotalAmount,
		ExpireTime:  time.Time{}, // 零值表示需要创建新的支付单
	}, nil
}
