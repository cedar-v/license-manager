package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"license-manager/internal/api/routes"
	"license-manager/internal/config"
	"license-manager/pkg/logger"
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
	// 加载配置
	if err := config.Load(""); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger.Init()
	log := logger.GetLogger()

	log.Info("启动 License Manager 服务器...")

	cfg := config.GetConfig()
	
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