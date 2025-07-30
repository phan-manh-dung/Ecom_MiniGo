package router

import (
	"gin/api-gateway/handler"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.Engine, productHandler *handler.ProductServiceClient) {
	productRoutes := r.Group("/api/product")
	{
		productRoutes.GET("/product/:id", productHandler.GetProduct)
		productRoutes.POST("/", productHandler.CreateProduct)
		productRoutes.PUT("/product/:id", productHandler.UpdateProduct)
		productRoutes.DELETE("/product/:id", productHandler.DeleteProduct)
		productRoutes.POST("/decrease-inventory", productHandler.DecreaseInventory)
	}
}
