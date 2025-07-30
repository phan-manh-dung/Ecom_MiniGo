package handler

import (
	"context"
	"gin/proto/generated/product"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductServiceClient struct {
	Client product.ProductServiceClient
}

func NewProductServiceClient(client product.ProductServiceClient) *ProductServiceClient {
	return &ProductServiceClient{Client: client}
}

func (u *ProductServiceClient) GetProduct(c *gin.Context) {
	productIdParam := c.Param("id")
	productId, err := strconv.ParseUint(productIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid productId"})
	}
	req := &product.GetProductRequest{Id: uint32(productId)}
	resp, err := u.Client.GetProduct(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": resp.Product, "message": resp.Message})
}

func (u *ProductServiceClient) CreateProduct(c *gin.Context) {
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
	resp, err := u.Client.CreateProduct(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"product": resp.Product, "message": resp.Message})
}

func (u *ProductServiceClient) UpdateProduct(c *gin.Context) {
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
	resp, err := u.Client.UpdateProduct(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": resp.Product, "message": resp.Message})
}

func (u *ProductServiceClient) DeleteProduct(c *gin.Context) {
	productIdParam := c.Param("id")
	productId, err := strconv.ParseUint(productIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid productId"})
		return
	}

	req := &product.DeleteProductRequest{Id: uint32(productId)}
	resp, err := u.Client.DeleteProduct(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": resp.Message})
}

func (u *ProductServiceClient) DecreaseInventory(c *gin.Context) {
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
	resp, err := u.Client.DecreaseInventory(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": resp.Message})
}
