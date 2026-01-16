package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"license-manager/internal/models"
	"license-manager/internal/repository"
	"license-manager/pkg/utils"

	"github.com/smartwalle/alipay/v3"
	"gorm.io/gorm"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, req *models.PaymentCreateRequest) (*models.Payment, error)
	GetPayment(ctx context.Context, paymentNo string) (*models.Payment, error)
	GetPaymentStatus(ctx context.Context, paymentNo string) (*models.PaymentStatusResponse, error)
	ProcessAlipayCallback(ctx context.Context, values url.Values) error
	GetUserPayments(ctx context.Context, customerID, cuUserID string, page, pageSize int) (*models.PaymentListResponse, error)
	CancelExpiredPayments(ctx context.Context) error
}

type paymentService struct {
	paymentRepo  repository.PaymentRepository
	cuOrderRepo  repository.CuOrderRepository
	alipayClient *alipay.Client
	config       *PaymentConfig
	db           *gorm.DB
}

type PaymentConfig struct {
	DefaultMethod string                           `json:"default_method"`
	Providers     map[string]*AlipayProviderConfig `json:"providers"`
	ExpireMinutes int                              `json:"expire_minutes"`
}

type AlipayProviderConfig struct {
	Enabled    bool   `json:"enabled"`
	AppID      string `json:"app_id"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	GatewayURL string `json:"gateway_url"`
	NotifyURL  string `json:"notify_url"`
	ReturnURL  string `json:"return_url"`
	SignType   string `json:"sign_type"`
	Charset    string `json:"charset"`
	Format     string `json:"format"`
}

func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	cuOrderRepo repository.CuOrderRepository,
	config *PaymentConfig,
	db *gorm.DB,
) PaymentService {
	var alipayClient *alipay.Client

	if config != nil && config.Providers != nil {
		if alipayConfig, exists := config.Providers["alipay"]; exists && alipayConfig.Enabled {
			isProduction := strings.Contains(strings.ToLower(alipayConfig.GatewayURL), "openapi.alipay.com") &&
				!strings.Contains(strings.ToLower(alipayConfig.GatewayURL), "sandbox")

			client, err := alipay.New(
				alipayConfig.AppID,
				alipayConfig.PrivateKey,
				isProduction,
				alipay.WithSandboxGateway(alipayConfig.GatewayURL),
				alipay.WithProductionGateway(alipayConfig.GatewayURL),
			)
			if err != nil {
				fmt.Printf("Failed to create Alipay client: %v\n", err)
			} else {
				if alipayConfig.PublicKey != "" {
					if err := client.LoadAliPayPublicKey(alipayConfig.PublicKey); err != nil {
						fmt.Printf("Failed to load Alipay public key: %v\n", err)
					}
				}
				alipayClient = client
			}
		}
	}

	return &paymentService{
		paymentRepo:  paymentRepo,
		cuOrderRepo:  cuOrderRepo,
		alipayClient: alipayClient,
		config:       config,
		db:           db,
	}
}

func (s *paymentService) CreatePayment(ctx context.Context, req *models.PaymentCreateRequest) (*models.Payment, error) {
	// 设置默认值
	if req.Currency == "" {
		req.Currency = "CNY"
	}
	if s.config != nil && req.PaymentMethod == "" {
		req.PaymentMethod = s.config.DefaultMethod
	}

	// 生成支付单号
	paymentNo := utils.GeneratePaymentNo()

	// 计算过期时间
	expireMinutes := 30
	if s.config != nil && s.config.ExpireMinutes > 0 {
		expireMinutes = s.config.ExpireMinutes
	}
	expireTime := time.Now().Add(time.Duration(expireMinutes) * time.Minute)

	// 创建支付单
	payment := &models.Payment{
		PaymentNo:       paymentNo,
		BusinessType:    req.BusinessType,
		BusinessID:      req.BusinessID,
		CustomerID:      req.CustomerID,
		CuUserID:        req.CuUserID,
		Amount:          req.Amount,
		Currency:        req.Currency,
		PaymentMethod:   req.PaymentMethod,
		PaymentProvider: req.PaymentMethod,
		Status:          models.PaymentStatusPending,
		ExpireTime:      expireTime,
	}

	// 如果是支付宝，生成支付URL
	if req.PaymentMethod == models.PaymentMethodAlipay {
		if s.alipayClient == nil {
			return nil, fmt.Errorf("alipay client not initialized")
		}

		var notifyURL, returnURL string
		if s.config != nil && s.config.Providers != nil {
			if provider, ok := s.config.Providers["alipay"]; ok && provider != nil {
				notifyURL = strings.TrimSpace(provider.NotifyURL)
				returnURL = strings.TrimSpace(provider.ReturnURL)
			}
		}

		payRequest := alipay.TradePagePay{
			Trade: alipay.Trade{
				OutTradeNo:  paymentNo,
				TotalAmount: fmt.Sprintf("%.2f", req.Amount),
				Subject:     "产品套餐购买",
				ProductCode: "FAST_INSTANT_TRADE_PAY",
				NotifyURL:   notifyURL,
				ReturnURL:   returnURL,
			},
		}

		payURL, err := s.alipayClient.TradePagePay(payRequest)
		if err != nil {
			return nil, fmt.Errorf("generate payment URL failed: %w", err)
		}
		payURLStr := payURL.String()
		payment.PaymentURL = &payURLStr
	}

	// 保存到数据库
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("create payment failed: %w", err)
	}

	return payment, nil
}

func (s *paymentService) GetPayment(ctx context.Context, paymentNo string) (*models.Payment, error) {
	return s.paymentRepo.GetByPaymentNo(ctx, paymentNo)
}

func (s *paymentService) GetPaymentStatus(ctx context.Context, paymentNo string) (*models.PaymentStatusResponse, error) {
	payment, err := s.paymentRepo.GetByPaymentNo(ctx, paymentNo)
	if err != nil {
		return nil, err
	}

	response := &models.PaymentStatusResponse{
		PaymentNo:   payment.PaymentNo,
		Status:      payment.Status,
		Amount:      payment.Amount,
		PaymentTime: payment.PaymentTime,
		TradeNo:     payment.TradeNo,
	}

	// 如果支付成功，获取业务订单信息
	if payment.Status == models.PaymentStatusPaid && payment.BusinessID != nil {
		if payment.BusinessType == models.BusinessTypePackageOrder {
			order, err := s.cuOrderRepo.GetByID(*payment.BusinessID)
			if err == nil && order.AuthorizationCode != nil {
				response.BusinessOrder = &models.PaymentBusinessOrder{
					OrderNo:           order.OrderNo,
					AuthorizationCode: order.AuthorizationCode,
					Status:            order.Status,
				}
			}
		}
	}

	return response, nil
}

func (s *paymentService) ProcessAlipayCallback(ctx context.Context, values url.Values) error {
	if s.alipayClient == nil {
		return errors.New("alipay client not initialized")
	}
	if s.db == nil {
		return errors.New("db not initialized")
	}
	if values == nil || len(values) == 0 {
		return errors.New("empty callback values")
	}

	notification, err := s.alipayClient.DecodeNotification(values)
	if err != nil {
		return err
	}

	paymentNo := strings.TrimSpace(notification.OutTradeNo)
	if paymentNo == "" {
		return errors.New("missing out_trade_no")
	}

	notifyDataBytes, err := json.Marshal(values)
	if err != nil {
		return err
	}
	notifyDataStr := string(notifyDataBytes)

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var payment models.Payment
		if err := tx.Where("payment_no = ?", paymentNo).First(&payment).Error; err != nil {
			return err
		}

		// 金额二次校验
		expectedAmount := fmt.Sprintf("%.2f", payment.Amount)
		if notification.TotalAmount != "" && expectedAmount != notification.TotalAmount {
			return fmt.Errorf("payment amount mismatch: expected %s, got %s", expectedAmount, notification.TotalAmount)
		}

		updates := map[string]interface{}{
			"notify_data": notifyDataStr,
		}
		if strings.TrimSpace(notification.TradeNo) != "" {
			updates["trade_no"] = strings.TrimSpace(notification.TradeNo)
		}

		isPaid := false
		switch notification.TradeStatus {
		case alipay.TradeStatusSuccess, alipay.TradeStatusFinished:
			isPaid = true
			updates["status"] = models.PaymentStatusPaid
			if notification.GmtPayment != "" {
				if t, err := time.ParseInLocation("2006-01-02 15:04:05", notification.GmtPayment, time.Local); err == nil {
					updates["payment_time"] = t
				}
			}
		case alipay.TradeStatusClosed:
			// 已支付的订单不允许被回调覆盖为取消
			if payment.Status != models.PaymentStatusPaid {
				updates["status"] = models.PaymentStatusCancelled
			}
		default:
			// WAIT_BUYER_PAY / 未知状态：只记录回调数据，不改变业务状态
		}

		if err := tx.Model(&models.Payment{}).Where("payment_no = ?", paymentNo).Updates(updates).Error; err != nil {
			return err
		}

		if !isPaid {
			return nil
		}

		// 支付成功：根据业务类型处理
		if payment.BusinessType != models.BusinessTypePackageOrder || payment.BusinessID == nil || strings.TrimSpace(*payment.BusinessID) == "" {
			return nil
		}

		var order models.CuOrder
		if err := tx.Where("id = ? AND deleted_at IS NULL", *payment.BusinessID).First(&order).Error; err != nil {
			return err
		}

		// 幂等：订单已生成授权码则不重复创建
		if order.AuthorizationCode != nil && strings.TrimSpace(*order.AuthorizationCode) != "" && order.Status == "paid" {
			return nil
		}

		// 生成授权码并写入授权码表 & 订单（沿用当前 free order 的规则）
		authCode, err := utils.GenerateLegacyAuthorizationCode(order.CustomerID)
		if err != nil {
			return err
		}

		now := time.Now()

		var startDate, endDate time.Time
		var expiredAt *time.Time
		if order.PackageID == "trial" {
			trialExpiry := time.Date(now.Year(), now.Month(), 25, 23, 59, 59, 0, now.Location())
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			endDate = trialExpiry
			expiredAt = &trialExpiry
		} else {
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			endDate = startDate.AddDate(1000, 0, 0)
			expiredAt = nil
		}

		var featureConfigMap, usageLimitsMap map[string]interface{}
		if order.PackageID == "basic" {
			usageLimitsMap = map[string]interface{}{
				"type": "standard",
			}
		}
		customParametersMap := map[string]interface{}{
			"package_id":    order.PackageID,
			"package_name":  order.PackageName,
			"license_count": order.LicenseCount,
		}

		var featureConfig, usageLimits, customParameters models.JSON
		if featureConfigMap != nil {
			b, err := json.Marshal(featureConfigMap)
			if err != nil {
				return err
			}
			featureConfig = models.JSON(b)
		}
		if usageLimitsMap != nil {
			b, err := json.Marshal(usageLimitsMap)
			if err != nil {
				return err
			}
			usageLimits = models.JSON(b)
		}
		b, err := json.Marshal(customParametersMap)
		if err != nil {
			return err
		}
		customParameters = models.JSON(b)

		description := fmt.Sprintf("%s - %d个授权", order.PackageName, order.LicenseCount)
		encryptionType := "standard"
		authCodeEntity := &models.AuthorizationCode{
			Code:             authCode,
			CustomerID:       order.CustomerID,
			CreatedBy:        order.CuUserID,
			Description:      &description,
			StartDate:        startDate,
			EndDate:          endDate,
			DeploymentType:   "cloud",
			EncryptionType:   &encryptionType,
			MaxActivations:   order.LicenseCount,
			IsLocked:         false,
			FeatureConfig:    featureConfig,
			UsageLimits:      usageLimits,
			CustomParameters: customParameters,
		}
		if err := tx.Create(authCodeEntity).Error; err != nil {
			return err
		}

		orderUpdates := map[string]interface{}{
			"status":             "paid",
			"authorization_code": authCode,
			"expired_at":         expiredAt,
			"updated_at":         now,
		}
		if err := tx.Model(&models.CuOrder{}).Where("id = ?", order.ID).Updates(orderUpdates).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *paymentService) GetUserPayments(ctx context.Context, customerID, cuUserID string, page, pageSize int) (*models.PaymentListResponse, error) {
	offset := (page - 1) * pageSize
	payments, total, err := s.paymentRepo.GetByCustomerAndCuUserID(ctx, customerID, cuUserID, offset, pageSize)
	if err != nil {
		return nil, err
	}

	paymentResponses := make([]*models.PaymentResponse, len(payments))
	for i, payment := range payments {
		paymentResponses[i] = payment.ToResponse()
	}

	return &models.PaymentListResponse{
		Payments:   paymentResponses,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *paymentService) CancelExpiredPayments(ctx context.Context) error {
	expiredPayments, err := s.paymentRepo.GetExpiredPayments(ctx)
	if err != nil {
		return err
	}

	for _, payment := range expiredPayments {
		err = s.paymentRepo.UpdateStatus(ctx, payment.PaymentNo, models.PaymentStatusExpired, nil, nil)
		if err != nil {
			fmt.Printf("Failed to cancel expired payment %s: %v\n", payment.PaymentNo, err)
		}
	}

	return nil
}
