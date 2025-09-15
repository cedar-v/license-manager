package repository

import (
	"context"
	"time"

	"license-manager/internal/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// GetUserByUsername 根据用户名获取用户
func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据ID获取用户
func (r *userRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// IncrementLoginAttempts 增加登录失败次数
func (r *userRepository) IncrementLoginAttempts(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", id).
		UpdateColumn("login_attempts", gorm.Expr("login_attempts + 1")).Error
}

// ResetLoginAttempts 重置登录失败次数
func (r *userRepository) ResetLoginAttempts(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"login_attempts": 0,
			"last_login_at":  &now,
			"locked_until":   nil,
		}).Error
}

// LockUser 锁定用户账号
func (r *userRepository) LockUser(ctx context.Context, id string, lockDuration int) error {
	lockUntil := time.Now().Add(time.Duration(lockDuration) * time.Minute)
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       "locked",
			"locked_until": &lockUntil,
		}).Error
}