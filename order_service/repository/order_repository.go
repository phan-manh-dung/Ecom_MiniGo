package repository

import (
	"context"
	"gin/order_service/model"
	"gin/shared/generic"

	"gorm.io/gorm"
)

type OrderRepository struct {
	*generic.BaseRepository
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		BaseRepository: generic.NewBaseRepository(db),
	}
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, id uint) (*model.Order, error) {
	return generic.GenericGetByID[model.Order, uint](ctx, r.GetDB(), id)
}

func (r *OrderRepository) GetOrdersByUserID(userID uint) ([]*model.Order, error) {
	var orders []*model.Order
	result := r.GetDB().Preload("OrderDetails").Where("user_id = ?", userID).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *model.Order) error {
	return generic.GenericCreate[model.Order](ctx, r.GetDB(), order)
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID uint, status string) error {
	// Lấy order hiện tại
	order, err := r.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Cập nhật status
	order.Status = status
	return generic.GenericUpdate(ctx, r.GetDB(), order)
}

func (r *OrderRepository) GetOrderDetails(orderID uint) ([]*model.OrderDetail, error) {
	var details []*model.OrderDetail
	result := r.GetDB().Where("order_id = ?", orderID).Find(&details)
	if result.Error != nil {
		return nil, result.Error
	}
	return details, nil
}
