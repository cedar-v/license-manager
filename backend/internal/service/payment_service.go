package service

import (
	"context"
	"fmt"
	"license-manager/internal/models"
	"license-manager/internal/repository"
	"license-manager/pkg/utils"
	"time"

	"github.com/smartwalle/alipay/v3"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, req *models.PaymentCreateRequest) (*models.Payment, error)
	GetPayment(ctx context.Context, paymentNo string) (*models.Payment, error)
	GetPaymentStatus(ctx context.Context, paymentNo string) (*models.PaymentStatusResponse, error)
	ProcessPaymentCallback(ctx context.Context, notification interface{}) error
	GetUserPayments(ctx context.Context, customerID, cuUserID string, page, pageSize int) (*models.PaymentListResponse, error)
	CancelExpiredPayments(ctx context.Context) error
}

type paymentService struct {
	paymentRepo  repository.PaymentRepository
	cuOrderRepo  repository.CuOrderRepository
	alipayClient *alipay.Client
	config       *PaymentConfig
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
) PaymentService {
	var alipayClient *alipay.Client
	fmt.Printf("PaymentService: Checking providers, config.Providers is nil: %v\n", config.Providers == nil)
	if config.Providers != nil {
		fmt.Printf("PaymentService: Found %d providers\n", len(config.Providers))
		for providerName := range config.Providers {
			fmt.Printf("PaymentService: Provider: %s\n", providerName)
		}
		if alipayConfig, exists := config.Providers["alipay"]; exists && alipayConfig.Enabled {
			fmt.Printf("Initializing Alipay client with APPID: %s\n", alipayConfig.AppID)
			fmt.Printf("Private key length: %d\n", len(alipayConfig.PrivateKey))
			fmt.Printf("Public key length: %d\n", len(alipayConfig.PublicKey))

			// 使用官方SDK创建客户端
			fmt.Printf("Creating Alipay client with APPID: %s\n", alipayConfig.AppID)
			client, err := alipay.New(alipayConfig.AppID, alipayConfig.PrivateKey, false)
			if err != nil {
				fmt.Printf("Failed to create Alipay client: %v\n", err)
				fmt.Printf("Private key length: %d\n", len(alipayConfig.PrivateKey))
			} else {
				fmt.Printf("Alipay client created successfully\n")
				// 先不加载公钥，测试基本功能
				alipayClient = client
				fmt.Println("Alipay client initialized successfully")
			}
		} else {
			fmt.Println("Alipay provider not found or disabled")
		}
	} else {
		fmt.Println("No payment providers configured")
	}

	return &paymentService{
		paymentRepo:  paymentRepo,
		cuOrderRepo:  cuOrderRepo,
		alipayClient: alipayClient,
		config:       config,
	}
}

func (s *paymentService) CreatePayment(ctx context.Context, req *models.PaymentCreateRequest) (*models.Payment, error) {
	// 设置默认值
	if req.Currency == "" {
		req.Currency = "CNY"
	}
	if req.PaymentMethod == "" {
		req.PaymentMethod = s.config.DefaultMethod
	}

	// 生成支付单号
	paymentNo := utils.GeneratePaymentNo()

	// 计算过期时间
	expireTime := time.Now().Add(time.Duration(s.config.ExpireMinutes) * time.Minute)

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
	if req.PaymentMethod == "alipay" {
		fmt.Printf("Payment method is alipay, checking client...\n")
		if s.alipayClient == nil {
			fmt.Printf("Alipay client is nil!\n")
			return nil, fmt.Errorf("alipay client not initialized")
		}

		fmt.Printf("Alipay client is available, creating payment request...\n")

		// 使用官方SDK创建支付请求
		// TradePagePay继承了Trade，需要通过Trade字段设置
		payRequest := alipay.TradePagePay{
			Trade: alipay.Trade{
				OutTradeNo:  paymentNo,
				TotalAmount: fmt.Sprintf("%.2f", req.Amount),
				Subject:     "产品套餐购买",
				ProductCode: "FAST_INSTANT_TRADE_PAY",
			},
		}

		fmt.Printf("Creating Alipay payment request: OutTradeNo=%s, TotalAmount=%s\n", paymentNo, fmt.Sprintf("%.2f", req.Amount))

		// 生成支付URL
		payURL, err := s.alipayClient.TradePagePay(payRequest)
		if err != nil {
			fmt.Printf("TradePagePay failed: %v\n", err)
			return nil, fmt.Errorf("generate payment URL failed: %w", err)
		}
		payURLStr := payURL.String()
		fmt.Printf("Generated payment URL: %s\n", payURLStr[:100]+"...")
		payment.PaymentURL = &payURLStr
	} else {
		fmt.Printf("Payment method is: %s\n", req.PaymentMethod)
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

func (s *paymentService) ProcessPaymentCallback(ctx context.Context, notification interface{}) error {
	// 这里暂时简化处理回调逻辑
	// 实际项目中应该进行完整的验证和处理
	// 包括签名验证、状态检查、业务逻辑处理等

	// 暂时返回成功，表示回调处理完成
	return nil
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
