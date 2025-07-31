package handler

import (
	"context"
	"gin/product_service/service"
	"gin/proto/generated/product"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductHandler struct {
	productService *service.ProductService
	product.UnimplementedProductServiceServer
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// get product
func (h *ProductHandler) GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.GetProductResponse, error) {
	response, err := h.productService.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get product")
	}
	return response, nil
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	response, err := h.productService.CreateProduct(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, " Failed to create product")
	}
	return response, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (*product.UpdateProductResponse, error) {
	response, err := h.productService.UpdateProduct(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update product")
	}
	return response, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	response, err := h.productService.DeleteProduct(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete product")
	}
	return response, nil
}

func (h *ProductHandler) DecreaseInventory(ctx context.Context, req *product.DecreaseInventoryRequest) (*product.DecreaseInventoryResponse, error) {
	response, err := h.productService.DecreaseInventory(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to decrease inventory")
	}
	return response, nil
}

func (h *ProductHandler) IncreaseInventory(ctx context.Context, req *product.IncreaseInventoryRequest) (*product.IncreaseInventoryResponse, error) {
	response, err := h.productService.IncreaseInventory(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to increase inventory")
	}
	return response, nil
}
