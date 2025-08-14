package handler

import (
	"context"
	"gin/api-gateway/service_manager"
	"gin/proto/generated/product"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductServiceClient struct {
	ServiceManager *service_manager.ServiceManager
}

func NewProductServiceClient(serviceManager *service_manager.ServiceManager) *ProductServiceClient {
	return &ProductServiceClient{ServiceManager: serviceManager}
}

func (u *ProductServiceClient) GetProduct(c *gin.Context) {
	if u.ServiceManager == nil || u.ServiceManager.ProductClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Product service not available"})
		return
	}

	productIdParam := c.Param("id")
	productId, err := strconv.ParseUint(productIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid productId"})
		return
	}
	req := &product.GetProductRequest{Id: uint32(productId)}
	resp, err := u.ServiceManager.ProductClient.GetProduct(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": resp.Product, "message": resp.Message})
}

func (u *ProductServiceClient) CreateProduct(c *gin.Context) {
	if u.ServiceManager == nil || u.ServiceManager.ProductClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Product service not available"})
		return
	}

	var createReq struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Price       float32 `json:"price" binding:"required"`
		Image       string  `json:"image" binding:"required"`
	}
	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req := &product.CreateProductRequest{
		Name:        createReq.Name,
		Description: createReq.Description,
		Price:       createReq.Price,
		Image:       createReq.Image,
	}
	resp, err := u.ServiceManager.ProductClient.CreateProduct(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"product": resp.Product, "message": resp.Message})
}

func (u *ProductServiceClient) UpdateProduct(c *gin.Context) {
	if u.ServiceManager == nil || u.ServiceManager.ProductClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Product service not available"})
		return
	}

	productIdParam := c.Param("id")
	productId, err := strconv.ParseUint(productIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid productId"})
		return
	}

	var updateReq struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float32 `json:"price"`
		Image       string  `json:"image"`
	}
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	req := &product.UpdateProductRequest{
		Id:          uint32(productId),
		Name:        updateReq.Name,
		Description: updateReq.Description,
		Price:       updateReq.Price,
		Image:       updateReq.Image,
	}
	resp, err := u.ServiceManager.ProductClient.UpdateProduct(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": resp.Product, "message": resp.Message})
}

func (u *ProductServiceClient) DeleteProduct(c *gin.Context) {
	if u.ServiceManager == nil || u.ServiceManager.ProductClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Product service not available"})
		return
	}

	productIdParam := c.Param("id")
	productId, err := strconv.ParseUint(productIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid productId"})
		return
	}

	req := &product.DeleteProductRequest{Id: uint32(productId)}
	resp, err := u.ServiceManager.ProductClient.DeleteProduct(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": resp.Message})
}

func (u *ProductServiceClient) DecreaseInventory(c *gin.Context) {
	if u.ServiceManager == nil || u.ServiceManager.ProductClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Product service not available"})
		return
	}

	var reqBody struct {
		ProductId uint32 `json:"product_id" binding:"required"`
		Quantity  uint32 `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req := &product.DecreaseInventoryRequest{
		ProductId: reqBody.ProductId,
		Quantity:  reqBody.Quantity,
	}
	resp, err := u.ServiceManager.ProductClient.DecreaseInventory(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": resp.Message})
}
