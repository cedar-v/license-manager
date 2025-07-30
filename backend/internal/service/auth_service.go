package service

import (
	"errors"

	"license-manager/internal/config"
	"license-manager/internal/models"
	"license-manager/pkg/utils"
)

type authService struct{}

// NewAuthService 创建认证服务实例
func NewAuthService() AuthService {
	return &authService{}
}

// Login 用户登录
func (s *authService) Login(req *models.LoginRequest) (*models.LoginData, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return nil, errors.New("配置未初始化")
	}

	// 验证用户名和密码
	if req.Username != cfg.Auth.Admin.Username || req.Password != cfg.Auth.Admin.Password {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成JWT Token (用户ID设为1，因为只有一个管理员用户)
	token, err := utils.GenerateToken(1, req.Username, "administrator")
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	// 计算过期时间（秒）
	expiresIn := cfg.Auth.JWT.ExpireHours * 3600

	return &models.LoginData{
		Token:     token,
		ExpiresIn: expiresIn,
		UserInfo: models.UserInfo{
			Username: req.Username,
			Role:     "administrator",
		},
	}, nil
}

// RefreshToken 刷新Token
func (s *authService) RefreshToken(token string) (string, error) {
	// 解析当前Token
	claims, err := utils.ParseToken(token)
	if err != nil {
		return "", err
	}

	// 生成新Token
	newToken, err := utils.GenerateToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return "", errors.New("生成新令牌失败")
	}

	return newToken, nil
}

// ValidateToken 验证Token
func (s *authService) ValidateToken(token string) error {
	_, err := utils.ValidateToken(token)
	return err
}