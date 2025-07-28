package model

import (
	"gorm.io/gorm"
)

// Order đại diện cho 1 đơn hàng
type Order struct {
	gorm.Model
	UserID       uint
	TotalPrice   float64
	Status       string
	OrderDetails []*OrderDetail
}

type OrderDetail struct {
	gorm.Model
	OrderID   uint
	ProductID uint
	Quantity  uint
	UnitPrice float64
	Order     *Order
}
