package main

import (
	"fmt"
	"log"
	"os"

	"license-manager/internal/api/routes"
	"license-manager/internal/config"
	"license-manager/internal/database"
	"license-manager/pkg/logger"
	"license-manager/pkg/i18n"

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

// @host localhost:18888
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token for authentication

func main() {
	// 加载配置 - 支持Docker环境和本地开发环境
	configPath := ""
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	} else {
		// 尝试不同的配置文件路径
		possiblePaths := []string{
			"config.yaml",              // Docker 环境
			"../configs/config.yaml",   // 本地开发环境
			"configs/config.yaml",      // 其他环境
		}
		
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				break
			}
		}
	}
	
	if err := config.Load(configPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger.Init()
	log := logger.GetLogger()

	log.Info("启动 License Manager 服务器...")

	cfg := config.GetConfig()

	// 初始化数据库
	if err := database.InitDatabase(&cfg.Database); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
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
