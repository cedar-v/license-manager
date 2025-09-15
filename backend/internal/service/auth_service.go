package service

import (
	"context"
	"errors"
	"time"

	"license-manager/internal/config"
	"license-manager/internal/models"
	"license-manager/internal/repository"
	"license-manager/pkg/utils"

	"gorm.io/gorm"
)

type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// Login 用户登录
func (s *authService) Login(ctx context.Context, req *models.LoginRequest, clientIP string) (*models.LoginData, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return nil, errors.New("配置未初始化")
	}

	// 根据用户名查找用户
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, errors.New("查询用户失败")
	}

	// 检查账号状态
	if user.Status == "disabled" {
		return nil, errors.New("账号已被禁用")
	}

	// 检查账号是否被锁定
	if user.IsAccountLocked() {
		return nil, errors.New("账号已被锁定，请稍后重试")
	}

	// 验证密码
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		// 增加登录失败次数
		s.userRepo.IncrementLoginAttempts(ctx, user.ID)
		
		// 检查是否需要锁定账号
		if user.LoginAttempts >= 4 { // 第5次失败时锁定
			s.userRepo.LockUser(ctx, user.ID, 30) // 锁定30分钟
			return nil, errors.New("密码错误次数过多，账号已被锁定30分钟")
		}
		
		return nil, errors.New("用户名或密码错误")
	}

	// 登录成功，重置登录失败次数并更新登录信息
	user.LastLoginIP = &clientIP
	s.userRepo.ResetLoginAttempts(ctx, user.ID)

	// 如果账号状态是locked但锁定时间已过期，更新状态为active
	if user.Status == "locked" && (user.LockedUntil == nil || time.Now().After(*user.LockedUntil)) {
		user.Status = "active"
		s.userRepo.UpdateUser(ctx, user)
	}

	// 生成JWT Token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	// 计算过期时间（秒）
	expiresIn := cfg.Auth.JWT.ExpireHours * 3600

	return &models.LoginData{
		Token:     token,
		ExpiresIn: expiresIn,
		UserInfo: models.UserInfo{
			Username: user.Username,
			Role:     user.Role,
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
