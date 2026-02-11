package routes

import (
	"fmt"
	_ "license-manager/docs/swagger" // swagger docs
	"license-manager/internal/api/handlers"
	"license-manager/internal/api/middleware"
	"license-manager/internal/config"
	"license-manager/internal/database"
	"license-manager/internal/repository"
	"license-manager/internal/service"
	"license-manager/pkg/cache"
	"license-manager/pkg/logger"
	"license-manager/pkg/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	router := gin.New()

	// 全局中间件
	router.Use(middleware.CustomLoggerMiddleware())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// 多语言中间件
	cfg := config.GetConfig()
	if cfg.I18n.Enable {
		i18nConfig := &middleware.I18nConfig{
			Enable:       cfg.I18n.Enable,
			DefaultLang:  cfg.I18n.DefaultLang,
			SupportLangs: cfg.I18n.SupportLangs,
		}
		router.Use(middleware.I18nMiddleware(i18nConfig))
	}

	// 初始化数据访问层
	db := database.GetDB()
	customerRepo := repository.NewCustomerRepository(db)
	userRepo := repository.NewUserRepository(db)
	cuUserRepo := repository.NewCuUserRepository(db)
	cuOrderRepo := repository.NewCuOrderRepository(db)
	authCodeRepo := repository.NewAuthorizationCodeRepository(db)
	licenseRepo := repository.NewLicenseRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	cuInvoiceRepo := repository.NewCuInvoiceRepository(db)
	adminInvoiceRepo := repository.NewAdminInvoiceRepository(db)
	packageRepo := repository.NewPackageRepository(db)
	leadRepo := repository.NewLeadRepository(db)

	// 获取logger实例
	log := logger.GetLogger()

	// 初始化缓存
	cacheInstance, err := cache.NewCache(cache.CacheConfig{
		Type:    cfg.Cache.Type,
		TTL:     cfg.Cache.TTL,
		Enabled: cfg.Cache.Enabled,
		Redis: cache.RedisConfig{
			Host:            cfg.Cache.Redis.Host,
			Port:            cfg.Cache.Redis.Port,
			Password:        cfg.Cache.Redis.Password,
			DB:              cfg.Cache.Redis.DB,
			PoolSize:        cfg.Cache.Redis.PoolSize,
			MinIdleConns:    cfg.Cache.Redis.MinIdleConns,
			ConnMaxIdleTime: cfg.Cache.Redis.ConnMaxIdleTime,
			DialTimeout:     cfg.Cache.Redis.DialTimeout,
			ReadTimeout:     cfg.Cache.Redis.ReadTimeout,
			WriteTimeout:    cfg.Cache.Redis.WriteTimeout,
		},
		Memory: cache.MemoryConfig{
			MaxSize: cfg.Cache.Memory.MaxSize,
		},
	})
	if err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}

	// 初始化SMS服务
	smsService, err := utils.NewSMSService(&cfg.SMS, cacheInstance, log)
	if err != nil {
		log.Fatalf("Failed to initialize SMS service: %v", err)
	}

	// 初始化服务层
	authService := service.NewAuthService(userRepo)
	systemService := service.NewSystemService()
	customerService := service.NewCustomerService(customerRepo)
	cuUserService := service.NewCuUserService(cuUserRepo, customerRepo, smsService, db)
	authCodeService := service.NewAuthorizationCodeService(authCodeRepo, customerRepo, cuUserRepo, licenseRepo)
	packageService := service.NewPackageService(packageRepo, db)
	cuOrderService := service.NewCuOrderService(cuOrderRepo, cuUserRepo, authCodeRepo, packageRepo, paymentRepo, cuInvoiceRepo, db)
	cuDeviceService := service.NewCuDeviceService(licenseRepo)
	paymentService := service.NewPaymentService(paymentRepo, cuOrderRepo, convertPaymentConfig(cfg.Payment), db)
	enumService := service.NewEnumService()
	licenseService := service.NewLicenseService(licenseRepo, db, log)
	dashboardService := service.NewDashboardService(dashboardRepo)

	// 初始化处理器层
	authHandler := handlers.NewAuthHandler(authService)
	systemHandler := handlers.NewSystemHandler(systemService)
	customerHandler := handlers.NewCustomerHandler(customerService)
	cuAuthHandler := handlers.NewCuAuthHandler(cuUserService)
	cuProfileHandler := handlers.NewCuProfileHandler(cuUserService)
	cuOrderHandler := handlers.NewCuOrderHandler(cuOrderService, paymentService, packageService)
	cuDeviceHandler := handlers.NewCuDeviceHandler(cuDeviceService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	cuAuthorizationHandler := handlers.NewCuAuthorizationHandler(authCodeService)
	cuInvoiceService := service.NewCuInvoiceService(cuInvoiceRepo, cuOrderRepo, db)
	cuInvoiceHandler := handlers.NewCuInvoiceHandler(cuInvoiceService)
	enumHandler := handlers.NewEnumHandler(enumService)
	authCodeHandler := handlers.NewAuthorizationCodeHandler(authCodeService)
	licenseHandler := handlers.NewLicenseHandler(licenseService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	adminInvoiceService := service.NewAdminInvoiceService(adminInvoiceRepo, cuOrderRepo, userRepo, cuUserRepo, db)
	adminInvoiceHandler := handlers.NewAdminInvoiceHandler(adminInvoiceService)
	packageHandler := handlers.NewPackageHandler(packageService)
	leadService := service.NewLeadService(leadRepo, db)
	leadHandler := handlers.NewLeadHandler(leadService)

	// 健康检测接口（无需认证）
	router.GET("/health", systemHandler.HealthCheck)

	// 静态文件服务 - 发票文件等
	// 发票PDF上传保存路径为 "../files/invoices"
	// 这里将URL前缀 /files 映射到 "../files" 目录，保证 /files/invoices 下的文件可被访问
	router.Static("/files", "../files")

	// Swagger文档路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	// API路由组
	api := router.Group("/api")
	{
		// 公开接口（无需认证）
		public := api.Group("")
		{
			public.POST("/v1/login", authHandler.Login)

			// 许可证激活接口（客户端软件使用）
			public.POST("/v1/activate", licenseHandler.ActivateLicense)
			public.POST("/v1/heartbeat", licenseHandler.Heartbeat)

			// 支付回调接口
			public.POST("/payment/alipay/callback", paymentHandler.AlipayCallback)
			public.POST("/v1/payment/alipay/callback", paymentHandler.AlipayCallback)

			// 公共发票下载接口
			public.GET("/public/invoices/download", adminInvoiceHandler.DownloadByToken)

			// 线索收集（无需鉴权）
			public.POST("/leads", leadHandler.CreateLead)
		}

		// 需要认证的接口
		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			// 认证相关
			auth.POST("/v1/logout", authHandler.Logout)
			auth.POST("/v1/auth/refresh", authHandler.RefreshToken)

			// 客户管理
			auth.GET("/customers", customerHandler.GetCustomerList)
			auth.GET("/customers/:id", customerHandler.GetCustomer)
			auth.POST("/customers", customerHandler.CreateCustomer)
			auth.PUT("/customers/:id", customerHandler.UpdateCustomer)
			auth.DELETE("/customers/:id", customerHandler.DeleteCustomer)
			auth.PATCH("/customers/:id/status", customerHandler.UpdateCustomerStatus)

			// 枚举管理
			auth.GET("/enums", enumHandler.GetAllEnums)
			auth.GET("/enums/:type", enumHandler.GetEnumsByType)

			// 授权码管理
			auth.GET("/v1/authorization-codes", authCodeHandler.GetAuthorizationCodeList)
			auth.POST("/v1/authorization-codes", authCodeHandler.CreateAuthorizationCode)
			auth.GET("/v1/authorization-codes/:id", authCodeHandler.GetAuthorizationCode)
			auth.GET("/v1/authorization-codes/:id/download", authCodeHandler.DownloadAuthorizationFile)
			auth.PUT("/v1/authorization-codes/:id", authCodeHandler.UpdateAuthorizationCode)
			auth.PUT("/v1/authorization-codes/:id/lock", authCodeHandler.LockUnlockAuthorizationCode)
			auth.DELETE("/v1/authorization-codes/:id", authCodeHandler.DeleteAuthorizationCode)
			auth.GET("/v1/authorization-codes/:id/changes", authCodeHandler.GetAuthorizationChangeList)

			// 许可证管理
			auth.GET("/v1/licenses", licenseHandler.GetLicenseList)
			auth.GET("/v1/licenses/:id", licenseHandler.GetLicense)
			auth.POST("/v1/licenses", licenseHandler.CreateLicense)
			auth.PUT("/v1/licenses/:id/revoke", licenseHandler.RevokeLicense)
			auth.GET("/v1/licenses/:id/download", licenseHandler.DownloadLicenseFile)

			// 统计分析
			auth.GET("/v1/stats/overview", licenseHandler.GetStatsOverview)

			// 仪表盘接口
			auth.GET("/v1/dashboard/authorization-trend", dashboardHandler.GetAuthorizationTrend)
			auth.GET("/v1/dashboard/recent-authorizations", dashboardHandler.GetRecentAuthorizations)

			// 发票管理（管理员）
			auth.GET("/v1/invoices", adminInvoiceHandler.GetAdminInvoices)
			auth.GET("/v1/invoices/:id", adminInvoiceHandler.GetAdminInvoiceDetail)
			auth.GET("/v1/invoices/summary", adminInvoiceHandler.GetAdminInvoiceSummary)
			auth.POST("/v1/invoices/:id/reject", adminInvoiceHandler.RejectInvoice)
			auth.POST("/v1/invoices/:id/issue", adminInvoiceHandler.IssueInvoice)
			auth.POST("/v1/invoices/upload", adminInvoiceHandler.UploadInvoiceFile)

			// 套餐管理（管理员）
			auth.GET("/packages", packageHandler.GetPackages)
			auth.GET("/packages/:id", packageHandler.GetPackageDetail)
			auth.POST("/packages", packageHandler.CreatePackage)
			auth.PUT("/packages/:id", packageHandler.UpdatePackage)
			auth.DELETE("/packages/:id", packageHandler.DeletePackage)
			auth.PUT("/packages/:id/status", packageHandler.UpdatePackageStatus)

			// 线索管理（管理员）
			auth.GET("/leads", leadHandler.GetLeads)
			auth.GET("/leads/summary", leadHandler.GetLeadSummary)
			auth.GET("/leads/:id", leadHandler.GetLead)
			auth.PUT("/leads/:id", leadHandler.UpdateLead)
			auth.PUT("/leads/:id/status", leadHandler.UpdateLeadStatus)
			auth.DELETE("/leads/:id", leadHandler.DeleteLead)
		}

		// 管理员接口
		admin := api.Group("/v1/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())
		{
			admin.GET("/system/info", systemHandler.GetSystemInfo)
		}
	}

	// 客户用户API路由组
	cuGroup := router.Group("/api/cu")
	{
		// 公开接口（无需认证）
		cuPublic := cuGroup.Group("")
		{
			cuPublic.POST("/register", cuAuthHandler.CuUserRegister)
			cuPublic.POST("/login", cuAuthHandler.CuUserLogin)
			cuPublic.POST("/send-login-sms", cuAuthHandler.CuUserSendLoginSms)
			cuPublic.POST("/send-register-sms", cuAuthHandler.CuUserSendRegisterSms)
			cuPublic.POST("/forgot-password", cuAuthHandler.CuUserForgotPassword)
			cuPublic.POST("/reset-password", cuAuthHandler.CuUserResetPassword)
			cuPublic.GET("/packages", cuOrderHandler.GetPackages)
		}

		// 需要认证的接口
		cuAuth := cuGroup.Group("")
		cuAuth.Use(middleware.CustomerAuth())
		{
			// 用户个人资料
			cuAuth.GET("/profile", cuProfileHandler.GetCuUserProfile)
			cuAuth.PUT("/profile", cuProfileHandler.UpdateCuUserProfile)
			cuAuth.POST("/profile/send-current-phone-sms", cuProfileHandler.CuUserSendCurrentPhoneSms)
			cuAuth.POST("/profile/send-new-phone-sms", cuProfileHandler.CuUserSendNewPhoneSms)
			cuAuth.PUT("/profile/phone", cuProfileHandler.UpdateCuUserPhone)
			cuAuth.PUT("/profile/password", cuProfileHandler.ChangeCuUserPassword)

			// 套餐和订单管理
			cuAuth.POST("/orders/calculate", cuOrderHandler.CalculatePrice)
			cuAuth.POST("/orders", cuOrderHandler.CreateOrder)
			cuAuth.GET("/orders/:order_id", cuOrderHandler.GetOrder)
			cuAuth.PUT("/orders/:order_id/cancel", cuOrderHandler.CancelOrder)
			cuAuth.POST("/orders/:order_id/pay", cuOrderHandler.ContinuePay)
			cuAuth.GET("/orders", cuOrderHandler.GetUserOrders)
			cuAuth.GET("/orders/summary", cuOrderHandler.GetOrderSummary)

			// 支付管理
			cuAuth.GET("/payment/:payment_no/status", paymentHandler.GetPaymentStatus)
			cuAuth.GET("/payment/history", paymentHandler.GetUserPayments)

			// 授权码管理
			cuAuth.POST("/authorization-codes/:codeId/share", cuAuthorizationHandler.ShareAuthorizationCode)
			cuAuth.POST("/authorization-codes/product-activation-code", cuAuthorizationHandler.GetProductActivationCode)
			cuAuth.GET("/authorization-codes", cuAuthorizationHandler.GetCuAuthorizationCodes)
			cuAuth.GET("/authorization-codes/summary", cuAuthorizationHandler.GetCuAuthorizationCodeSummary)

			// 设备管理
			cuAuth.GET("/devices", cuDeviceHandler.GetDevices)
			cuAuth.GET("/devices/summary", cuDeviceHandler.GetDeviceSummary)
			cuAuth.DELETE("/devices/:id", cuDeviceHandler.UnbindDevice)

			// 发票管理
			cuAuth.POST("/invoices", cuInvoiceHandler.CreateInvoice)
			cuAuth.GET("/invoices", cuInvoiceHandler.GetUserInvoices)
			cuAuth.GET("/invoices/:id", cuInvoiceHandler.GetUserInvoiceDetail)
			cuAuth.GET("/invoices/summary", cuInvoiceHandler.GetUserInvoiceSummary)
			cuAuth.GET("/invoices/:id/download", cuInvoiceHandler.DownloadInvoice)
		}
	}

	return router
}

// convertPaymentConfig 转换支付配置
func convertPaymentConfig(cfg config.PaymentConfig) *service.PaymentConfig {
	fmt.Printf("convertPaymentConfig: cfg.DefaultMethod=%s, cfg.Providers is nil: %v\n", cfg.DefaultMethod, cfg.Providers == nil)
	if cfg.Providers != nil {
		fmt.Printf("convertPaymentConfig: Found %d providers\n", len(cfg.Providers))
		for name, provider := range cfg.Providers {
			fmt.Printf("convertPaymentConfig: Provider %s, Enabled=%v\n", name, provider.Enabled)
		}
	}

	config := &service.PaymentConfig{
		DefaultMethod: cfg.DefaultMethod,
		ExpireMinutes: cfg.ExpireMinutes,
	}

	if cfg.Providers != nil {
		config.Providers = make(map[string]*service.AlipayProviderConfig)
		for name, provider := range cfg.Providers {
			if provider.Enabled {
				config.Providers[name] = &service.AlipayProviderConfig{
					Enabled:    provider.Enabled,
					AppID:      provider.AppID,
					PrivateKey: provider.PrivateKey,
					PublicKey:  provider.PublicKey,
					GatewayURL: provider.GatewayURL,
					NotifyURL:  provider.NotifyURL,
					ReturnURL:  provider.ReturnURL,
					SignType:   provider.SignType,
					Charset:    provider.Charset,
					Format:     provider.Format,
				}
			}
		}
	}

	fmt.Printf("convertPaymentConfig: Final config.Providers is nil: %v\n", config.Providers == nil)
	if config.Providers != nil {
		fmt.Printf("convertPaymentConfig: Final config has %d providers\n", len(config.Providers))
	}

	return config
}
