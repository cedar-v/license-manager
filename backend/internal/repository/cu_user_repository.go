package repository

import (
	"errors"

	"license-manager/internal/models"

	"gorm.io/gorm"
)

type CuUserRepository interface {
	Create(user *models.CuUser) error
	GetByID(id string) (*models.CuUser, error)
	GetByPhone(phone, countryCode string) (*models.CuUser, error)
	GetByEmail(email string) (*models.CuUser, error)
	GetByCustomerID(customerID string, offset, limit int) ([]*models.CuUser, int64, error)
	Update(user *models.CuUser) error
	Delete(id string) error
	UpdateLoginInfo(id, ip string) error
	IncrementLoginAttempts(id string) error
	ResetLoginAttempts(id string) error
	LockAccount(id string, until interface{}) error
	CheckPhoneExists(phone, countryCode string, excludeID string) (bool, error)
	CheckEmailExists(email string, excludeID string) (bool, error)
}

type cuUserRepository struct {
	db *gorm.DB
}

func NewCuUserRepository(db *gorm.DB) CuUserRepository {
	return &cuUserRepository{db: db}
}

func (r *cuUserRepository) Create(user *models.CuUser) error {
	return r.db.Create(user).Error
}

func (r *cuUserRepository) GetByID(id string) (*models.CuUser, error) {
	var user models.CuUser
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *cuUserRepository) GetByPhone(phone, countryCode string) (*models.CuUser, error) {
	var user models.CuUser
	err := r.db.Where("phone = ? AND phone_country_code = ? AND deleted_at IS NULL", phone, countryCode).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *cuUserRepository) GetByEmail(email string) (*models.CuUser, error) {
	var user models.CuUser
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *cuUserRepository) GetByCustomerID(customerID string, offset, limit int) ([]*models.CuUser, int64, error) {
	var users []*models.CuUser
	var total int64

	query := r.db.Model(&models.CuUser{}).Where("customer_id = ? AND deleted_at IS NULL", customerID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *cuUserRepository) Update(user *models.CuUser) error {
	return r.db.Save(user).Error
}

func (r *cuUserRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.CuUser{}).Error
}

func (r *cuUserRepository) UpdateLoginInfo(id, ip string) error {
	return r.db.Model(&models.CuUser{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_at": gorm.Expr("NOW()"),
		"last_login_ip": ip,
	}).Error
}

func (r *cuUserRepository) IncrementLoginAttempts(id string) error {
	return r.db.Model(&models.CuUser{}).Where("id = ?", id).UpdateColumn("login_attempts", gorm.Expr("login_attempts + 1")).Error
}

func (r *cuUserRepository) ResetLoginAttempts(id string) error {
	return r.db.Model(&models.CuUser{}).Where("id = ?", id).Update("login_attempts", 0).Error
}

func (r *cuUserRepository) LockAccount(id string, until interface{}) error {
	return r.db.Model(&models.CuUser{}).Where("id = ?", id).Update("locked_until", until).Error
}

func (r *cuUserRepository) CheckPhoneExists(phone, countryCode string, excludeID string) (bool, error) {
	var count int64
	query := r.db.Model(&models.CuUser{}).Where("phone = ? AND phone_country_code = ? AND deleted_at IS NULL", phone, countryCode)
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *cuUserRepository) CheckEmailExists(email string, excludeID string) (bool, error) {
	var count int64
	query := r.db.Model(&models.CuUser{}).Where("email = ? AND deleted_at IS NULL", email)
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
