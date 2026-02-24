package repository

import (
	"context"
	"license-manager/internal/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

type CuOrderRepository interface {
	Create(order *models.CuOrder) error
	GetByID(id string) (*models.CuOrder, error)
	GetByOrderNo(orderNo string) (*models.CuOrder, error)
	GetByCuUserID(ctx context.Context, cuUserID string, req *models.CuOrderListRequest, createdAtStart, createdAtEnd *time.Time) ([]*models.CuOrder, int64, error)
	GetByCustomerID(customerID string, offset, limit int) ([]*models.CuOrder, int64, error)
	GetCustomerOrderSummary(ctx context.Context, customerID string) (*models.OrderSummaryResponse, error)
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

func (r *cuOrderRepository) GetByCuUserID(ctx context.Context, cuUserID string, req *models.CuOrderListRequest, createdAtStart, createdAtEnd *time.Time) ([]*models.CuOrder, int64, error) {
	var orders []*models.CuOrder
	var total int64

	// 默认值
	page := 1
	pageSize := 10
	search := ""
	status := ""
	if req != nil {
		if req.Page > 0 {
			page = req.Page
		}
		if req.PageSize > 0 {
			pageSize = req.PageSize
		}
		if pageSize > 100 {
			pageSize = 100
		}
		search = strings.TrimSpace(req.Search)
		status = strings.TrimSpace(req.Status)
	}
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx).Model(&models.CuOrder{}).Where("cu_user_id = ? AND deleted_at IS NULL", cuUserID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("(order_no LIKE ? OR authorization_code LIKE ?)", like, like)
	}

	if createdAtStart != nil && createdAtEnd != nil {
		query = query.Where("created_at >= ? AND created_at < ?", *createdAtStart, *createdAtEnd)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
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
	// 使用 Unscoped() 进行物理删除而不是软删除
	return r.db.Unscoped().Where("id = ?", id).Delete(&models.CuOrder{}).Error
}

// GetPendingOrders 移除：不再有pending状态的订单

// UpdateStatus 移除：不再需要更新订单状态

func (r *cuOrderRepository) CheckTrialOrderExists(cuUserID string, currentMonth string) (bool, error) {
	var count int64
	startOfMonth := currentMonth + "-01 00:00:00"
	endOfMonth := currentMonth + "-31 23:59:59"

	err := r.db.Model(&models.CuOrder{}).
		Where("cu_orders.cu_user_id = ? AND (cu_orders.status = 'paid' OR cu_orders.status = 'pending') AND cu_orders.created_at >= ? AND cu_orders.created_at <= ? AND cu_orders.deleted_at IS NULL",
			cuUserID, startOfMonth, endOfMonth).
		// We need to join with package to check package type because package_id might be a UUID, not 'trial'.
		// But if the order already recorded the type somewhere or if we just assume 'trial' is the package ID?
		// Actually, let's just use a subquery or join to verify package type is trial.
		Joins("JOIN packages ON packages.id = cu_orders.package_id").
		Where("packages.type = ?", string(models.PackageTypeTrial)).
		Count(&count).Error

	return count > 0, err
}

// GetCustomerOrderSummary 获取客户订单汇总统计
func (r *cuOrderRepository) GetCustomerOrderSummary(ctx context.Context, customerID string) (*models.OrderSummaryResponse, error) {
	// 查询订单总数
	var totalOrders int64
	err := r.db.Model(&models.CuOrder{}).
		Where("customer_id = ?", customerID).
		Count(&totalOrders).Error
	if err != nil {
		return nil, err
	}

	// 查询待支付订单数
	var pendingOrders int64
	err = r.db.Model(&models.CuOrder{}).
		Where("customer_id = ? AND status = ?", customerID, "pending").
		Count(&pendingOrders).Error
	if err != nil {
		return nil, err
	}

	// 查询已支付订单数
	var paidOrders int64
	err = r.db.Model(&models.CuOrder{}).
		Where("customer_id = ? AND status = ?", customerID, "paid").
		Count(&paidOrders).Error
	if err != nil {
		return nil, err
	}

	return &models.OrderSummaryResponse{
		TotalOrders:   totalOrders,
		PendingOrders: pendingOrders,
		PaidOrders:    paidOrders,
	}, nil
}
