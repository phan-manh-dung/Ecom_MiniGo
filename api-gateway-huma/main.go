package main

import (
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Khởi tạo handlers (mock services for now)
	handlers := NewHandlers(nil)

	// Tạo Chi router và Huma API
	r := chi.NewRouter()
	api := humachi.New(r, huma.DefaultConfig("E-commerce API Gateway (Huma)", "1.0.0"))

	// ==================== USER ROUTES ====================
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/users", Summary: "Create a new user"}, handlers.CreateUserHandler)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/users/{id}", Summary: "Get user by ID"}, handlers.GetUserHandler)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/users/{id}", Summary: "Update user"}, handlers.UpdateUserHandler)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/users/{id}", Summary: "Delete user"}, handlers.DeleteUserHandler)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/users", Summary: "List all users"}, handlers.ListUsersHandler)

	// ==================== PRODUCT ROUTES ====================
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/products", Summary: "Create a new product"}, handlers.CreateProductHandler)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/products/{id}", Summary: "Get product by ID"}, handlers.GetProductHandler)
	huma.Register(api, huma.Operation{Method: http.MethodPut, Path: "/products/{id}", Summary: "Update product"}, handlers.UpdateProductHandler)
	huma.Register(api, huma.Operation{Method: http.MethodDelete, Path: "/products/{id}", Summary: "Delete product"}, handlers.DeleteProductHandler)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/products", Summary: "List all products"}, handlers.ListProductsHandler)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/products/inventory/decrease", Summary: "Decrease product inventory"}, handlers.DecreaseInventoryHandler)

	// ==================== ORDER ROUTES ====================
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/orders", Summary: "Create a new order"}, handlers.CreateOrderHandler)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/orders/{id}", Summary: "Get order by ID"}, handlers.GetOrderHandler)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/orders/{id}/cancel", Summary: "Cancel order"}, handlers.CancelOrderHandler)
	huma.Register(api, huma.Operation{Method: http.MethodGet, Path: "/orders", Summary: "List all orders"}, handlers.ListOrdersHandler)

	// Khởi động server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Starting Huma API Gateway on port %s", port)
	log.Printf("API Documentation available at: http://localhost:%s/docs", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
