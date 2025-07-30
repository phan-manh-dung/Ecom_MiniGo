package repository

import (
	"gin/order_service/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetOrderByID(id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("OrderDetails").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetOrdersByUserID(userID uint) ([]*model.Order, error) {
	var orders []*model.Order
	err := r.db.Preload("OrderDetails").Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) CreateOrder(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) UpdateOrderStatus(orderID uint, status string) error {
	return r.db.Model(&model.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *OrderRepository) GetOrderDetails(orderID uint) ([]*model.OrderDetail, error) {
	var details []*model.OrderDetail
	err := r.db.Where("order_id = ?", orderID).Find(&details).Error
	if err != nil {
		return nil, err
	}
	return details, nil
}
