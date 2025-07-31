package handler

import (
	"context"
	"gin/proto/generated/order"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderServiceClient struct {
	Client order.OrderServiceClient
}

func NewOrderServiceClient(client order.OrderServiceClient) *OrderServiceClient {
	return &OrderServiceClient{Client: client}
}

func (o *OrderServiceClient) GetOrder(c *gin.Context) {
	orderIdParam := c.Param("id")
	orderId, err := strconv.ParseUint(orderIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid orderId"})
		return
	}
	req := &order.GetOrderRequest{Id: uint32(orderId)}
	resp, err := o.Client.GetOrder(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order": resp.Order, "message": resp.Message})
}

func (o *OrderServiceClient) GetOrdersByUser(c *gin.Context) {
	userIdParam := c.Param("user_id")
	userId, err := strconv.ParseUint(userIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}
	req := &order.GetOrdersByUserRequest{UserId: uint32(userId)}
	resp, err := o.Client.GetOrdersByUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"orders": resp.Orders, "message": resp.Message})
}

func (o *OrderServiceClient) CreateOrder(c *gin.Context) {
	var createReq struct {
		UserId uint32 `json:"user_id" binding:"required"`
		Items  []struct {
			ProductId uint32  `json:"product_id" binding:"required"`
			Quantity  uint32  `json:"quantity" binding:"required"`
			UnitPrice float64 `json:"unit_price" binding:"required"`
		} `json:"items" binding:"required"`
	}
	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	var items []*order.OrderItem
	for _, i := range createReq.Items {
		items = append(items, &order.OrderItem{
			ProductId: i.ProductId,
			Quantity:  i.Quantity,
			UnitPrice: i.UnitPrice,
		})
	}
	req := &order.CreateOrderRequest{
		UserId: createReq.UserId,
		Items:  items,
	}
	resp, err := o.Client.CreateOrder(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"order": resp.Order, "message": resp.Message})
}

func (o *OrderServiceClient) UpdateOrderStatus(c *gin.Context) {
	orderIdParam := c.Param("id")
	orderId, err := strconv.ParseUint(orderIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid orderId"})
		return
	}
	var statusReq struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&statusReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req := &order.UpdateOrderStatusRequest{
		OrderId: uint32(orderId),
		Status:  statusReq.Status,
	}
	resp, err := o.Client.UpdateOrderStatus(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order": resp.Order, "message": resp.Message})
}

func (o *OrderServiceClient) GetOrderDetails(c *gin.Context) {
	orderIdParam := c.Param("id")
	orderId, err := strconv.ParseUint(orderIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid orderId"})
		return
	}
	req := &order.GetOrderDetailsRequest{OrderId: uint32(orderId)}
	resp, err := o.Client.GetOrderDetails(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order_details": resp.OrderDetails, "message": resp.Message})
}

func (o *OrderServiceClient) CancelOrder(c *gin.Context) {
	orderIdParam := c.Param("id")
	orderId, err := strconv.ParseUint(orderIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid orderId"})
		return
	}

	req := &order.CancelOrderRequest{
		OrderId: uint32(orderId),
	}
	resp, err := o.Client.CancelOrder(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order": resp.Order, "message": resp.Message})
}
