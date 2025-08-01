package service

import "license-manager/internal/models"

// AuthService 认证服务接口
type AuthService interface {
	Login(req *models.LoginRequest) (*models.LoginData, error)
	RefreshToken(token string) (string, error)
	ValidateToken(token string) error
}

// SystemService 系统服务接口
type SystemService interface {
	GetHealthStatus() *models.HealthResponse
	GetSystemInfo() map[string]interface{}
}