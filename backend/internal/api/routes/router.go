package routes

import (
	_ "license-manager/docs/swagger" // swagger docs
	"license-manager/internal/api/handlers"
	"license-manager/internal/api/middleware"
	"license-manager/internal/config"
	"license-manager/internal/database"
	"license-manager/internal/repository"
	"license-manager/internal/service"
	"license-manager/pkg/logger"

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

	// 获取logger实例
	log := logger.GetLogger()

	// 初始化服务层
	authService := service.NewAuthService(userRepo)
	systemService := service.NewSystemService()
	customerService := service.NewCustomerService(customerRepo)
	cuUserService := service.NewCuUserService(cuUserRepo, customerRepo, db)
	authCodeService := service.NewAuthorizationCodeService(authCodeRepo, customerRepo, licenseRepo)
	cuOrderService := service.NewCuOrderService(cuOrderRepo, cuUserRepo, authCodeRepo, db)
	enumService := service.NewEnumService()
	licenseService := service.NewLicenseService(licenseRepo, db, log)
	dashboardService := service.NewDashboardService(dashboardRepo)

	// 初始化处理器层
	authHandler := handlers.NewAuthHandler(authService)
	systemHandler := handlers.NewSystemHandler(systemService)
	customerHandler := handlers.NewCustomerHandler(customerService)
	cuAuthHandler := handlers.NewCuAuthHandler(cuUserService)
	cuProfileHandler := handlers.NewCuProfileHandler(cuUserService)
	cuOrderHandler := handlers.NewCuOrderHandler(cuOrderService)
	enumHandler := handlers.NewEnumHandler(enumService)
	authCodeHandler := handlers.NewAuthorizationCodeHandler(authCodeService)
	licenseHandler := handlers.NewLicenseHandler(licenseService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// 健康检测接口（无需认证）
	router.GET("/health", systemHandler.HealthCheck)

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
			cuAuth.PUT("/profile/phone", cuProfileHandler.UpdateCuUserPhone)
			cuAuth.PUT("/profile/password", cuProfileHandler.ChangeCuUserPassword)

			// 套餐和订单管理
			cuAuth.POST("/orders/calculate", cuOrderHandler.CalculatePrice)
			cuAuth.POST("/orders", cuOrderHandler.CreateOrder)
			cuAuth.GET("/orders/:order_id", cuOrderHandler.GetOrder)
			cuAuth.GET("/orders", cuOrderHandler.GetUserOrders)
		}
	}

	return router
}
