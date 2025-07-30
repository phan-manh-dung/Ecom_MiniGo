package service

import (
	"context"
	"fmt"
	"gin/product_service/model"
	"gin/product_service/repository"
	"gin/proto/generated/product"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, id uint32) (*product.GetProductResponse, error) {
	productModel, err := s.productRepo.GetProduct(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &product.GetProductResponse{
				Product: nil,
				Message: "Product not found",
			}, nil
		}
		return nil, fmt.Errorf("failed to get product: %v", err)
	}

	// Convert model to proto
	protoProduct := convertToProtoProduct(productModel)

	return &product.GetProductResponse{
		Product: protoProduct,
		Message: "User retrieved successfully",
	}, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	// validation
	if req.Name == "" || req.Description == "" || req.Price <= 0 {
		return nil, fmt.Errorf("Name or Des or Price are required")
	}

	// CreateProduct model save to database
	productModel := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
		Image:       req.Image,
	}

	// error
	if err := s.productRepo.Create(productModel); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// convert to proto
	protoProduct := convertToProtoProduct(productModel)

	return &product.CreateProductResponse{
		Product: protoProduct,
		Message: "Product create successfully",
	}, nil

}

func (s *ProductService) UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (*product.UpdateProductResponse, error) {
	// Get existing product
	existingProduct, err := s.productRepo.GetProduct(uint(req.Id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %v", err)
	}

	// Update fields if provided
	if req.Name != "" {
		existingProduct.Name = req.Name
	}
	if req.Description != "" {
		existingProduct.Description = req.Description
	}
	if req.Price > 0 {
		existingProduct.Price = float64(req.Price)
	}
	if req.Image != "" {
		existingProduct.Image = req.Image
	}

	// Save updated product
	if err := s.productRepo.Update(existingProduct); err != nil {
		return nil, fmt.Errorf("failed to update product: %v", err)
	}

	// Convert to proto
	protoProduct := convertToProtoProduct(existingProduct)

	return &product.UpdateProductResponse{
		Product: protoProduct,
		Message: "Product updated successfully",
	}, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uint32) (*product.DeleteProductResponse, error) {
	// Check if product exists
	_, err := s.productRepo.GetProduct(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %v", err)
	}

	// Delete product
	if err := s.productRepo.Delete(uint(id)); err != nil {
		return nil, fmt.Errorf("failed to delete product: %v", err)
	}

	return &product.DeleteProductResponse{
		Message: "Product deleted successfully",
	}, nil
}

func convertToProtoProduct(p *model.Product) *product.Product {
	protoProduct := &product.Product{
		Id:          uint32(p.ID),
		Name:        p.Name,
		Description: p.Description,
		Price:       float32(p.Price),
		Image:       p.Image,
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
	if p.Inventory != nil {
		protoProduct.Inventory = &product.Inventory{
			ProductId: uint32(p.Inventory.ProductID),
			Quantity:  uint32(p.Inventory.Quantity),
			CreatedAt: timestamppb.New(p.Inventory.CreatedAt),
			UpdatedAt: timestamppb.New(p.Inventory.UpdatedAt),
		}
	}
	return protoProduct
}
