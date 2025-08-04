package database

import (
	"fmt"
	"log"
	"time"

	"license-manager/internal/config"
	"license-manager/internal/models"
)

// AutoMigrate 执行自动数据库迁移
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	log.Println("Starting database auto migration...")

	// 迁移数据库表结构
	err := DB.AutoMigrate(
		&models.Customer{},
		&models.CustomerCodeSequence{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// 初始化客户编码序列
	if err := initCustomerCodeSequence(); err != nil {
		return fmt.Errorf("failed to initialize customer code sequence: %w", err)
	}

	log.Println("Database auto migration completed successfully")
	return nil
}

// initCustomerCodeSequence 初始化客户编码序列
func initCustomerCodeSequence() error {
	currentYear := time.Now().Year()
	
	var sequence models.CustomerCodeSequence
	result := DB.Where("year = ?", currentYear).First(&sequence)
	
	if result.Error != nil {
		// 记录不存在，创建新记录
		sequence = models.CustomerCodeSequence{
			Year:           currentYear,
			SequenceNumber: 0,
		}
		if err := DB.Create(&sequence).Error; err != nil {
			return fmt.Errorf("failed to create customer code sequence: %w", err)
		}
		log.Printf("Initialized customer code sequence for year %d", currentYear)
	}
	
	return nil
}

// RunMigrationIfEnabled 根据配置决定是否运行迁移
func RunMigrationIfEnabled() error {
	cfg := config.GetConfig()
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	if cfg.Database.AutoMigrate {
		return AutoMigrate()
	}

	log.Println("Auto migration is disabled in config")
	return nil
}