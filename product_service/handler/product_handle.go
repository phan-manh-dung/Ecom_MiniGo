package handler

import (
	"context"
	"gin/product_service/service"
	"gin/proto/generated/product"
	"gin/shared/generic"
)

type ProductHandler struct {
	productService *service.ProductService
	generic        *generic.GenericHandler
	product.UnimplementedProductServiceServer
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		generic:        generic.NewGenericHandler(),
	}
}

// get product
func (h *ProductHandler) GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.GetProductResponse, error) {
	return generic.HandleOperationWithID[uint32, *product.GetProductResponse, uint32](ctx, req.Id, h.productService.GetProduct, "get product")
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	return generic.HandleOperation(ctx, req, h.productService.CreateProduct, "create product")
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (*product.UpdateProductResponse, error) {
	return generic.HandleOperation(ctx, req, h.productService.UpdateProduct, "update product")
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	return generic.HandleOperationWithID[uint32, *product.DeleteProductResponse, uint32](ctx, req.Id, h.productService.DeleteProduct, "delete product")
}

func (h *ProductHandler) DecreaseInventory(ctx context.Context, req *product.DecreaseInventoryRequest) (*product.DecreaseInventoryResponse, error) {
	return generic.HandleOperation(ctx, req, h.productService.DecreaseInventory, "decrease inventory")
}

func (h *ProductHandler) IncreaseInventory(ctx context.Context, req *product.IncreaseInventoryRequest) (*product.IncreaseInventoryResponse, error) {
	return generic.HandleOperation(ctx, req, h.productService.IncreaseInventory, "increase inventory")
}
