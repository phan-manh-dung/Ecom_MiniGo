package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string     `json:"name" binding:"required" validate:"required"`
	Description string     `json:"description"`
	Price       float64    `json:"price" binding:"required" validate:"required,gt=0"`
	Image       string     `json:"image"`
	Inventory   *Inventory // 1-1: Mỗi product có một inventory
}

type Inventory struct {
	gorm.Model
	ProductID uint     // Foreign key
	Product   *Product // Con trỏ để dễ preload và tránh vòng lặp
	Quantity  int      `json:"quantity" validate:"required,gte=0"`
}
