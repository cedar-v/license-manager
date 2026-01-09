package database

import (
	"fmt"
	"log"
	"time"

	"license-manager/internal/config"
	"license-manager/internal/models"
	"license-manager/pkg/utils"
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
		&models.User{},
		&models.CuUser{},
		&models.CuOrder{},
		&models.Payment{},
		&models.AuthorizationCode{},
		&models.License{},
		&models.AuthorizationChange{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// 创建客户用户表的复合索引（GORM不支持复合索引，需要手动创建）
	if err := createCuUserIndexes(); err != nil {
		return fmt.Errorf("failed to create cu_user indexes: %w", err)
	}

	// 初始化客户编码序列
	if err := initCustomerCodeSequence(); err != nil {
		return fmt.Errorf("failed to initialize customer code sequence: %w", err)
	}

	// 初始化默认管理员用户
	if err := initDefaultAdminUser(); err != nil {
		return fmt.Errorf("failed to initialize default admin user: %w", err)
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

// createCuUserIndexes 创建客户用户表的复合索引
func createCuUserIndexes() error {
	// 检查是否支持复合索引（MySQL等）
	dbType := DB.Dialector.Name()

	switch dbType {
	case "mysql":
		// 注意：所有必要的索引都已在SQL迁移文件中创建
		// 这里不再重复创建，避免"Duplicate key name"错误
		log.Println("All cu_users indexes already created in SQL migration file")
	default:
		log.Printf("Skipping composite index creation for database type: %s", dbType)
	}

	return nil
}

// initDefaultAdminUser 初始化默认管理员用户
func initDefaultAdminUser() error {
	var user models.User
	result := DB.Where("username = ?", "admin").First(&user)

	if result.Error != nil {
		// 用户不存在，创建默认管理员
		passwordHash, err := utils.HashPassword("admin@123")
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		user = models.User{
			Username:     "admin",
			Email:        "admin@example.com",
			PasswordHash: passwordHash,
			FullName:     "系统管理员",
			Role:         "admin",
			Status:       "active",
		}

		if err := DB.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create default admin user: %w", err)
		}

		log.Println("Created default admin user: admin / admin@123")
	} else {
		log.Println("Default admin user already exists")
	}

	return nil
}
