package main

import (
	"time"
)

// User Models
type CreateUserRequest struct {
	Name   string `json:"name" example:"John Doe" doc:"User's full name"`
	SDT    string `json:"sdt" example:"0123456789" doc:"User's phone number"`
	RoleID uint32 `json:"role_id" example:"1" doc:"User's role ID"`
}

type UpdateUserRequest struct {
	Name string `json:"name" example:"John Doe Updated" doc:"Updated user name"`
	SDT  string `json:"sdt" example:"0987654321" doc:"Updated phone number"`
}

type UserResponse struct {
	ID        uint32    `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	SDT       string    `json:"sdt" example:"0123456789"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int32          `json:"total" example:"100"`
}

// Product Models
type CreateProductRequest struct {
	Name        string  `json:"name" example:"iPhone 15" doc:"Product name"`
	Description string  `json:"description" example:"Latest iPhone model" doc:"Product description"`
	Price       float32 `json:"price" example:"999.99" doc:"Product price"`
	Image       string  `json:"image" example:"https://example.com/iphone.jpg" doc:"Product image URL"`
	Inventory   uint32  `json:"inventory" example:"100" doc:"Available inventory"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" example:"iPhone 15 Pro"`
	Description string  `json:"description" example:"Latest iPhone Pro model"`
	Price       float32 `json:"price" example:"1199.99"`
	Image       string  `json:"image" example:"https://example.com/iphone-pro.jpg"`
	Inventory   uint32  `json:"inventory" example:"50"`
}

type ProductResponse struct {
	ID          uint32    `json:"id" example:"1"`
	Name        string    `json:"name" example:"iPhone 15"`
	Description string    `json:"description" example:"Latest iPhone model"`
	Price       float32   `json:"price" example:"999.99"`
	Image       string    `json:"image" example:"https://example.com/iphone.jpg"`
	Inventory   uint32    `json:"inventory" example:"100"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListProductsResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int32             `json:"total" example:"50"`
}

// Order Models
type CreateOrderRequest struct {
	UserID uint32      `json:"user_id" example:"1" doc:"User ID"`
	Items  []OrderItem `json:"items" doc:"Order items"`
}

type OrderItem struct {
	ProductID uint32  `json:"product_id" example:"1" doc:"Product ID"`
	Quantity  uint32  `json:"quantity" example:"2" doc:"Quantity to order"`
	UnitPrice float32 `json:"unit_price" example:"999.99" doc:"Unit price"`
}

type OrderResponse struct {
	ID          uint32      `json:"id" example:"1"`
	UserID      uint32      `json:"user_id" example:"1"`
	TotalPrice  float32     `json:"total_price" example:"1999.98"`
	OrderStatus int         `json:"status" example:"1"`
	Items       []OrderItem `json:"items"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type ListOrdersResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int32           `json:"total" example:"25"`
}

// Inventory Models
type DecreaseInventoryRequest struct {
	ProductID uint32 `json:"product_id" example:"1" doc:"Product ID"`
	Quantity  uint32 `json:"quantity" example:"5" doc:"Quantity to decrease"`
}

type IncreaseInventoryRequest struct {
	ProductID uint32 `json:"product_id" example:"1" doc:"Product ID"`
	Quantity  uint32 `json:"quantity" example:"5" doc:"Quantity to increase"`
}

// Common Response
type ErrorResponse struct {
	Error   string `json:"error" example:"Internal server error"`
	Message string `json:"message" example:"Something went wrong"`
}

type SuccessResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// Order status constants
const (
	OrderStatusPending   = 0
	OrderStatusCompleted = 1
	OrderStatusCancelled = 2
)
