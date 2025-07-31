package service

import (
	"context"
	"fmt"
	"gin/order_service/model"
	"gin/order_service/repository"
	"gin/proto/generated/order"
	"gin/proto/generated/product"
)

type OrderService struct {
	orderRepo     *repository.OrderRepository
	productClient product.ProductServiceClient
}

func NewOrderService(orderRepo *repository.OrderRepository, productClient product.ProductServiceClient) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		productClient: productClient,
	}
}

func (s *OrderService) GetOrder(ctx context.Context, id uint32) (*order.GetOrderResponse, error) {
	ord, err := s.orderRepo.GetOrderByID(uint(id))
	if err != nil {
		return nil, fmt.Errorf("order not found")
	}
	return &order.GetOrderResponse{
		Order:   convertToProtoOrder(ord),
		Message: "Order found",
	}, nil
}

func (s *OrderService) GetOrdersByUser(ctx context.Context, userID uint32) (*order.GetOrdersByUserResponse, error) {
	orders, err := s.orderRepo.GetOrdersByUserID(uint(userID))
	if err != nil {
		return nil, fmt.Errorf("orders not found")
	}
	var protoOrders []*order.Order
	for _, o := range orders {
		protoOrders = append(protoOrders, convertToProtoOrder(o))
	}
	return &order.GetOrdersByUserResponse{
		Orders:  protoOrders,
		Message: "Orders found",
	}, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error) {
	err := s.orderRepo.UpdateOrderStatus(uint(req.OrderId), req.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to update status")
	}
	ord, _ := s.orderRepo.GetOrderByID(uint(req.OrderId))
	return &order.UpdateOrderStatusResponse{
		Order:   convertToProtoOrder(ord),
		Message: "Status updated",
	}, nil
}

func (s *OrderService) GetOrderDetails(ctx context.Context, req *order.GetOrderDetailsRequest) (*order.GetOrderDetailsResponse, error) {
	details, err := s.orderRepo.GetOrderDetails(uint(req.OrderId))
	if err != nil {
		return nil, fmt.Errorf("order details not found")
	}
	var protoDetails []*order.OrderDetail
	for _, d := range details {
		protoDetails = append(protoDetails, convertToProtoOrderDetail(d))
	}
	return &order.GetOrderDetailsResponse{
		OrderDetails: protoDetails,
		Message:      "Order details found",
	}, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	// Tạo order model
	ord := &model.Order{
		UserID: uint(req.UserId),
		Status: "pending",
	}
	var total float64
	for _, item := range req.Items {
		ord.OrderDetails = append(ord.OrderDetails, &model.OrderDetail{
			ProductID: uint(item.ProductId),
			Quantity:  uint(item.Quantity),
			UnitPrice: float64(item.UnitPrice),
		})
		total += float64(item.UnitPrice) * float64(item.Quantity)
	}
	ord.TotalPrice = total

	// 1. Lưu order vào database trước
	if err := s.orderRepo.CreateOrder(ord); err != nil {
		return nil, fmt.Errorf("failed to create order")
	}

	// 2. Gọi gRPC sang Product Service để giảm tồn kho từng sản phẩm
	for _, item := range req.Items {
		_, err := s.productClient.DecreaseInventory(ctx, &product.DecreaseInventoryRequest{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to decrease inventory for product %d: %v", item.ProductId, err)
		}
	}

	return &order.CreateOrderResponse{
		Order:   convertToProtoOrder(ord),
		Message: "Order created and inventory updated",
	}, nil
}

// Helper: convert model to proto
func convertToProtoOrder(o *model.Order) *order.Order {
	var details []*order.OrderDetail
	for _, d := range o.OrderDetails {
		details = append(details, convertToProtoOrderDetail(d))
	}
	return &order.Order{
		Id:           uint32(o.ID),
		UserId:       uint32(o.UserID),
		TotalPrice:   o.TotalPrice,
		Status:       o.Status,
		OrderDetails: details,
		// Bổ sung created_at, updated_at nếu cần
	}
}

func convertToProtoOrderDetail(d *model.OrderDetail) *order.OrderDetail {
	return &order.OrderDetail{
		Id:        uint32(d.ID),
		OrderId:   uint32(d.OrderID),
		ProductId: uint32(d.ProductID),
		Quantity:  uint32(d.Quantity),
		UnitPrice: d.UnitPrice,
	}
}
