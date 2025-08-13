package handler

import (
	"context"
	"gin/order_service/service"
	"gin/proto/generated/order"
	"gin/shared/generic"
)

type OrderHandler struct {
	orderService *service.OrderService
	generic      *generic.GenericHandler
	order.UnimplementedOrderServiceServer
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		generic:      generic.NewGenericHandler(),
	}
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	return generic.HandleOperationWithID[uint32, *order.GetOrderResponse, uint32](ctx, req.Id, h.orderService.GetOrder, "get order")
}

func (h *OrderHandler) GetOrdersByUser(ctx context.Context, req *order.GetOrdersByUserRequest) (*order.GetOrdersByUserResponse, error) {
	return generic.HandleOperationWithID[uint32, *order.GetOrdersByUserResponse, uint32](ctx, req.UserId, h.orderService.GetOrdersByUser, "get orders by user")
}

func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error) {
	return generic.HandleOperation(ctx, req, h.orderService.UpdateOrderStatus, "update order status")
}

func (h *OrderHandler) GetOrderDetails(ctx context.Context, req *order.GetOrderDetailsRequest) (*order.GetOrderDetailsResponse, error) {
	return generic.HandleOperation(ctx, req, h.orderService.GetOrderDetails, "get order details")
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	return generic.HandleOperation(ctx, req, h.orderService.CreateOrder, "create order")
}

func (h *OrderHandler) CancelOrder(ctx context.Context, req *order.CancelOrderRequest) (*order.CancelOrderResponse, error) {
	return generic.HandleOperation(ctx, req, h.orderService.CancelOrder, "cancel order")
}
