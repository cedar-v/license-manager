package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"license-manager/internal/api/routes"
	"license-manager/internal/config"
	"license-manager/internal/database"
	"license-manager/pkg/i18n"
	"license-manager/pkg/logger"

	"github.com/gin-gonic/gin"
)

// @title License Manager API
// @version 1.0
// @description 软件授权管理平台API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token for authentication

func main() {
	// 加载配置
	configPath := "config.yaml" // Docker 环境默认路径
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "../configs/config.yaml" // 本地开发环境
	}

	if err := config.Load(configPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger.Init()
	log := logger.GetLogger()

	log.Info("启动 License Manager 服务器...")

	cfg := config.GetConfig()

	// 初始化数据库（失败时重试）
	const (
		maxDBRetries  = 200
		retryInterval = 3 * time.Second
	)

	var dbErr error
	for attempt := 1; attempt <= maxDBRetries; attempt++ {
		dbErr = database.InitDatabase(&cfg.Database)
		if dbErr == nil {
			break
		}

		if attempt == maxDBRetries {
			log.Fatalf("数据库初始化失败: %v", dbErr)
		}

		log.Warnf("数据库初始化失败(第%d次): %v，将在 %v 后重试...", attempt, dbErr, retryInterval)
		time.Sleep(retryInterval)
	}

	defer database.Close()

	// 运行数据库迁移
	if err := database.RunMigrationIfEnabled(); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化多语言支持
	if cfg.I18n.Enable {
		if err := i18n.InitGlobalManager(cfg.I18n.ConfigPath, cfg.I18n.DefaultLang); err != nil {
			log.Fatalf("多语言初始化失败: %v", err)
		}

		// 预加载支持的语言
		for _, lang := range cfg.I18n.SupportLangs {
			if err := i18n.LoadLanguage(lang); err != nil {
				log.Warnf("加载语言包 %s 失败: %v", lang, err)
			} else {
				log.Infof("成功加载语言包: %s", lang)
			}
		}
		log.Info("多语言支持已启用")
	} else {
		log.Info("多语言支持已禁用")
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 设置路由
	router := routes.SetupRouter()

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Infof("服务器启动在 %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
