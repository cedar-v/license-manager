package routes

import (
	_ "license-manager/docs/swagger" // swagger docs
	"license-manager/internal/api/handlers"
	"license-manager/internal/api/middleware"
	"license-manager/internal/config"
	"license-manager/internal/database"
	"license-manager/internal/repository"
	"license-manager/internal/service"

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

	// 初始化服务层
	authService := service.NewAuthService(userRepo)
	systemService := service.NewSystemService()
	customerService := service.NewCustomerService(customerRepo)
	enumService := service.NewEnumService()

	// 初始化处理器层
	authHandler := handlers.NewAuthHandler(authService)
	systemHandler := handlers.NewSystemHandler(systemService)
	customerHandler := handlers.NewCustomerHandler(customerService)
	enumHandler := handlers.NewEnumHandler(enumService)

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
		}

		// 管理员接口
		admin := api.Group("/v1/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())
		{
			admin.GET("/system/info", systemHandler.GetSystemInfo)
		}
	}

	return router
}
