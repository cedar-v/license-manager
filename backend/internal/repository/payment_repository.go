package repository

import (
	"context"
	"license-manager/internal/models"

	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *paymentRepository) GetByID(ctx context.Context, id int) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) GetByPaymentNo(ctx context.Context, paymentNo string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).Where("payment_no = ?", paymentNo).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) GetByBusinessID(ctx context.Context, businessType, businessID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).Where("business_type = ? AND business_id = ?", businessType, businessID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) GetByCustomerAndCuUserID(ctx context.Context, customerID, cuUserID string, offset, limit int) ([]*models.Payment, int64, error) {
	var payments []*models.Payment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Payment{}).Where("customer_id = ? AND cu_user_id = ?", customerID, cuUserID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&payments).Error
	if err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, paymentNo, status string, tradeNo *string, paymentTime *string) error {
	updateData := map[string]interface{}{
		"status": status,
	}

	if tradeNo != nil {
		updateData["trade_no"] = *tradeNo
	}

	if paymentTime != nil {
		updateData["payment_time"] = *paymentTime
	}

	return r.db.WithContext(ctx).Model(&models.Payment{}).Where("payment_no = ?", paymentNo).Updates(updateData).Error
}

func (r *paymentRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Payment{}).Error
}

func (r *paymentRepository) GetExpiredPayments(ctx context.Context) ([]*models.Payment, error) {
	var payments []*models.Payment

	err := r.db.WithContext(ctx).Where("status = ? AND expire_time < NOW()", models.PaymentStatusPending).Find(&payments).Error
	if err != nil {
		return nil, err
	}

	return payments, nil
}
