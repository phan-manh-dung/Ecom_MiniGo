package router

import (
	"gin/api-gateway/handler"
	"gin/api-gateway/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, userHandler *handler.UserServiceClient) {
	userRoutes := r.Group("/api/users")
	{
		// Public routes (không cần authentication)
		userRoutes.POST("/login", userHandler.LoginUser)
		userRoutes.POST("/register", userHandler.CreateUser)

		// Protected routes (cần authentication)
		userRoutes.GET("/id/:id", userHandler.GetUser)
		userRoutes.GET("/sdt/:sdt", userHandler.GetUserBySDT)
		userRoutes.PUT("/id/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/id/:id", userHandler.DeleteUser)

		// Admin only routes
		userRoutes.GET("/", middleware.AdminOnlyMiddleware(), userHandler.ListUsers)
	}
}
