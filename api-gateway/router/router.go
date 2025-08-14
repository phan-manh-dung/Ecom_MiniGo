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

// SetupRoutes thiết lập tất cả routes với middleware phù hợp
func (r *Router) SetupRoutes(
	userHandler *handler.UserServiceClient,
	productHandler *handler.ProductServiceClient,
	orderHandler *handler.OrderServiceClient,
) *gin.Engine {

	// 1. Setup public routes (không cần auth)
	r.setupPublicRoutes(userHandler)

	// 2. Setup protected routes (cần auth)
	r.setupProtectedRoutes(userHandler, productHandler, orderHandler)

	return r.engine
}

// setupPublicRoutes - các routes không cần authentication
func (r *Router) setupPublicRoutes(userHandler *handler.UserServiceClient) {
	// Health check endpoint
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public auth routes (không cần auth)
	authGroup := r.engine.Group("/api/auth")
	{
		authGroup.POST("/login", userHandler.Login)
		authGroup.POST("/register", userHandler.CreateUser)
	}
}

// setupProtectedRoutes - các routes cần authentication
func (r *Router) setupProtectedRoutes(
	userHandler *handler.UserServiceClient,
	productHandler *handler.ProductServiceClient,
	orderHandler *handler.OrderServiceClient,
) {
	// Tạo protected group với auth middleware
	protectedGroup := r.engine.Group("/api")
	protectedGroup.Use(middleware.NewAuthMiddleware())
	{
		// User routes
		userGroup := protectedGroup.Group("/users")
		{
			userGroup.GET("/:id", userHandler.GetUser)
			userGroup.PUT("/:id", userHandler.UpdateUser)
			userGroup.DELETE("/:id", userHandler.DeleteUser)
			userGroup.GET("/sdt/:id", userHandler.GetUserBySDT)
		}

		// Product routes
		productGroup := protectedGroup.Group("/products")
		{
			productGroup.GET("/:id", productHandler.GetProduct)
			productGroup.POST("/:id/inventory/decrease", productHandler.DecreaseInventory)
		}

		// Order routes
		orderGroup := protectedGroup.Group("/orders")
		{
			orderGroup.POST("/", orderHandler.CreateOrder)
			orderGroup.GET("/:id", orderHandler.GetOrder)
			orderGroup.PUT("/:id/status", orderHandler.UpdateOrderStatus)
			orderGroup.DELETE("/:id", orderHandler.CancelOrder)
			orderGroup.GET("/user/:userId", orderHandler.GetOrdersByUser)
		}

		// Admin routes với role check
		adminGroup := protectedGroup.Group("/admin")
		adminGroup.Use(middleware.RequireRoleMiddleware("ADMIN"))
		{
			adminGroup.GET("/users", userHandler.ListUsers)
			adminGroup.PUT("/products/:id", productHandler.UpdateProduct)
			adminGroup.DELETE("/products/:id", productHandler.DeleteProduct)
			adminGroup.POST("/products/", productHandler.CreateProduct)

			// Report routes
			reportHandler := handler.NewReportHandler()
			adminGroup.POST("/reports/generate", reportHandler.GenerateReport)
		}
	}
}
