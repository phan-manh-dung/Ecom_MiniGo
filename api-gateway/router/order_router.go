package router

import (
	"gin/api-gateway/handler"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(r *gin.Engine, orderHandler *handler.OrderServiceClient) {
	orderRoutes := r.Group("/api/order")
	{
		orderRoutes.GET("/:id", orderHandler.GetOrder)
		orderRoutes.GET("/user/:user_id", orderHandler.GetOrdersByUser)
		orderRoutes.POST("/", orderHandler.CreateOrder)
		orderRoutes.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		orderRoutes.PUT("/:id/cancel", orderHandler.CancelOrder)
		orderRoutes.GET("/:id/details", orderHandler.GetOrderDetails)
	}
}
