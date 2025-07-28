package router

import (
	"gin/api-gateway/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, userHandler *handler.UserServiceClient) {
	userRoutes := r.Group("/api/users")
	{
		userRoutes.GET("/id/:id", userHandler.GetUser)
		userRoutes.GET("/sdt/:sdt", userHandler.GetUserBySDT)
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.PUT("/id/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/id/:id", userHandler.DeleteUser)
		userRoutes.GET("/", userHandler.ListUsers)
	}
}
