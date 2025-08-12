package routes

import (
	_ "license-manager/docs/swagger" // swagger docs
	"license-manager/internal/api/handlers"
	"license-manager/internal/api/middleware"
	"license-manager/internal/config"
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

	// 初始化服务
	authService := service.NewAuthService()
	systemService := service.NewSystemService()

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(authService)
	systemHandler := handlers.NewSystemHandler(systemService)

	// 健康检测接口（无需认证）
	router.GET("/health", systemHandler.HealthCheck)

	// Swagger文档路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		// 公开接口（无需认证）
		public := v1.Group("")
		{
			public.POST("/login", authHandler.Login)
		}

		// 需要认证的接口
		auth := v1.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			// 认证相关
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/auth/refresh", authHandler.RefreshToken)
		}

		// 管理员接口
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())
		{
			admin.GET("/system/info", systemHandler.GetSystemInfo)
		}
	}

	return router
}
