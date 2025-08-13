package router

import (
	"gin/api-gateway/handler"
	"gin/api-gateway/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter() *Router {
	return &Router{
		engine: gin.New(),
	}
}

// SetupRoutes thiết lập routes với middleware và handlers
func (r *Router) SetupRoutes(
	userHandler *handler.UserServiceClient,
	productHandler *handler.ProductServiceClient,
	orderHandler *handler.OrderServiceClient,
) *gin.Engine {

	// 1. Áp dụng middleware theo thứ tự (từ ngoài vào trong)
	r.engine.Use(middleware.NewCORSMiddleware())      // CORS middleware (outermost)
	r.engine.Use(middleware.NewRequestIDMiddleware()) // Request ID middleware
	r.engine.Use(middleware.NewAuthMiddleware())      // Authentication middleware (innermost)

	// 2. Setup routes với handlers đã được inject
	r.setupUserRoutes(userHandler)
	r.setupProductRoutes(productHandler)
	r.setupOrderRoutes(orderHandler)

	return r.engine
}

func (r *Router) setupUserRoutes(userHandler *handler.UserServiceClient) {
	userGroup := r.engine.Group("/api/users")
	{
		userGroup.POST("/", userHandler.CreateUser)
		userGroup.GET("/:id", userHandler.GetUser)
		userGroup.PUT("/:id", userHandler.UpdateUser)
		userGroup.DELETE("/:id", userHandler.DeleteUser)
		userGroup.GET("/", userHandler.ListUsers)
	}
}

func (r *Router) setupProductRoutes(productHandler *handler.ProductServiceClient) {
	productGroup := r.engine.Group("/api/products")
	{
		productGroup.POST("/", productHandler.CreateProduct)
		productGroup.GET("/:id", productHandler.GetProduct)
		productGroup.PUT("/:id", productHandler.UpdateProduct)
		productGroup.DELETE("/:id", productHandler.DeleteProduct)
		productGroup.POST("/:id/inventory/decrease", productHandler.DecreaseInventory)
		//	productGroup.POST("/:id/inventory/increase", productHandler.IncreaseInventory)
	}
}

func (r *Router) setupOrderRoutes(orderHandler *handler.OrderServiceClient) {
	orderGroup := r.engine.Group("/api/orders")
	{
		orderGroup.POST("/", orderHandler.CreateOrder)
		orderGroup.GET("/:id", orderHandler.GetOrder)
		orderGroup.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		orderGroup.DELETE("/:id", orderHandler.CancelOrder)
		orderGroup.GET("/user/:userId", orderHandler.GetOrdersByUser)
	}
}
