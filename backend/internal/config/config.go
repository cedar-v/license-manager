package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Auth   AuthConfig   `mapstructure:"auth"`
	Log    LogConfig    `mapstructure:"log"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type AuthConfig struct {
	JWT      JWTConfig      `mapstructure:"jwt"`
	Admin    AdminConfig    `mapstructure:"admin"`
	Security SecurityConfig `mapstructure:"security"`
}

type JWTConfig struct {
	Secret                   string `mapstructure:"secret"`
	ExpireHours              int    `mapstructure:"expire_hours"`
	RefreshThresholdMinutes  int    `mapstructure:"refresh_threshold_minutes"`
}

type AdminConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type SecurityConfig struct {
	MaxLoginAttempts        int `mapstructure:"max_login_attempts"`
	LockoutDurationMinutes  int `mapstructure:"lockout_duration_minutes"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
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
		log.Printf("Warning: Config file not found, using defaults: %v", err)
		log.Printf("Searched paths: %v", viper.ConfigFileUsed())
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

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

	// Auth defaults
	viper.SetDefault("auth.jwt.secret", "license-manager-default-secret-key")
	viper.SetDefault("auth.jwt.expire_hours", 1)
	viper.SetDefault("auth.jwt.refresh_threshold_minutes", 30)
	
	viper.SetDefault("auth.admin.username", "admin")
	viper.SetDefault("auth.admin.password", "admin@123")
	
	viper.SetDefault("auth.security.max_login_attempts", 5)
	viper.SetDefault("auth.security.lockout_duration_minutes", 15)

	// Log defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
}

func GetConfig() *Config {
	return AppConfig
}