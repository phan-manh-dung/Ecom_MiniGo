package service

import (
	"gin/proto/generated/product"

	"google.golang.org/grpc"
)

// NewProductServiceClient tạo gRPC client cho Product Service
func NewProductServiceClient(addr string) (product.ProductServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return product.NewProductServiceClient(conn), nil
}
