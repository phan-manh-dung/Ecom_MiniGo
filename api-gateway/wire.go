//go:build wireinject
// +build wireinject

package main

import (
	"gin/api-gateway/handler"
	"gin/api-gateway/middleware"
	"gin/api-gateway/router"
	"gin/api-gateway/service_manager"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// wireApp khởi tạo toàn bộ application với dependency injection
func wireApp() (*gin.Engine, error) {
	wire.Build(
		// 1. Tạo Service Manager
		service_manager.NewServiceManager,

		// 2. Tạo gRPC Clients từ Service Manager
		provideUserClient,
		provideProductClient,
		provideOrderClient,

		// 3. Tạo Handlers
		handler.NewUserServiceClient,
		handler.NewProductServiceClient,
		handler.NewOrderServiceClient,

		// 4. Tạo Middleware
		middleware.NewCORSMiddleware,
		middleware.NewRequestIDMiddleware,
		middleware.NewLoggingMiddleware,
		middleware.NewAuthMiddleware,

		// 5. Tạo Router
		router.NewRouter,

		// 6. Tạo App với tất cả dependencies
		provideApp,
	)
	return &gin.Engine{}, nil
}

// provideUserClient tạo User gRPC client
func provideUserClient(serviceManager *service_manager.ServiceManager) interface{} {
	return serviceManager.UserClient
}

// provideProductClient tạo Product gRPC client
func provideProductClient(serviceManager *service_manager.ServiceManager) interface{} {
	return serviceManager.ProductClient
}

// provideOrderClient tạo Order gRPC client
func provideOrderClient(serviceManager *service_manager.ServiceManager) interface{} {
	return serviceManager.OrderClient
}

// provideApp tạo gin engine với tất cả dependencies
func provideApp(
	userHandler *handler.UserServiceClient,
	productHandler *handler.ProductServiceClient,
	orderHandler *handler.OrderServiceClient,
	corsMiddleware gin.HandlerFunc,
	requestIDMiddleware gin.HandlerFunc,
	loggingMiddleware gin.HandlerFunc,
	authMiddleware gin.HandlerFunc,
	router *router.Router,
) *gin.Engine {

	// 1. Setup routes trước (không có middleware)
	engine := router.SetupRoutes(userHandler, productHandler, orderHandler)

	// 2. Sau đó áp dụng middleware theo thứ tự (từ ngoài vào trong)
	engine.Use(corsMiddleware)      // CORS middleware (outermost)
	engine.Use(requestIDMiddleware) // Request ID middleware
	engine.Use(loggingMiddleware)   // Logging middleware
	engine.Use(authMiddleware)      // Authentication middleware (innermost)

	return engine
}
