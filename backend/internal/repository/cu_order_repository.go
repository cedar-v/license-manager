package repository

import (
	"license-manager/internal/models"

	"gorm.io/gorm"
)

type CuOrderRepository interface {
	Create(order *models.CuOrder) error
	GetByID(id string) (*models.CuOrder, error)
	GetByOrderNo(orderNo string) (*models.CuOrder, error)
	GetByCuUserID(cuUserID string, offset, limit int) ([]*models.CuOrder, int64, error)
	GetByCustomerID(customerID string, offset, limit int) ([]*models.CuOrder, int64, error)
	Update(order *models.CuOrder) error
	Delete(id string) error
	CheckTrialOrderExists(cuUserID string, currentMonth string) (bool, error)
}

type cuOrderRepository struct {
	db *gorm.DB
}

func NewCuOrderRepository(db *gorm.DB) CuOrderRepository {
	return &cuOrderRepository{db: db}
}

func (r *cuOrderRepository) Create(order *models.CuOrder) error {
	return r.db.Create(order).Error
}

func (r *cuOrderRepository) GetByID(id string) (*models.CuOrder, error) {
	var order models.CuOrder
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *cuOrderRepository) GetByOrderNo(orderNo string) (*models.CuOrder, error) {
	var order models.CuOrder
	err := r.db.Where("order_no = ? AND deleted_at IS NULL", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *cuOrderRepository) GetByCuUserID(cuUserID string, offset, limit int) ([]*models.CuOrder, int64, error) {
	var orders []*models.CuOrder
	var total int64

	query := r.db.Model(&models.CuOrder{}).Where("cu_user_id = ? AND deleted_at IS NULL", cuUserID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *cuOrderRepository) GetByCustomerID(customerID string, offset, limit int) ([]*models.CuOrder, int64, error) {
	var orders []*models.CuOrder
	var total int64

	query := r.db.Model(&models.CuOrder{}).Where("customer_id = ? AND deleted_at IS NULL", customerID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *cuOrderRepository) Update(order *models.CuOrder) error {
	return r.db.Save(order).Error
}

func (r *cuOrderRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.CuOrder{}).Error
}

// GetPendingOrders 移除：不再有pending状态的订单

// UpdateStatus 移除：不再需要更新订单状态

func (r *cuOrderRepository) CheckTrialOrderExists(cuUserID string, currentMonth string) (bool, error) {
	var count int64
	startOfMonth := currentMonth + "-01"
	endOfMonth := currentMonth + "-31"

	err := r.db.Model(&models.CuOrder{}).
		Where("cu_user_id = ? AND package_id = ? AND status = ? AND created_at >= ? AND created_at <= ? AND deleted_at IS NULL",
			cuUserID, "trial", "paid", startOfMonth, endOfMonth).
		Count(&count).Error

	return count > 0, err
}
