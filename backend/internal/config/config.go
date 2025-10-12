package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Log      LogConfig      `mapstructure:"log"`
	I18n     I18nConfig     `mapstructure:"i18n"`
	License  LicenseConfig  `mapstructure:"license"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	Charset         string `mapstructure:"charset"`
	ParseTime       bool   `mapstructure:"parse_time"`
	Loc             string `mapstructure:"loc"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	AutoMigrate     bool   `mapstructure:"auto_migrate"`
}

type AuthConfig struct {
	JWT      JWTConfig      `mapstructure:"jwt"`
	Security SecurityConfig `mapstructure:"security"`
}

type JWTConfig struct {
	Secret                  string `mapstructure:"secret"`
	ExpireHours             int    `mapstructure:"expire_hours"`
	RefreshThresholdMinutes int    `mapstructure:"refresh_threshold_minutes"`
}

type SecurityConfig struct {
	MaxLoginAttempts       int `mapstructure:"max_login_attempts"`
	LockoutDurationMinutes int `mapstructure:"lockout_duration_minutes"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type I18nConfig struct {
	Enable       bool     `mapstructure:"enable"`        // 是否启用多语言
	DefaultLang  string   `mapstructure:"default_lang"`  // 默认语言
	ConfigPath   string   `mapstructure:"config_path"`   // 语言包路径
	SupportLangs []string `mapstructure:"support_langs"` // 支持的语言列表
}

type LicenseConfig struct {
	EncryptionKey    string `mapstructure:"encryption_key"`    // 许可证文件加密密钥
	HeartbeatTimeout int    `mapstructure:"heartbeat_timeout"` // 心跳超时时间(秒)
	OfflineTimeout   int    `mapstructure:"offline_timeout"`   // 离线超时时间(分钟)
	ExpiringDays     int    `mapstructure:"expiring_days"`     // 即将过期天数
}

var AppConfig *Config

func Load(configPath string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.AddConfigPath("./configs")
		viper.AddConfigPath(".")
	}

	// 设置默认值
	setDefaults()

	// 支持环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	log.Printf("Using config file: %s", viper.ConfigFileUsed())

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 18888)
	viper.SetDefault("server.mode", "debug")

	// Database defaults
	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.database", "license_manager")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.parse_time", true)
	viper.SetDefault("database.loc", "Local")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", 3600)
	viper.SetDefault("database.auto_migrate", true)

	// Auth defaults
	viper.SetDefault("auth.jwt.secret", "license-manager-default-secret-key")
	viper.SetDefault("auth.jwt.expire_hours", 1)
	viper.SetDefault("auth.jwt.refresh_threshold_minutes", 30)

	viper.SetDefault("auth.security.max_login_attempts", 5)
	viper.SetDefault("auth.security.lockout_duration_minutes", 15)

	// Log defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")

	// I18n defaults
	viper.SetDefault("i18n.enable", true)
	viper.SetDefault("i18n.default_lang", "zh-CN")
	viper.SetDefault("i18n.config_path", "../configs/i18n/errors")
	viper.SetDefault("i18n.support_langs", []string{"zh-CN", "en-US", "ja-JP"})

	// License defaults
	viper.SetDefault("license.encryption_key", "license-manager-secret-key-32b")
	viper.SetDefault("license.heartbeat_timeout", 300)
	viper.SetDefault("license.offline_timeout", 1440)
	viper.SetDefault("license.expiring_days", 30)
}

func GetConfig() *Config {
	return AppConfig
}
