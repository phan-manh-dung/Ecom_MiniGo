package handler

import (
	"context"
	"gin/order_service/service"
	"gin/proto/generated/order"
)

type OrderHandler struct {
	orderService *service.OrderService
	order.UnimplementedOrderServiceServer
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	return h.orderService.GetOrder(ctx, req.Id)
}

func (h *OrderHandler) GetOrdersByUser(ctx context.Context, req *order.GetOrdersByUserRequest) (*order.GetOrdersByUserResponse, error) {
	return h.orderService.GetOrdersByUser(ctx, req.UserId)
}

func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error) {
	return h.orderService.UpdateOrderStatus(ctx, req)
}

func (h *OrderHandler) GetOrderDetails(ctx context.Context, req *order.GetOrderDetailsRequest) (*order.GetOrderDetailsResponse, error) {
	return h.orderService.GetOrderDetails(ctx, req)
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	return h.orderService.CreateOrder(ctx, req)
}

func (h *OrderHandler) CancelOrder(ctx context.Context, req *order.CancelOrderRequest) (*order.CancelOrderResponse, error) {
	return h.orderService.CancelOrder(ctx, req)
}
