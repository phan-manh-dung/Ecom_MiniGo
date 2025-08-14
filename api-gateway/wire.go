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

// ==== Định nghĩa type mới cho từng middleware để tránh conflict ====
type CORSMiddleware gin.HandlerFunc
type RequestIDMiddleware gin.HandlerFunc
type LoggingMiddleware gin.HandlerFunc

// wireApp khởi tạo toàn bộ application với dependency injection
func wireApp() (*gin.Engine, error) {
	wire.Build(
		// 1. Tạo Service Manager
		service_manager.NewServiceManager,

		// 2. Tạo Handlers với Service Manager
		handler.NewUserServiceClient,
		handler.NewProductServiceClient,
		handler.NewOrderServiceClient,

		// 3. Tạo Global Middleware (mỗi cái một type riêng)
		provideCORSMiddleware,
		provideRequestIDMiddleware,
		provideLoggingMiddleware,

		// 4. Tạo Router
		router.NewRouter,

		// 5. Tạo App với tất cả dependencies
		provideApp,
	)
	return &gin.Engine{}, nil
}

// ==== Provider functions cho từng middleware ====
func provideCORSMiddleware() CORSMiddleware {
	return CORSMiddleware(middleware.NewCORSMiddleware())
}

func provideRequestIDMiddleware() RequestIDMiddleware {
	return RequestIDMiddleware(middleware.NewRequestIDMiddleware())
}

func provideLoggingMiddleware() LoggingMiddleware {
	return LoggingMiddleware(middleware.NewLoggingMiddleware())
}

// ==== Provider App ====
func provideApp(
	serviceManager *service_manager.ServiceManager,
	userHandler *handler.UserServiceClient,
	productHandler *handler.ProductServiceClient,
	orderHandler *handler.OrderServiceClient,
	corsMiddleware CORSMiddleware,
	requestIDMiddleware RequestIDMiddleware,
	loggingMiddleware LoggingMiddleware,
	router *router.Router,
) *gin.Engine {

	// 1. Setup routes với handlers và middleware (đã được xử lý trong router)
	engine := router.SetupRoutes(userHandler, productHandler, orderHandler)

	// 2. Áp dụng global middleware theo thứ tự đúng
	engine.Use(gin.HandlerFunc(corsMiddleware))
	engine.Use(gin.HandlerFunc(requestIDMiddleware))
	engine.Use(gin.HandlerFunc(loggingMiddleware))

	return engine
}
